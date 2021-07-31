// Copyright 2016 The go-bfe Authors
// This file is part of the go-bfe library.
//
// The go-bfe library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-bfe library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-bfe library. If not, see <http://www.gnu.org/licenses/>.

// Package les implements the Light Bfedu Subprotocol.
package les

import (
	"fmt"
	"time"

	"github.com/bfe2021/go-bfe/accounts"
	"github.com/bfe2021/go-bfe/bfe/bfeconfig"
	"github.com/bfe2021/go-bfe/bfe/downloader"
	"github.com/bfe2021/go-bfe/bfe/filters"
	"github.com/bfe2021/go-bfe/bfe/gasprice"
	"github.com/bfe2021/go-bfe/common"
	"github.com/bfe2021/go-bfe/common/hexutil"
	"github.com/bfe2021/go-bfe/common/mclock"
	"github.com/bfe2021/go-bfe/consensus"
	"github.com/bfe2021/go-bfe/core"
	"github.com/bfe2021/go-bfe/core/bloombits"
	"github.com/bfe2021/go-bfe/core/rawdb"
	"github.com/bfe2021/go-bfe/core/types"
	"github.com/bfe2021/go-bfe/event"
	"github.com/bfe2021/go-bfe/internal/bfeapi"
	"github.com/bfe2021/go-bfe/les/vflux"
	vfc "github.com/bfe2021/go-bfe/les/vflux/client"
	"github.com/bfe2021/go-bfe/light"
	"github.com/bfe2021/go-bfe/log"
	"github.com/bfe2021/go-bfe/node"
	"github.com/bfe2021/go-bfe/p2p"
	"github.com/bfe2021/go-bfe/p2p/enode"
	"github.com/bfe2021/go-bfe/p2p/enr"
	"github.com/bfe2021/go-bfe/params"
	"github.com/bfe2021/go-bfe/rlp"
	"github.com/bfe2021/go-bfe/rpc"
)

type LightBfedu struct {
	lesCommons

	peers              *serverPeerSet
	reqDist            *requestDistributor
	retriever          *retrieveManager
	odr                *LesOdr
	relay              *lesTxRelay
	handler            *clientHandler
	txPool             *light.TxPool
	blockchain         *light.LightChain
	serverPool         *vfc.ServerPool
	serverPoolIterator enode.Iterator
	pruner             *pruner

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	ApiBackend     *LesApiBackend
	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager
	netRPCService  *bfeapi.PublicNetAPI

	p2pServer *p2p.Server
	p2pConfig *p2p.Config
}

