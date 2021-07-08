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

	"github.com/AlayaNetwork/Alaya-Go/log"

	"github.com/AlayaNetwork/Alaya-Go/params"

	"github.com/AlayaNetwork/Alaya-Go/x/reward"

	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/x/xcom"
)

//给没有领取委托奖励的账户平账 , https://github.com/PlatONnetwork/PlatON-Go/issues/1583
func NewFixIssue1583Plugin() *FixIssue1583Plugin {
	fix := new(FixIssue1583Plugin)
	return fix
}

type FixIssue1583Plugin struct{}

func (a *FixIssue1583Plugin) fix(blockHash common.Hash, chainID *big.Int, state xcom.StateDB) error {
	if chainID.Cmp(params.AlayaChainConfig.ChainID) != 0 {
		return nil
	}
	accounts, err := newIssue1583Accounts()
	if err != nil {
		return err
	}
	for _, account := range accounts {
		receiveReward := account.RewardPer.CalDelegateReward(account.delegationAmount)
		if err := rm.ReturnDelegateReward(account.addr, receiveReward, state); err != nil {
			log.Error("fix issue 1583,return delegate reward fail", "account", account.addr, "err", err)
			return common.InternalError
		}
	}
	return nil
}

type issue1583Accounts struct {
	addr             common.Address
	delegationAmount *big.Int
	RewardPer        reward.DelegateRewardPer
}

func newIssue1583Accounts() ([]issue1583Accounts, error) {
	type delegationInfo struct {
		account          string
		delegationAmount string
	}

	//node f2ec2830850  in Epoch216
	node1DelegationInfo := []delegationInfo{
		{"atp1a78ls0s4zrw66dwmx0h6vu6lp2wffjn52ftfkk", "1000000000000000000"},
		{"atp1tmdug9sv5kv06wu87yvhuedzat7mqwtk2runc9", "1000000000000000000"},
		{"atp1vtesk9c9eqq3aalajlnnvgxw7a86gq8dw423x9", "1000000000000000000"},
		{"atp14hmsk2qne5ps3epf9jx8cyz94mkq00y4jlzvl8", "1000000000000000000"},

	}
	node1DelegationAmount, _ := new(big.Int).SetString("5000000000000000000", 10)
	node1DelegationReward, _ := new(big.Int).SetString("2072192513368983957211", 10)
	node1RewardPer := reward.DelegateRewardPer{
		Delegate: node1DelegationAmount,
		Reward:   node1DelegationReward,
	}

	accounts := make([]issue1583Accounts, 0)
	for _, c := range node1DelegationInfo {
		addr, err := common.Bech32ToAddress(c.account)
		if err != nil {
			return nil, err
		}
		amount, _ := new(big.Int).SetString(c.delegationAmount, 10)
		accounts = append(accounts, issue1583Accounts{
			addr:             addr,
			delegationAmount: amount,
			RewardPer:        node1RewardPer,
		})
	}

	//fff1010bbf176 in epoch475
	node2DelegationInfos := []delegationInfo{
		{"atp1m0qy0ylh7evjk2yfrezkmaylwz32c4n3ttqukp", "1000000000000000000"},
		{"atp1nmy3lfpxk0r7rezr5hjdscqsulwjlfprmjundm", "1000000000000000000"},
		{"atp16a5f649vmwk3842jcwajrfqtjrsu5rlfxx8v2h", "1000000000000000000"},
		{"atp1xwew5ywq4lcy3tt4epay5dm8mx46cs5t389y7v", "1000000000000000000"},
		{"atp102rkam0m4ayasfe5awklaf3y367zedr6vc236c", "1000000000000000000"},
		{"atp19nzzh4p3s87jf436hs2wmjjwta7k5yak2p98az", "1000000000000000000"},

	}

	node2DelegationAmount, _ := new(big.Int).SetString("7000000000000000000", 10)
	node2DelegationReward, _ := new(big.Int).SetString("1737967914438502673791", 10)
	node2RewardPer := reward.DelegateRewardPer{
		Delegate: node2DelegationAmount,
		Reward:   node2DelegationReward,
	}

	for _, c := range node2DelegationInfos {
		addr, err := common.Bech32ToAddress(c.account)
		if err != nil {
			return nil, err
		}
		amount, _ := new(big.Int).SetString(c.delegationAmount, 10)
		accounts = append(accounts, issue1583Accounts{
			addr:             addr,
			delegationAmount: amount,
			RewardPer:        node2RewardPer,
		})
	}

	return accounts, nil
}
