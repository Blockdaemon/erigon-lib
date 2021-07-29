/*
   Copyright 2021 Erigon contributors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package txpool

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubPoolMarkerOrder(t *testing.T) {
	require := require.New(t)
	require.Less(
		NewSubPoolMarker(true, true, true, true, false),
		NewSubPoolMarker(true, true, true, true, true),
	)
	require.Less(
		NewSubPoolMarker(true, true, true, false, true),
		NewSubPoolMarker(true, true, true, true, true),
	)
	require.Less(
		NewSubPoolMarker(true, true, true, false, true),
		NewSubPoolMarker(true, true, true, true, false),
	)
	require.Less(
		NewSubPoolMarker(false, true, true, true, true),
		NewSubPoolMarker(true, false, true, true, true),
	)
	require.Less(
		NewSubPoolMarker(false, false, false, true, true),
		NewSubPoolMarker(false, false, true, true, true),
	)
	require.Less(
		NewSubPoolMarker(false, false, true, true, false),
		NewSubPoolMarker(false, false, true, true, true),
	)
}

func TestSubPool(t *testing.T) {
	sub := NewSubPool()
	sub.Add(&MetaTx{SubPool: 0b10101})
	sub.Add(&MetaTx{SubPool: 0b11110})
	sub.Add(&MetaTx{SubPool: 0b11101})
	sub.Add(&MetaTx{SubPool: 0b10001})
	require.Equal(t, uint8(0b11110), uint8(sub.Best().SubPool))
	require.Equal(t, uint8(0b10001), uint8(sub.Worst().SubPool))

	sub = NewSubPool()
	sub.Add(&MetaTx{SubPool: 0b00001})
	sub.Add(&MetaTx{SubPool: 0b01110})
	sub.Add(&MetaTx{SubPool: 0b01101})
	sub.Add(&MetaTx{SubPool: 0b00101})
	require.Equal(t, uint8(0b00001), uint8(sub.Worst().SubPool))
	require.Equal(t, uint8(0b01110), uint8(sub.Best().SubPool))
}