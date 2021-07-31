// Copyright 2014 The go-bfe Authors
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

// Package  bfe  implements the Bfedu protocol.
package bfe

import (
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bfe2021/go-bfe/accounts"
	"github.com/bfe2021/go-bfe/bfe/bfeconfig"
	"github.com/bfe2021/go-bfe/bfe/downloader"
	"github.com/bfe2021/go-bfe/bfe/filters"
	"github.com/bfe2021/go-bfe/bfe/gasprice"
	"github.com/bfe2021/go-bfe/bfe/protocols/bfe"
	"github.com/bfe2021/go-bfe/bfe/protocols/snap"
	"github.com/bfe2021/go-bfe/bfedb"
	"github.com/bfe2021/go-bfe/common"
	"github.com/bfe2021/go-bfe/common/hexutil"
	"github.com/bfe2021/go-bfe/consensus"
	"github.com/bfe2021/go-bfe/consensus/clique"
	"github.com/bfe2021/go-bfe/core"
	"github.com/bfe2021/go-bfe/core/bloombits"
	"github.com/bfe2021/go-bfe/core/rawdb"
	"github.com/bfe2021/go-bfe/core/state/pruner"
	"github.com/bfe2021/go-bfe/core/types"
	"github.com/bfe2021/go-bfe/core/vm"
	"github.com/bfe2021/go-bfe/event"
	"github.com/bfe2021/go-bfe/internal/bfeapi"
	"github.com/bfe2021/go-bfe/log"
	"github.com/bfe2021/go-bfe/miner"
	"github.com/bfe2021/go-bfe/node"
	"github.com/bfe2021/go-bfe/p2p"
	"github.com/bfe2021/go-bfe/p2p/enode"
	"github.com/bfe2021/go-bfe/params"
	"github.com/bfe2021/go-bfe/rlp"
	"github.com/bfe2021/go-bfe/rpc"
)

// Config contains the configuration options of the BFE protocol.
// Deprecated: use bfeconfig.Config instead.
type Config = bfeconfig.Config

// Bfedu implements the Bfedu full node service.
type Bfedu struct {
	config *bfeconfig.Config

	// Handlers
	txPool             *core.TxPool
	blockchain         *core.BlockChain
	handler            *handler
	bfeDialCandidates  enode.Iterator
	snapDialCandidates enode.Iterator

	// DB interfaces
	chainDb bfedb.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests     chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer      *core.ChainIndexer             // Bloom indexer operating during block imports
	closeBloomHandler chan struct{}

	APIBackend *BfeAPIBackend

	miner     *miner.Miner
	gasPrice  *big.Int
	bfeerbase common.Address

	networkID     uint64
	netRPCService *bfeapi.PublicNetAPI

	p2pServer *p2p.Server

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and bfeerbase)
}

