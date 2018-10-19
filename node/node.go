// Copyright 2014 The Metabase Authors
// This file is part of vasuki p2p library.
//
// The vasuki library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The vasuki library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package node

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/bits"

	"github.com/Metabase-Network/vasuki/common"
	"github.com/Metabase-Network/vasuki/crypto"
	"github.com/Metabase-Network/vasuki/stor"
)

// Def is an identity of nodes, using its public key hash and network address.
type Def struct {
	NodeID     []byte
	NodeAddr   common.Address
	NodePubKey ecdsa.PublicKey
}

var errInvalidPubkey = errors.New("Error Generating Key for Node ")

// CreateNode Is a factory function which initializes Def
func CreateNode(path string) Def {
	var res *ecdsa.PrivateKey
	var err error
	// Gen Key -> Convert to HEX string -> Convert to
	db, err := stor.Start(path)
	key, _ := db.Get([]byte("NodePvk"))
	if err == nil {
		res1, _ := crypto.HexToECDSA(hex.EncodeToString(key))
		res = res1
	} else {
		key, _ := crypto.GenerateKey()
		key1 := crypto.FromECDSA(key)
		db.Put([]byte("NodePvk"), key1)
		res1, _ := crypto.HexToECDSA(hex.EncodeToString(key1))
		res = res1
	}
	if err != nil {
		return Def{}
	}
	return Def{NodePubKey: res.PublicKey, NodeAddr: crypto.PubkeyToAddress(res.PublicKey), NodeID: crypto.Keccak256(crypto.PubkeyToAddress(res.PublicKey).Bytes())}
}

//Equals Compares 2 Node ID
func (id Def) Equals(other []byte) bool {
	return bytes.Equal(id.NodeID, other)
}

//XorID XOR's ID
func (id Def) XorID(other []byte) []byte {
	result := make([]byte, len(id.NodeID))

	for i := 0; i < len(id.NodeID) && i < len(other); i++ {
		result[i] = id.NodeID[i] ^ other[i]
	}
	return result
}

//AddressHex Converts the address to hex String
func (id Def) AddressHex() string {
	return hex.EncodeToString(id.NodeAddr.Bytes())
}

//IDHex Converts the Node ID to Hex String
func (id Def) IDHex() string {
	return hex.EncodeToString(id.NodeID)
}

// PrefixLen returns the number of prefixed zeros in a peer Def.
func (id Def) PrefixLen() int {
	for i, b := range id.NodeID {
		if b != 0 {
			return i*8 + bits.LeadingZeros8(uint8(b))
		}
	}
	return len(id.NodeID)*8 - 1
}
