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
	"fmt"

	"github.com/ledgerwatch/erigon-lib/rlp"
)

type NewPooledTransactionHashesPacket [][32]byte

// ParseHashesCount looks at the RLP length Prefix for list of 32-byte hashes
// and returns number of hashes in the list to expect
func ParseHashesCount(payload Hashes, pos int) (count int, dataPos int, err error) {
	dataPos, dataLen, err := rlp.List(payload, pos)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: hashes len: %w", rlp.ParseHashErrorPrefix, err)
	}
	if dataLen%33 != 0 {
		return 0, 0, fmt.Errorf("%s: hashes len must be multiple of 33", rlp.ParseHashErrorPrefix)
	}
	return dataLen / 33, dataPos, nil
}

// EncodeHashes produces RLP encoding of given number of hashes, as RLP list
// It appends encoding to the given given slice (encodeBuf), reusing the space
// there is there is enough capacity.
// The first returned value is the slice where encodinfg
func EncodeHashes(hashes []byte, encodeBuf []byte) []byte {
	hashesLen := len(hashes) / 32 * 33
	dataLen := hashesLen
	encodeBuf = ensureEnoughSize(encodeBuf, rlp.ListPrefixLen(hashesLen)+dataLen)
	rlp.EncodeHashes(hashes, encodeBuf)
	return encodeBuf
}

// ParseHash extracts the next hash from the RLP encoding (payload) from a given position.
// It appends the hash to the given slice, reusing the space if there is enough capacity
// The first returned value is the slice where hash is appended to.
// The second returned value is the new position in the RLP payload after the extraction
// of the hash.
func ParseHash(payload []byte, pos int, hashbuf []byte) ([]byte, int, error) {
	hashbuf = ensureEnoughSize(hashbuf, 32)
	pos, err := rlp.ParseHash(payload, pos, hashbuf)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: hash len: %w", rlp.ParseHashErrorPrefix, err)
	}
	return hashbuf, pos, nil
}

func ensureEnoughSize(in []byte, size int) []byte {
	if cap(in) < size {
		newBuf := make([]byte, size)
		copy(newBuf, in)
		return newBuf
	}
	return in[:size] // Reuse the space if it has enough capacity
}

// EncodeGetPooledTransactions66 produces encoding of GetPooledTransactions66 packet
func EncodeGetPooledTransactions66(hashes []byte, requestId uint64, encodeBuf []byte) ([]byte, error) {
	pos := 0
	hashesLen := len(hashes) / 32 * 33
	dataLen := rlp.ListPrefixLen(hashesLen) + hashesLen + rlp.U64Len(requestId)
	encodeBuf = ensureEnoughSize(encodeBuf, rlp.ListPrefixLen(dataLen)+dataLen)
	// Length Prefix for the entire structure
	pos += rlp.EncodeListPrefix(dataLen, encodeBuf[pos:])
	pos += rlp.EncodeU64(requestId, encodeBuf[pos:])
	pos += rlp.EncodeHashes(hashes, encodeBuf[pos:])
	_ = pos
	return encodeBuf, nil
}

func ParseGetPooledTransactions66(payload []byte, pos int, hashbuf []byte) (requestID uint64, hashes []byte, newPos int, err error) {
	pos, _, err = rlp.List(payload, pos)
	if err != nil {
		return 0, hashes, 0, err
	}

	pos, requestID, err = rlp.U64(payload, pos)
	if err != nil {
		return 0, hashes, 0, err
	}
	var hashesCount int
	hashesCount, pos, err = ParseHashesCount(payload, pos)
	if err != nil {
		return 0, hashes, 0, err
	}
	hashes = ensureEnoughSize(hashbuf, 32*hashesCount)

	for i := 0; pos != len(payload); i++ {
		pos, err = rlp.ParseHash(payload, pos, hashes[i*32:])
		if err != nil {
			return 0, hashes, 0, err
		}
	}
	return requestID, hashes, pos, nil
}

func ParseGetPooledTransactions65(payload []byte, pos int, hashbuf []byte) (hashes []byte, newPos int, err error) {
	pos, _, err = rlp.List(payload, pos)
	if err != nil {
		return hashes, 0, err
	}

	var hashesCount int
	hashesCount, pos, err = ParseHashesCount(payload, pos)
	if err != nil {
		return hashes, 0, err
	}
	hashes = ensureEnoughSize(hashbuf, 32*hashesCount)

	for i := 0; pos != len(payload); i++ {
		pos, err = rlp.ParseHash(payload, pos, hashes[i*32:])
		if err != nil {
			return hashes, 0, err
		}
	}
	return hashes, pos, nil
}
func EncodePooledTransactions66(txsRlp [][]byte, requestId uint64, encodeBuf []byte) []byte {
	pos := 0
	txsRlpLen := 0
	for i := range txsRlp {
		txsRlpLen += len(txsRlp[i])
	}
	dataLen := rlp.U64Len(requestId) + rlp.ListPrefixLen(txsRlpLen) + txsRlpLen

	encodeBuf = ensureEnoughSize(encodeBuf, rlp.ListPrefixLen(dataLen)+dataLen)

	// Length Prefix for the entire structure
	pos += rlp.EncodeListPrefix(dataLen, encodeBuf[pos:])
	pos += rlp.EncodeU64(requestId, encodeBuf[pos:])
	pos += rlp.EncodeListPrefix(txsRlpLen, encodeBuf[pos:])
	for i := range txsRlp {
		copy(encodeBuf[pos:], txsRlp[i])
		pos += len(txsRlp[i])
	}
	_ = pos
	return encodeBuf
}
func EncodePooledTransactions65(txsRlp [][]byte, encodeBuf []byte) []byte {
	pos := 0
	dataLen := 0
	for i := range txsRlp {
		dataLen += len(txsRlp[i])
	}

	encodeBuf = ensureEnoughSize(encodeBuf, rlp.ListPrefixLen(dataLen)+dataLen)
	// Length Prefix for the entire structure
	pos += rlp.EncodeListPrefix(dataLen, encodeBuf[pos:])
	for i := range txsRlp {
		copy(encodeBuf[pos:], txsRlp[i])
		pos += len(txsRlp[i])
	}
	_ = pos
	return encodeBuf
}