// New creates an instance of the light client.
func New(stack *node.Node, config *bfeconfig.Config) (*LightBfedu, error) {
	chainDb, err := stack.OpenDatabase("lightchaindata", config.DatabaseCache, config.DatabaseHandles, "bfe/db/chaindata/")
	if err != nil {
		return nil, err
	}
	lesDb, err := stack.OpenDatabase("les.client", 0, 0, "bfe/db/lesclient/")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlockWithOverride(chainDb, config.Genesis, config.OverrideBerlin)
	if _, isCompat := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !isCompat {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	peers := newServerPeerSet()
	long := &LightBfedu{
		lesCommons: lesCommons{
			genesis:     genesisHash,
			config:      config,
			chainConfig: chainConfig,
			iConfig:     light.DefaultClientIndexerConfig,
			chainDb:     chainDb,
			lesDb:       lesDb,
			closeCh:     make(chan struct{}),
		},
		peers:          peers,
		eventMux:       stack.EventMux(),
		reqDist:        newRequestDistributor(peers, &mclock.System{}),
		accountManager: stack.AccountManager(),
		engine:         bfeconfig.CreateConsensusEngine(stack, chainConfig, &config.Bfeash, nil, false, chainDb),
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   core.NewBloomIndexer(chainDb, params.BloomBitsBlocksClient, params.HelperTrieConfirmations),
		p2pServer:      stack.Server(),
		p2pConfig:      &stack.Config().P2P,
	}

	var prenegQuery vfc.QueryFunc
	if lbfe.p2pServer.DiscV5 != nil {
		prenegQuery = lbfe.prenegQuery
	}
	lbfe.serverPool, lbfe.serverPoolIterator = vfc.NewServerPool(lesDb, []byte("serverpool:"), time.Second, prenegQuery, &mclock.System{}, config.UltraLightServers, requestList)
	lbfe.serverPool.AddMetrics(suggestedTimeoutGauge, totalValueGauge, serverSelectableGauge, serverConnectedGauge, sessionValueMeter, serverDialedMeter)

	lbfe.retriever = newRetrieveManager(peers, lbfe.reqDist, lbfe.serverPool.GetTimeout)
	lbfe.relay = newLesTxRelay(peers, lbfe.retriever)

	lbfe.odr = NewLesOdr(chainDb, light.DefaultClientIndexerConfig, lbfe.peers, lbfe.retriever)
	lbfe.chtIndexer = light.NewChtIndexer(chainDb, lbfe.odr, params.CHTFrequency, params.HelperTrieConfirmations, config.LightNoPrune)
	lbfe.bloomTrieIndexer = light.NewBloomTrieIndexer(chainDb, lbfe.odr, params.BloomBitsBlocksClient, params.BloomTrieFrequency, config.LightNoPrune)
	lbfe.odr.SetIndexers(lbfe.chtIndexer, lbfe.bloomTrieIndexer, lbfe.bloomIndexer)

	checkpoint := config.Checkpoint
	if checkpoint == nil {
		checkpoint = params.TrustedCheckpoints[genesisHash]
	}
	// Note: NewLightChain adds the trusted checkpoint so it needs an ODR with
	// indexers already set but not started yet
	if lbfe.blockchain, err = light.NewLightChain(lbfe.odr, lbfe.chainConfig, lbfe.engine, checkpoint); err != nil {
		return nil, err
	}
	lbfe.chainReader = lbfe.blockchain
	lbfe.txPool = light.NewTxPool(lbfe.chainConfig, lbfe.blockchain, lbfe.relay)

	// Set up checkpoint oracle.
	lbfe.oracle = lbfe.setupOracle(stack, genesisHash, config)

	// Note: AddChildIndexer starts the update process for the child
	lbfe.bloomIndexer.AddChildIndexer(lbfe.bloomTrieIndexer)
	lbfe.chtIndexer.Start(lbfe.blockchain)
	lbfe.bloomIndexer.Start(lbfe.blockchain)

	// Start a light chain pruner to delete useless historical data.
	lbfe.pruner = newPruner(chainDb, lbfe.chtIndexer, lbfe.bloomTrieIndexer)

	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		lbfe.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}

	lbfe.ApiBackend = &LesApiBackend{stack.Config().ExtRPCEnabled(), stack.Config().AllowUnprotectedTxs, long, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.Miner.GasPrice
	}
	lbfe.ApiBackend.gpo = gasprice.NewOracle(lbfe.ApiBackend, gpoParams)

	lbfe.handler = newClientHandler(config.UltraLightServers, config.UltraLightFraction, checkpoint, long)
	if lbfe.handler.ulc != nil {
		log.Warn("Ultra light client is enabled", "trustedNodes", len(lbfe.handler.ulc.keys), "minTrustedFraction", lbfe.handler.ulc.fraction)
		lbfe.blockchain.DisableCheckFreq()
	}

	lbfe.netRPCService = bfeapi.NewPublicNetAPI(lbfe.p2pServer, lbfe.config.NetworkId)

	// Register the backend on the node
	stack.RegisterAPIs(lbfe.APIs())
	stack.RegisterProtocols(lbfe.Protocols())
	stack.RegisterLifecycle(long)

	// Check for unclean shutdown
	if uncleanShutdowns, discards, err := rawdb.PushUncleanShutdownMarker(chainDb); err != nil {
		log.Error("Could not update unclean-shutdown-marker list", "error", err)
	} else {
		if discards > 0 {
			log.Warn("Old unclean shutdowns found", "count", discards)
		}
		for _, tstamp := range uncleanShutdowns {
			t := time.Unix(int64(tstamp), 0)
			log.Warn("Unclean shutdown detected", "booted", t,
				"age", common.PrettyAge(t))
		}
	}
	return long, nil
}

// VfluxRequest sends a batch of requests to the given node through discv5 UDP TalkRequest and returns the responses
func (s *LightBfedu) VfluxRequest(n *enode.Node, reqs vflux.Requests) vflux.Replies {
	if s.p2pServer.DiscV5 == nil {
		return nil
	}
	reqsEnc, _ := rlp.EncodeToBytes(&reqs)
	repliesEnc, _ := s.p2pServer.DiscV5.TalkRequest(s.serverPool.DialNode(n), "vfx", reqsEnc)
	var replies vflux.Replies
	if len(repliesEnc) == 0 || rlp.DecodeBytes(repliesEnc, &replies) != nil {
		return nil
	}
	return replies
}

// vfxVersion returns the version number of the "les" service subdomain of the vflux UDP
// service, as advertised in the ENR record
func (s *LightBfedu) vfxVersion(n *enode.Node) uint {
	if n.Seq() == 0 {
		var err error
		if s.p2pServer.DiscV5 == nil {
			return 0
		}
		if n, err = s.p2pServer.DiscV5.RequestENR(n); n != nil && err == nil && n.Seq() != 0 {
			s.serverPool.Persist(n)
		} else {
			return 0
		}
	}

	var les []rlp.RawValue
	if err := n.Load(enr.WithEntry("les", &les)); err != nil || len(les) < 1 {
		return 0
	}
	var version uint
	rlp.DecodeBytes(les[0], &version) // Ignore additional fields (for forward compatibility).
	return version
}

