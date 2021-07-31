// Copyright 2015 The go-bfe Authors
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

package bfe

import (
	"context"
	"errors"
	"math/big"

	"github.com/bfe2021/go-bfe/accounts"
	"github.com/bfe2021/go-bfe/bfe/downloader"
	"github.com/bfe2021/go-bfe/bfe/gasprice"
	"github.com/bfe2021/go-bfe/bfedb"
	"github.com/bfe2021/go-bfe/common"
	"github.com/bfe2021/go-bfe/consensus"
	"github.com/bfe2021/go-bfe/core"
	"github.com/bfe2021/go-bfe/core/bloombits"
	"github.com/bfe2021/go-bfe/core/rawdb"
	"github.com/bfe2021/go-bfe/core/state"
	"github.com/bfe2021/go-bfe/core/types"
	"github.com/bfe2021/go-bfe/core/vm"
	"github.com/bfe2021/go-bfe/event"
	"github.com/bfe2021/go-bfe/miner"
	"github.com/bfe2021/go-bfe/params"
	"github.com/bfe2021/go-bfe/rpc"
)

// BfeAPIBackend implements bfeapi.Backend for full nodes
type BfeAPIBackend struct {
	extRPCEnabled       bool
	allowUnprotectedTxs bool
	bfe                 *Bfedu
	gpo                 *gasprice.Oracle
}

// ChainConfig returns the active chain configuration.
func (b *BfeAPIBackend) ChainConfig() *params.ChainConfig {
	return b.bfe.blockchain.Config()
}

func (b *BfeAPIBackend) CurrentBlock() *types.Block {
	return b.bfe.blockchain.CurrentBlock()
}

func (b *BfeAPIBackend) SetHead(number uint64) {
	b.bfe.handler.downloader.Cancel()
	b.bfe.blockchain.SetHead(number)
}

func (b *BfeAPIBackend) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if number == rpc.PendingBlockNumber {
		block := b.bfe.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if number == rpc.LatestBlockNumber {
		return b.bfe.blockchain.CurrentBlock().Header(), nil
	}
	return b.bfe.blockchain.GetHeaderByNumber(uint64(number)), nil
}

func (b *BfeAPIBackend) HeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.HeaderByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header := b.bfe.blockchain.GetHeaderByHash(hash)
		if header == nil {
			return nil, errors.New("header for hash not found")
		}
		if blockNrOrHash.RequireCanonical && b.bfe.blockchain.GetCanonicalHash(header.Number.Uint64()) != hash {
			return nil, errors.New("hash is not currently canonical")
		}
		return header, nil
	}
	return nil, errors.New("invalid arguments; neither block nor hash specified")
}

func (b *BfeAPIBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return b.bfe.blockchain.GetHeaderByHash(hash), nil
}

func (b *BfeAPIBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if number == rpc.PendingBlockNumber {
		block := b.bfe.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if number == rpc.LatestBlockNumber {
		return b.bfe.blockchain.CurrentBlock(), nil
	}
	return b.bfe.blockchain.GetBlockByNumber(uint64(number)), nil
}

func (b *BfeAPIBackend) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.bfe.blockchain.GetBlockByHash(hash), nil
}

func (b *BfeAPIBackend) BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.BlockByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header := b.bfe.blockchain.GetHeaderByHash(hash)
		if header == nil {
			return nil, errors.New("header for hash not found")
		}
		if blockNrOrHash.RequireCanonical && b.bfe.blockchain.GetCanonicalHash(header.Number.Uint64()) != hash {
			return nil, errors.New("hash is not currently canonical")
		}
		block := b.bfe.blockchain.GetBlock(hash, header.Number.Uint64())
		if block == nil {
			return nil, errors.New("header found, but block body is missing")
		}
		return block, nil
	}
	return nil, errors.New("invalid arguments; neither block nor hash specified")
}

func (b *BfeAPIBackend) StateAndHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if number == rpc.PendingBlockNumber {
		block, state := b.bfe.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, nil, err
	}
	if header == nil {
		return nil, nil, errors.New("header not found")
	}
	stateDb, err := b.bfe.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

func (b *BfeAPIBackend) StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*state.StateDB, *types.Header, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.StateAndHeaderByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header, err := b.HeaderByHash(ctx, hash)
		if err != nil {
			return nil, nil, err
		}
		if header == nil {
			return nil, nil, errors.New("header for hash not found")
		}
		if blockNrOrHash.RequireCanonical && b.bfe.blockchain.GetCanonicalHash(header.Number.Uint64()) != hash {
			return nil, nil, errors.New("hash is not currently canonical")
		}
		stateDb, err := b.bfe.BlockChain().StateAt(header.Root)
		return stateDb, header, err
	}
	return nil, nil, errors.New("invalid arguments; neither block nor hash specified")
}

