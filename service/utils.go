package service

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ontio/ontology/common"
	"github.com/polynetwork/poly-go-sdk"
	mcom "github.com/polynetwork/poly/common"
	"github.com/polynetwork/poly/native/service/cross_chain_manager/btc"
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
		config.SetBech32PrefixForAccount("swth", "swthpub")
		config.SetBech32PrefixForValidator("swthvaloper", "swthvaloperpub")
		config.SetBech32PrefixForConsensusNode("swthvalcons", "swthvalconspub")
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

func SetPolyChainId(poly *poly_go_sdk.PolySdk) error {
	hdr, err := poly.GetHeaderByHeight(0)
	if err != nil {
		return err
	}
	poly.SetChainId(hdr.ChainID)
	return nil
}
