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

// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/bfe2021/go-bfe/metrics"
)

var (
	headerInMeter      = metrics.NewRegisteredMeter("bfe/downloader/headers/in", nil)
	headerReqTimer     = metrics.NewRegisteredTimer("bfe/downloader/headers/req", nil)
	headerDropMeter    = metrics.NewRegisteredMeter("bfe/downloader/headers/drop", nil)
	headerTimeoutMeter = metrics.NewRegisteredMeter("bfe/downloader/headers/timeout", nil)

	bodyInMeter      = metrics.NewRegisteredMeter("bfe/downloader/bodies/in", nil)
	bodyReqTimer     = metrics.NewRegisteredTimer("bfe/downloader/bodies/req", nil)
	bodyDropMeter    = metrics.NewRegisteredMeter("bfe/downloader/bodies/drop", nil)
	bodyTimeoutMeter = metrics.NewRegisteredMeter("bfe/downloader/bodies/timeout", nil)

	receiptInMeter      = metrics.NewRegisteredMeter("bfe/downloader/receipts/in", nil)
	receiptReqTimer     = metrics.NewRegisteredTimer("bfe/downloader/receipts/req", nil)
	receiptDropMeter    = metrics.NewRegisteredMeter("bfe/downloader/receipts/drop", nil)
	receiptTimeoutMeter = metrics.NewRegisteredMeter("bfe/downloader/receipts/timeout", nil)

	stateInMeter   = metrics.NewRegisteredMeter("bfe/downloader/states/in", nil)
	stateDropMeter = metrics.NewRegisteredMeter("bfe/downloader/states/drop", nil)

	throttleCounter = metrics.NewRegisteredCounter("bfe/downloader/throttle", nil)
)