// New creates a new Bfedu object (including the
// initialisation of the common Bfedu object)
func New(stack *node.Node, config *bfeconfig.Config) (*Bfedu, error) {
	// Ensure configuration values are compatible and sane
	if config.SyncMode == downloader.LightSync {
		return nil, errors.New("can't run bfe.Bfedu in light sync mode, use les.LightBfedu")
	}
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	if config.Miner.GasPrice == nil || config.Miner.GasPrice.Cmp(common.Big0) <= 0 {
		log.Warn("Sanitizing invalid miner gas price", "provided", config.Miner.GasPrice, "updated", bfeconfig.Defaults.Miner.GasPrice)
		config.Miner.GasPrice = new(big.Int).Set(bfeconfig.Defaults.Miner.GasPrice)
	}
	if config.NoPruning && config.TrieDirtyCache > 0 {
		if config.SnapshotCache > 0 {
			config.TrieCleanCache += config.TrieDirtyCache * 3 / 5
			config.SnapshotCache += config.TrieDirtyCache * 2 / 5
		} else {
			config.TrieCleanCache += config.TrieDirtyCache
		}
		config.TrieDirtyCache = 0
	}
	log.Info("Allocated trie memory caches", "clean", common.StorageSize(config.TrieCleanCache)*1024*1024, "dirty", common.StorageSize(config.TrieDirtyCache)*1024*1024)

	// Assemble the Bfedu object
	chainDb, err := stack.OpenDatabaseWithFreezer("chaindata", config.DatabaseCache, config.DatabaseHandles, config.DatabaseFreezer, "bfe/db/chaindata/")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlockWithOverride(chainDb, config.Genesis, config.OverrideBerlin)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	if err := pruner.RecoverPruning(stack.ResolvePath(""), chainDb, stack.ResolvePath(config.TrieCleanCacheJournal)); err != nil {
		log.Error("Failed to recover state", "error", err)
	}
	bfe := &Bfedu{
		config:            config,
		chainDb:           chainDb,
		eventMux:          stack.EventMux(),
		accountManager:    stack.AccountManager(),
		engine:            bfeconfig.CreateConsensusEngine(stack, chainConfig, &config.Bfeash, config.Miner.Notify, config.Miner.Noverify, chainDb),
		closeBloomHandler: make(chan struct{}),
		networkID:         config.NetworkId,
		gasPrice:          config.Miner.GasPrice,
		bfeerbase:         config.Miner.Bfeerbase,
		bloomRequests:     make(chan chan *bloombits.Retrieval),
		bloomIndexer:      core.NewBloomIndexer(chainDb, params.BloomBitsBlocks, params.BloomConfirms),
		p2pServer:         stack.Server(),
	}

	bcVersion := rawdb.ReadDatabaseVersion(chainDb)
	var dbVer = "<nil>"
	if bcVersion != nil {
		dbVer = fmt.Sprintf("%d", *bcVersion)
	}
	log.Info("Initialising Bfedu protocol", "network", config.NetworkId, "dbversion", dbVer)

	if !config.SkipBcVersionCheck {
		if bcVersion != nil && *bcVersion > core.BlockChainVersion {
			return nil, fmt.Errorf("database version is v%d, Gbfe %s only supports v%d", *bcVersion, params.VersionWithMeta, core.BlockChainVersion)
		} else if bcVersion == nil || *bcVersion < core.BlockChainVersion {
			log.Warn("Upgrade blockchain database version", "from", dbVer, "to", core.BlockChainVersion)
			rawdb.WriteDatabaseVersion(chainDb, core.BlockChainVersion)
		}
	}
	var (
		vmConfig = vm.Config{
			EnablePreimageRecording: config.EnablePreimageRecording,
			EWASMInterpreter:        config.EWASMInterpreter,
			EVMInterpreter:          config.EVMInterpreter,
		}
		cacheConfig = &core.CacheConfig{
			TrieCleanLimit:      config.TrieCleanCache,
			TrieCleanJournal:    stack.ResolvePath(config.TrieCleanCacheJournal),
			TrieCleanRejournal:  config.TrieCleanCacheRejournal,
			TrieCleanNoPrefetch: config.NoPrefetch,
			TrieDirtyLimit:      config.TrieDirtyCache,
			TrieDirtyDisabled:   config.NoPruning,
			TrieTimeLimit:       config.TrieTimeout,
			SnapshotLimit:       config.SnapshotCache,
			Preimages:           config.Preimages,
		}
	)
	bfe.blockchain, err = core.NewBlockChain(chainDb, cacheConfig, chainConfig, bfe.engine, vmConfig, bfe.shouldPreserve, &config.TxLookupLimit)
	if err != nil {
		return nil, err
	}
	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		bfe.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}
	bfe.bloomIndexer.Start(bfe.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = stack.ResolvePath(config.TxPool.Journal)
	}
	bfe.txPool = core.NewTxPool(config.TxPool, chainConfig, bfe.blockchain)

	// Permit the downloader to use the trie cache allowance during fast sync
	cacheLimit := cacheConfig.TrieCleanLimit + cacheConfig.TrieDirtyLimit + cacheConfig.SnapshotLimit
	checkpoint := config.Checkpoint
	if checkpoint == nil {
		checkpoint = params.TrustedCheckpoints[genesisHash]
	}
	if bfe.handler, err = newHandler(&handlerConfig{
		Database:   chainDb,
		Chain:      bfe.blockchain,
		TxPool:     bfe.txPool,
		Network:    config.NetworkId,
		Sync:       config.SyncMode,
		BloomCache: uint64(cacheLimit),
		EventMux:   bfe.eventMux,
		Checkpoint: checkpoint,
		Whitelist:  config.Whitelist,
	}); err != nil {
		return nil, err
	}
	bfe.miner = miner.New(bfe, &config.Miner, chainConfig, bfe.EventMux(), bfe.engine, bfe.isLocalBlock)
	bfe.miner.SetExtra(makeExtraData(config.Miner.ExtraData))

	bfe.APIBackend = &BfeAPIBackend{stack.Config().ExtRPCEnabled(), stack.Config().AllowUnprotectedTxs, bfe, nil}
	if bfe.APIBackend.allowUnprotectedTxs {
		log.Info("Unprotected transactions allowed")
	}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.Miner.GasPrice
	}
	bfe.APIBackend.gpo = gasprice.NewOracle(bfe.APIBackend, gpoParams)

	bfe.bfeDialCandidates, err = setupDiscovery(bfe.config.BfeDiscoveryURLs)
	if err != nil {
		return nil, err
	}
	bfe.snapDialCandidates, err = setupDiscovery(bfe.config.SnapDiscoveryURLs)
	if err != nil {
		return nil, err
	}
	// Start the RPC service
	bfe.netRPCService = bfeapi.NewPublicNetAPI(bfe.p2pServer, config.NetworkId)

	// Register the backend on the node
	stack.RegisterAPIs(bfe.APIs())
	stack.RegisterProtocols(bfe.Protocols())
	stack.RegisterLifecycle(bfe)
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
	return bfe, nil
}

