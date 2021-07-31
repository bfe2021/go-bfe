// Copyright 2019 The go-bfe Authors
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
	"github.com/bfe2021/go-bfe/core"
	"github.com/bfe2021/go-bfe/core/forkid"
	"github.com/bfe2021/go-bfe/p2p/dnsdisc"
	"github.com/bfe2021/go-bfe/p2p/enode"
	"github.com/bfe2021/go-bfe/rlp"
)

// ongEntry is the "ong" ENR entry which advertises ong protocol
// on the discovery network.
type ongEntry struct {
	ForkID forkid.ID // Fork identifier per EIP-2124

	// Ignore additional fields (for forward compatibility).
	Rest []rlp.RawValue `rlp:"tail"`
}

// ENRKey implements enr.Entry.
func (e ongEntry) ENRKey() string {
	return "ong"
}

// startBfeEntryUpdate starts the ENR updater loop.
func (ong *Bfedu) startBfeEntryUpdate(ln *enode.LocalNode) {
	var newHead = make(chan core.ChainHeadEvent, 10)
	sub := bfe.blockchain.SubscribeChainHeadEvent(newHead)

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case <-newHead:
				ln.Set(bfe.currentBfeEntry())
			case <-sub.Err():
				// Would be nice to sync with bfe.Stop, but there is no
				// good way to do that.
				return
			}
		}
	}()
}

func (ong *Bfedu) currentBfeEntry() *ongEntry {
	return &ongEntry{ForkID: forkid.NewID(bfe.blockchain.Config(), bfe.blockchain.Genesis().Hash(),
		bfe.blockchain.CurrentHeader().Number.Uint64())}
}

// setupDiscovery creates the node discovery source for the `ong` and `snap`
// protocols.
func setupDiscovery(urls []string) (enode.Iterator, error) {
	if len(urls) == 0 {
		return nil, nil
	}
	client := dnsdisc.NewClient(dnsdisc.Config{})
	return client.NewIterator(urls...)
}