// prenegQuery sends a capacity query to the given server node to determine whether
// a connection slot is immediately available
func (s *LightBfedu) prenegQuery(n *enode.Node) int {
	if s.vfxVersion(n) < 1 {
		// UDP query not supported, always try TCP connection
		return 1
	}

	var requests vflux.Requests
	requests.Add("les", vflux.CapacityQueryName, vflux.CapacityQueryReq{
		Bias:      180,
		AddTokens: []vflux.IntOrInf{{}},
	})
	replies := s.VfluxRequest(n, requests)
	var cqr vflux.CapacityQueryReply
	if replies.Get(0, &cqr) != nil || len(cqr) != 1 { // Note: Get returns an error if replies is nil
		return -1
	}
	if cqr[0] > 0 {
		return 1
	}
	return 0
}

type LightDummyAPI struct{}

// Bfeerbase is the address that mining rewards will be send to
func (s *LightDummyAPI) Bfeerbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("mining is not supported in light mode")
}

// Coinbase is the address that mining rewards will be send to (alias for Bfeerbase)
func (s *LightDummyAPI) Coinbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("mining is not supported in light mode")
}

// Hashrate returns the POW hashrate
func (s *LightDummyAPI) Hashrate() hexutil.Uint {
	return 0
}

// Mining returns an indication if this node is currently mining.
func (s *LightDummyAPI) Mining() bool {
	return false
}

// APIs returns the collection of RPC services the bfedu package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *LightBfedu) APIs() []rpc.API {
	apis := bfeapi.GetAPIs(s.ApiBackend)
	apis = append(apis, s.engine.APIs(s.BlockChain().HeaderChain())...)
	return append(apis, []rpc.API{
		{
			Namespace: "ong",
			Version:   "1.0",
			Service:   &LightDummyAPI{},
			Public:    true,
		}, {
			Namespace: "ong",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.handler.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "ong",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, true, 5*time.Minute),
			Public:    true,
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		}, {
			Namespace: "les",
			Version:   "1.0",
			Service:   NewPrivateLightAPI(&s.lesCommons),
			Public:    false,
		}, {
			Namespace: "vflux",
			Version:   "1.0",
			Service:   s.serverPool.API(),
			Public:    false,
		},
	}...)
}

func (s *LightBfedu) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *LightBfedu) BlockChain() *light.LightChain      { return s.blockchain }
func (s *LightBfedu) TxPool() *light.TxPool              { return s.txPool }
func (s *LightBfedu) Engine() consensus.Engine           { return s.engine }
func (s *LightBfedu) LesVersion() int                    { return int(ClientProtocolVersions[0]) }
func (s *LightBfedu) Downloader() *downloader.Downloader { return s.handler.downloader }
func (s *LightBfedu) EventMux() *event.TypeMux           { return s.eventMux }

// Protocols returns all the currently configured network protocols to start.
func (s *LightBfedu) Protocols() []p2p.Protocol {
	return s.makeProtocols(ClientProtocolVersions, s.handler.runPeer, func(id enode.ID) interface{} {
		if p := s.peers.peer(id.String()); p != nil {
			return p.Info()
		}
		return nil
	}, s.serverPoolIterator)
}

// Start implements node.Lifecycle, starting all internal goroutines needed by the
// light bfedu protocol implementation.
func (s *LightBfedu) Start() error {
	log.Warn("Light client mode is an experimental feature")

	discovery, err := s.setupDiscovery(s.p2pConfig)
	if err != nil {
		return err
	}
	s.serverPool.AddSource(discovery)
	s.serverPool.Start()
	// Start bloom request workers.
	s.wg.Add(bloomServiceThreads)
	s.startBloomHandlers(params.BloomBitsBlocksClient)
	s.handler.start()

	return nil
}

// Stop implements node.Lifecycle, terminating all internal goroutines used by the
// Bfedu protocol.
func (s *LightBfedu) Stop() error {
	close(s.closeCh)
	s.serverPool.Stop()
	s.peers.close()
	s.reqDist.close()
	s.odr.Stop()
	s.relay.Stop()
	s.bloomIndexer.Close()
	s.chtIndexer.Close()
	s.blockchain.Stop()
	s.handler.stop()
	s.txPool.Stop()
	s.engine.Close()
	s.pruner.close()
	s.eventMux.Stop()
	rawdb.PopUncleanShutdownMarker(s.chainDb)
	s.chainDb.Close()
	s.lesDb.Close()
	s.wg.Wait()
	log.Info("Light bfedu stopped")
	return nil
}