func makeExtraData(extra []byte) []byte {
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			"gbfe",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.MaximumExtraDataSize)
		extra = nil
	}
	return extra
}

// APIs return the collection of RPC services the bfedu package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *Bfedu) APIs() []rpc.API {
	apis := bfeapi.GetAPIs(s.APIBackend)

	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	// Append all the local APIs and return
	return append(apis, []rpc.API{
		{
			Namespace: "bfe",
			Version:   "1.0",
			Service:   NewPublicBfeduAPI(s),
			Public:    true,
		}, {
			Namespace: "bfe",
			Version:   "1.0",
			Service:   NewPublicMinerAPI(s),
			Public:    true,
		}, {
			Namespace: "bfe",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.handler.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "bfe",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.APIBackend, false, 5*time.Minute),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *Bfedu) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *Bfedu) Bfeerbase() (eb common.Address, err error) {
	s.lock.RLock()
	bfeerbase := s.bfeerbase
	s.lock.RUnlock()

	if bfeerbase != (common.Address{}) {
		return bfeerbase, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			bfeerbase := accounts[0].Address

			s.lock.Lock()
			s.bfeerbase = bfeerbase
			s.lock.Unlock()

			log.Info("Bfeerbase automatically configured", "address", bfeerbase)
			return bfeerbase, nil
		}
	}
	return common.Address{}, fmt.Errorf("bfeerbase must be explicitly specified")
}

// isLocalBlock checks whether the specified block is mined
// by local miner accounts.
//
// We regard two types of accounts as local miner account: bfeerbase
// and accounts specified via `txpool.locals` flag.
func (s *Bfedu) isLocalBlock(block *types.Block) bool {
	author, err := s.engine.Author(block.Header())
	if err != nil {
		log.Warn("Failed to retrieve block author", "number", block.NumberU64(), "hash", block.Hash(), "err", err)
		return false
	}
	// Check whether the given address is bfeerbase.
	s.lock.RLock()
	bfeerbase := s.bfeerbase
	s.lock.RUnlock()
	if author == bfeerbase {
		return true
	}
	// Check whether the given address is specified by `txpool.local`
	// CLI flag.
	for _, account := range s.config.TxPool.Locals {
		if account == author {
			return true
		}
	}
	return false
}

// shouldPreserve checks whether we should preserve the given block
// during the chain reorg depending on whether the author of block
// is a local account.
func (s *Bfedu) shouldPreserve(block *types.Block) bool {
	// The reason we need to disable the self-reorg preserving for clique
	// is it can be probable to introduce a deadlock.
	//
	// e.g. If there are 7 available signers
	//
	// r1   A
	// r2     B
	// r3       C
	// r4         D
	// r5   A      [X] F G
	// r6    [X]
	//
	// In the round5, the inturn signer E is offline, so the worst case
	// is A, F and G sign the block of round5 and reject the block of opponents
	// and in the round6, the last available signer B is offline, the whole
	// network is stuck.
	if _, ok := s.engine.(*clique.Clique); ok {
		return false
	}
	return s.isLocalBlock(block)
}

// SetBfeerbase sets the mining reward address.
func (s *Bfedu) SetBfeerbase(bfeerbase common.Address) {
	s.lock.Lock()
	s.bfeerbase = bfeerbase
	s.lock.Unlock()

	s.miner.SetBfeerbase(bfeerbase)
}

