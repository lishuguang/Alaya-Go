// Copyright 2021 The Alaya Network Authors
// This file is part of the Alaya-Go library.
//
// The Alaya-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Alaya-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Alaya-Go library. If not, see <http://www.gnu.org/licenses/>.

package plugin

import (
	"math/big"

	"github.com/AlayaNetwork/Alaya-Go/x/xcom"

	"github.com/AlayaNetwork/Alaya-Go/log"

	"github.com/AlayaNetwork/Alaya-Go/x/xutil"

	"github.com/AlayaNetwork/Alaya-Go/p2p/discover"

	"github.com/AlayaNetwork/Alaya-Go/params"

	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/core/snapshotdb"
)

//this is use fix validators staking shares error, https://github.com/PlatONnetwork/PlatON-Go/issues/1654
func NewFixIssue1654Plugin(sdb snapshotdb.DB) *FixIssue1654Plugin {
	fix := new(FixIssue1654Plugin)
	fix.sdb = sdb
	return fix
}

type FixIssue1654Plugin struct {
	sdb snapshotdb.DB
}

func (a *FixIssue1654Plugin) fix(blockHash common.Hash, chainID *big.Int, state xcom.StateDB) error {
	if chainID.Cmp(params.AlayaChainConfig.ChainID) != 0 {
		return nil
	}
	candidates, err := NewIssue1654Candidates()
	if err != nil {
		return err
	}
	for _, candidate := range candidates {
		canAddr, err := xutil.NodeId2Addr(candidate.nodeID)
		if nil != err {
			return err
		}
		can, err := stk.GetCandidateInfo(blockHash, canAddr)
		if snapshotdb.NonDbNotFoundErr(err) {
			return err
		}
		if can.IsNotEmpty() && can.StakingBlockNum == candidate.stakingNum {
			if can.Status.IsValid() {
				if err := stk.db.DelCanPowerStore(blockHash, can); nil != err {
					return err
				}
				can.SubShares(candidate.shouldSub)
				if err := stk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
					return err
				}
				if err := stk.db.SetCanMutableStore(blockHash, canAddr, can.CandidateMutable); nil != err {
					return err
				}
				log.Debug("fix issue1654,can is valid,update the can power", "nodeID", candidate.nodeID, "stakingNum", candidate.stakingNum, "sub", candidate.shouldSub, "newShare", can.Shares)
			} else {
				if can.Shares != nil {
					if can.Shares.Cmp(candidate.shouldSub)>=0{
						can.SubShares(candidate.shouldSub)
						if err := stk.db.SetCanMutableStore(blockHash, canAddr, can.CandidateMutable); nil != err {
							return err
						}
						log.Debug("fix issue1654,can is invalid", "nodeID", candidate.nodeID, "stakingNum", candidate.stakingNum, "sub", candidate.shouldSub, "newShare", can.Shares)
					}
				}
			}
		}
	}
	return nil
}

type issue1654Candidate struct {
	nodeID     discover.NodeID
	stakingNum uint64
	shouldSub  *big.Int
}

func NewIssue1654Candidates() ([]issue1654Candidate, error) {
	type candidate struct {
		Node   string
		Num    int
		Amount string
	}

	candidates := []candidate{
		{"493c66bd7d6051e42a68bffa5f70005555886f28a0d9f10afaca4abc45723a26d6b833126fb65f11e3be51613405df664e7cda12baad538dd08b0a5774aa22cf", 25, "798328877005347593582890"},
	}

	nodes := make([]issue1654Candidate, 0)
	for _, c := range candidates {
		amount, _ := new(big.Int).SetString(c.Amount, 10)
		nodes = append(nodes, issue1654Candidate{
			nodeID:     discover.MustHexID(c.Node),
			stakingNum: uint64(c.Num),
			shouldSub:  amount,
		})
	}
	return nodes, nil
}
