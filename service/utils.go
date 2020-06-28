package service

import (
	"encoding/hex"
	"fmt"
	"github.com/ConjurTech/switcheo-chain/cmd"
	"github.com/cosmos/cosmos-sdk/types"
	mcom "github.com/ontio/multi-chain/common"
	"github.com/ontio/multi-chain/native/service/cross_chain_manager/btc"
	"github.com/ontio/ontology/common"
	"strings"
)

func buildData(toChainId uint64, ccFee int64, toAddr string) ([]byte, error) {
	var data []byte
	ccflag := byte(0xcc)
	var args *btc.Args
	switch toChainId {
	case 2:
		toAddr = strings.ReplaceAll(toAddr, "0x", "")
		toAddrBytes, _ := hex.DecodeString(toAddr)
		args = &btc.Args{
			Address:   toAddrBytes[:],
			ToChainID: toChainId,
			Fee:       ccFee,
		}
	case 3:
		addrBytes, _ := common.AddressFromBase58(toAddr)
		args = &btc.Args{
			Address:   addrBytes[:],
			ToChainID: toChainId,
			Fee:       ccFee,
		}
	case 4:
		addrBytes, _ := hex.DecodeString(toAddr)
		args = &btc.Args{
			Address:   addrBytes[:],
			ToChainID: toChainId,
			Fee:       ccFee,
		}
	case 5:
		config := types.GetConfig()
		config.SetBech32PrefixForAccount(types.Bech32PrefixAccAddr, types.Bech32PrefixAccPub)
		config.SetBech32PrefixForValidator(types.Bech32PrefixValAddr, types.Bech32PrefixValPub)
		config.SetBech32PrefixForConsensusNode(types.Bech32PrefixConsAddr, types.Bech32PrefixConsPub)
		addr, _ := types.AccAddressFromBech32(toAddr)
		args = &btc.Args{
			Address:   addr[:],
			ToChainID: toChainId,
			Fee:       ccFee,
		}
	case 172:
		config := types.GetConfig()
		config.SetBech32PrefixForAccount(cmd.MainPrefix, cmd.MainPrefix+types.PrefixPublic)
		config.SetBech32PrefixForValidator(cmd.MainPrefix+types.PrefixValidator+types.PrefixOperator, cmd.MainPrefix+types.PrefixValidator+types.PrefixOperator+types.PrefixPublic)
		config.SetBech32PrefixForConsensusNode(cmd.MainPrefix+types.PrefixValidator+types.PrefixConsensus, cmd.MainPrefix+types.PrefixValidator+types.PrefixConsensus+types.PrefixPublic)
		addr, _ := types.AccAddressFromBech32(toAddr)
		args = &btc.Args{
			Address:   addr[:],
			ToChainID: toChainId,
			Fee:       ccFee,
		}
	default:
		return nil, fmt.Errorf("not supported chainid %d", toChainId)
	}
	var buf []byte
	sink := mcom.NewZeroCopySink(buf)
	args.Serialization(sink)
	data = append(append(data, ccflag), sink.Bytes()...)

	return data, nil
}

type Service interface {
	Run()
}