// StartMining starts the miner with the given number of CPU threads. If mining
// is already running, this Method adjust the number of threads allowed to use
// and updates the minimum price required by the transaction pool.
func (s *Bfedu) StartMining(threads int) error {
	// Update the thread count within the consensus engine
	type threaded interface {
		SetThreads(threads int)
	}
	if th, ok := s.engine.(threaded); ok {
		log.Info("Updated mining threads", "threads", threads)
		if threads == 0 {
			threads = -1 // Disable the miner from within
		}
		th.SetThreads(threads)
	}
	// If the miner was not running, initialize it
	if !s.IsMining() {
		// Propagate the initial price point to the transaction pool
		s.lock.RLock()
		price := s.gasPrice
		s.lock.RUnlock()
		s.txPool.SetGasPrice(price)

		// Configure the local mining address
		eb, err := s.Bfeerbase()
		if err != nil {
			log.Error("Cannot start mining without ongerbase", "err", err)
			return fmt.Errorf("ongerbase missing: %v", err)
		}
		if clique, ok := s.engine.(*clique.Clique); ok {
			wallet, err := s.accountManager.Find(accounts.Account{Address: eb})
			if wallet == nil || err != nil {
				log.Error("Bfeerbase account unavailable locally", "err", err)
				return fmt.Errorf("signer missing: %v", err)
			}
			clique.Authorize(eb, wallet.SignData)
		}
		// If mining is started, we can disable the transaction rejection mechanism
		// introduced to speed sync times.
		atomic.StoreUint32(&s.handler.acceptTxs, 1)

		go s.miner.Start(eb)
	}
	return nil
}

// StopMining terminates the miner, both at the consensus engine level as well as
// at the block creation level.
func (s *Bfedu) StopMining() {
	// Update the thread count within the consensus engine
	type threaded interface {
		SetThreads(threads int)
	}
	if th, ok := s.engine.(threaded); ok {
		th.SetThreads(-1)
	}
	// Stop the block creating itself
	s.miner.Stop()
}

func (s *Bfedu) IsMining() bool      { return s.miner.Mining() }
func (s *Bfedu) Miner() *miner.Miner { return s.miner }

func (s *Bfedu) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *Bfedu) BlockChain() *core.BlockChain       { return s.blockchain }
func (s *Bfedu) TxPool() *core.TxPool               { return s.txPool }
func (s *Bfedu) EventMux() *event.TypeMux           { return s.eventMux }
func (s *Bfedu) Engine() consensus.Engine           { return s.engine }
func (s *Bfedu) ChainDb() bfedb.Database            { return s.chainDb }
func (s *Bfedu) IsListening() bool                  { return true } // Always listening
func (s *Bfedu) Downloader() *downloader.Downloader { return s.handler.downloader }
func (s *Bfedu) Synced() bool                       { return atomic.LoadUint32(&s.handler.acceptTxs) == 1 }
func (s *Bfedu) ArchiveMode() bool                  { return s.config.NoPruning }
func (s *Bfedu) BloomIndexer() *core.ChainIndexer   { return s.bloomIndexer }

// Protocols returns all the currently configured
// network protocols to start.
func (s *Bfedu) Protocols() []p2p.Protocol {
	protos := bfe.MakeProtocols((*bfeHandler)(s.handler), s.networkID, s.bfeDialCandidates)
	if s.config.SnapshotCache > 0 {
		protos = append(protos, snap.MakeProtocols((*snapHandler)(s.handler), s.snapDialCandidates)...)
	}
	return protos
}

// Start implements node.Lifecycle, starting all internal goroutines needed by the
// Bfedu protocol implementation.
func (s *Bfedu) Start() error {
	bfe.StartENRUpdater(s.blockchain, s.p2pServer.LocalNode())

	// Start the bloom bits servicing goroutines
	s.startBloomHandlers(params.BloomBitsBlocks)

	// Figure out a max peers count based on the server limits
	maxPeers := s.p2pServer.MaxPeers
	if s.config.LightServ > 0 {
		if s.config.LightPeers >= s.p2pServer.MaxPeers {
			return fmt.Errorf("invalid peer config: light peer count (%d) >= total peer count (%d)", s.config.LightPeers, s.p2pServer.MaxPeers)
		}
		maxPeers -= s.config.LightPeers
	}
	// Start the networking layer and the light server if requested
	s.handler.Start(maxPeers)
	return nil
}

// Stop implements node.Lifecycle, terminating all internal goroutines used by the
// Bfedu protocol.
func (s *Bfedu) Stop() error {
	// Stop all the peer-related stuff first.
	s.handler.Stop()

	// Then stop everything else.
	s.bloomIndexer.Close()
	close(s.closeBloomHandler)
	s.txPool.Stop()
	s.miner.Stop()
	s.blockchain.Stop()
	s.engine.Close()
	rawdb.PopUncleanShutdownMarker(s.chainDb)
	s.chainDb.Close()
	s.eventMux.Stop()

	return nil
}
