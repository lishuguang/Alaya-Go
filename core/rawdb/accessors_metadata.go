// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rawdb

import (
	"encoding/json"

	"github.com/AlayaNetwork/Alaya-Go/ethdb"

	"github.com/AlayaNetwork/Alaya-Go/x/xcom"

	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/log"
	"github.com/AlayaNetwork/Alaya-Go/params"
	"github.com/AlayaNetwork/Alaya-Go/rlp"
)

// ReadDatabaseVersion retrieves the version number of the database.
func ReadDatabaseVersion(db ethdb.KeyValueReader) int {
    var version uint64

    enc, _ := db.Get(databaseVerisionKey)
    rlp.DecodeBytes(enc, &version)

    return int(version)
}

// WriteDatabaseVersion stores the version number of the database
func WriteDatabaseVersion(db ethdb.KeyValueWriter, version int) {
    enc, _ := rlp.EncodeToBytes(uint64(version))
    if err := db.Put(databaseVerisionKey, enc); err != nil {
        log.Crit("Failed to store the database version", "err", err)
    }
}

// ReadChainConfig retrieves the consensus settings based on the given genesis hash.
func ReadChainConfig(db ethdb.KeyValueReader, hash common.Hash) *params.ChainConfig {
	data, _ := db.Get(configKey(hash))
	if len(data) == 0 {
		return nil
	}
	var config params.ChainConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Error("Invalid chain config JSON", "hash", hash, "err", err)
		return nil
	}
	return &config
}

// WriteChainConfig writes the chain config settings to the database.
func WriteChainConfig(db ethdb.KeyValueWriter, hash common.Hash, cfg *params.ChainConfig) {
	if cfg == nil {
		return
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		log.Crit("Failed to JSON encode chain config", "err", err)
	}
	if err := db.Put(configKey(hash), data); err != nil {
		log.Crit("Failed to store chain config", "err", err)
	}
}

// WriteEconomicModel writes the EconomicModel settings to the database.
func WriteEconomicModel(db ethdb.Writer, hash common.Hash, ec *xcom.EconomicModel) {
	if ec == nil {
		return
	}

	data, err := json.Marshal(ec)
	if err != nil {
		log.Crit("Failed to JSON encode EconomicModel config", "err", err)
	}
	if err := db.Put(economicModelKey(hash), data); err != nil {
		log.Crit("Failed to store EconomicModel", "err", err)
	}
}

// WriteEconomicModelExtend writes the EconomicModelExtend settings to the database.
func WriteEconomicModelExtend(db ethdb.Writer, hash common.Hash, ec *xcom.EconomicModelExtend) {
	if ec == nil {
		return
	}

	data, err := json.Marshal(ec)
	if err != nil {
		log.Crit("Failed to JSON encode EconomicModelExtend config", "err", err)
	}
	if err := db.Put(economicModelExtendKey(hash), data); err != nil {
		log.Crit("Failed to store EconomicModelExtend", "err", err)
	}
}

// ReadEconomicModel retrieves the EconomicModel settings based on the given genesis hash.
func ReadEconomicModel(db ethdb.Reader, hash common.Hash) *xcom.EconomicModel {
	data, _ := db.Get(economicModelKey(hash))
	if len(data) == 0 {
		return nil
	}

	var ec xcom.EconomicModel
	// reset the global ec
	if err := json.Unmarshal(data, &ec); err != nil {
		log.Error("Invalid EconomicModel JSON", "hash", hash, "err", err)
		return nil
	}
	return &ec
}

// ReadEconomicModelExtend retrieves the EconomicModelExtend settings based on the given genesis hash.
func ReadEconomicModelExtend(db ethdb.Reader, hash common.Hash) *xcom.EconomicModelExtend {
	data, _ := db.Get(economicModelExtendKey(hash))
	if len(data) == 0 {
		return nil
	}

	var ec xcom.EconomicModelExtend
	// reset the global ec
	if err := json.Unmarshal(data, &ec); err != nil {
		log.Error("Invalid EconomicModelExtend JSON", "hash", hash, "err", err)
		return nil
	}
	return &ec
}

// ReadPreimage retrieves a single preimage of the provided hash.
func ReadPreimage(db ethdb.KeyValueReader, hash common.Hash) []byte {
	data, _ := db.Get(preimageKey(hash))
	return data
}

// WritePreimages writes the provided set of preimages to the database.
func WritePreimages(db ethdb.KeyValueWriter, preimages map[common.Hash][]byte) {
	for hash, preimage := range preimages {
		if err := db.Put(preimageKey(hash), preimage); err != nil {
			log.Crit("Failed to store trie preimage", "err", err)
		}
	}
	preimageCounter.Inc(int64(len(preimages)))
	preimageHitCounter.Inc(int64(len(preimages)))
}
