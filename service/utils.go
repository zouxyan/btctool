package service

import (
	"encoding/hex"
	mcom "github.com/ontio/multi-chain/common"
	"github.com/ontio/multi-chain/native/service/cross_chain_manager/btc"
	"github.com/ontio/ontology/common"

	"strings"
)

func buildData(toChainId uint64, ccFee int64, toAddr, contractAddr string) ([]byte, error) {
	var data []byte
	ccflag := byte(0x66)
	var args *btc.Args
	switch toChainId {
	case 1:
		toAddr = strings.ReplaceAll(toAddr, "0x", "")
		contractAddr = strings.ReplaceAll(contractAddr, "0x", "")

		toAddrBytes, _ := hex.DecodeString(toAddr)
		contract, _ := hex.DecodeString(contractAddr)
		args = &btc.Args{
			Address:           toAddrBytes[:],
			ToChainID:         toChainId,
			Fee:               ccFee,
			ToContractAddress: contract[:],
		}
	case 2:
		contractAddrBytes, err := common.AddressFromHexString(contractAddr)
		if err != nil {
			return nil, err
		}
		addrBytes, _ := common.AddressFromBase58(toAddr)
		args = &btc.Args{
			Address:           addrBytes[:],
			ToChainID:         toChainId,
			Fee:               ccFee,
			ToContractAddress: contractAddrBytes[:],
		}
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