func (b *BfeAPIBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	return b.bfe.blockchain.GetReceiptsByHash(hash), nil
}

func (b *BfeAPIBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	receipts := b.bfe.blockchain.GetReceiptsByHash(hash)
	if receipts == nil {
		return nil, nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (b *BfeAPIBackend) GetTd(ctx context.Context, hash common.Hash) *big.Int {
	return b.bfe.blockchain.GetTdByHash(hash)
}

func (b *BfeAPIBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header) (*vm.EVM, func() error, error) {
	vmError := func() error { return nil }

	txContext := core.NewEVMTxContext(msg)
	context := core.NewEVMBlockContext(header, b.bfe.BlockChain(), nil)
	return vm.NewEVM(context, txContext, state, b.bfe.blockchain.Config(), *b.bfe.blockchain.GetVMConfig()), vmError, nil
}

func (b *BfeAPIBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.bfe.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *BfeAPIBackend) SubscribePendingLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.bfe.miner.SubscribePendingLogs(ch)
}

func (b *BfeAPIBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.bfe.BlockChain().SubscribeChainEvent(ch)
}

func (b *BfeAPIBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.bfe.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *BfeAPIBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.bfe.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *BfeAPIBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.bfe.BlockChain().SubscribeLogsEvent(ch)
}

func (b *BfeAPIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.bfe.txPool.AddLocal(signedTx)
}

func (b *BfeAPIBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.bfe.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *BfeAPIBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.bfe.txPool.Get(hash)
}

func (b *BfeAPIBackend) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, error) {
	tx, blockHash, blockNumber, index := rawdb.ReadTransaction(b.bfe.ChainDb(), txHash)
	return tx, blockHash, blockNumber, index, nil
}

func (b *BfeAPIBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.bfe.txPool.Nonce(addr), nil
}

func (b *BfeAPIBackend) Stats() (pending int, queued int) {
	return b.bfe.txPool.Stats()
}

func (b *BfeAPIBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.bfe.TxPool().Content()
}

func (b *BfeAPIBackend) TxPool() *core.TxPool {
	return b.bfe.TxPool()
}

func (b *BfeAPIBackend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.bfe.TxPool().SubscribeNewTxsEvent(ch)
}

func (b *BfeAPIBackend) Downloader() *downloader.Downloader {
	return b.bfe.Downloader()
}

func (b *BfeAPIBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *BfeAPIBackend) ChainDb() bfedb.Database {
	return b.bfe.ChainDb()
}

func (b *BfeAPIBackend) EventMux() *event.TypeMux {
	return b.bfe.EventMux()
}

func (b *BfeAPIBackend) AccountManager() *accounts.Manager {
	return b.bfe.AccountManager()
}

func (b *BfeAPIBackend) ExtRPCEnabled() bool {
	return b.extRPCEnabled
}

func (b *BfeAPIBackend) UnprotectedAllowed() bool {
	return b.allowUnprotectedTxs
}

func (b *BfeAPIBackend) RPCGasCap() uint64 {
	return b.bfe.config.RPCGasCap
}

func (b *BfeAPIBackend) RPCTxFeeCap() float64 {
	return b.bfe.config.RPCTxFeeCap
}

func (b *BfeAPIBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.bfe.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *BfeAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.bfe.bloomRequests)
	}
}

func (b *BfeAPIBackend) Engine() consensus.Engine {
	return b.bfe.engine
}

func (b *BfeAPIBackend) CurrentHeader() *types.Header {
	return b.bfe.blockchain.CurrentHeader()
}

func (b *BfeAPIBackend) Miner() *miner.Miner {
	return b.bfe.Miner()
}

func (b *BfeAPIBackend) StartMining(threads int) error {
	return b.bfe.StartMining(threads)
}

func (b *BfeAPIBackend) StateAtBlock(ctx context.Context, block *types.Block, reexec uint64) (*state.StateDB, func(), error) {
	return b.bfe.stateAtBlock(block, reexec)
}

func (b *BfeAPIBackend) StatesInRange(ctx context.Context, fromBlock *types.Block, toBlock *types.Block, reexec uint64) ([]*state.StateDB, func(), error) {
	return b.bfe.statesInRange(fromBlock, toBlock, reexec)
}

func (b *BfeAPIBackend) StateAtTransaction(ctx context.Context, block *types.Block, txIndex int, reexec uint64) (core.Message, vm.BlockContext, *state.StateDB, func(), error) {
	return b.bfe.stateAtTransaction(block, txIndex, reexec)
}
