package service

import (
	"github.com/ontio/eth_tools/log"
	sdk "github.com/polynetwork/poly-go-sdk"
	"github.com/polynetwork/poly/common"
	"github.com/polynetwork/poly/native/service/cross_chain_manager/btc"
	"github.com/polynetwork/poly/native/service/utils"
)

func CountPolyUtxo(rpcAddr string) {
	poly := sdk.NewPolySdk()
	poly.NewRpcClient().SetAddress(rpcAddr)
	if err := SetPolyChainId(poly); err != nil {
		log.Errorf("failed to set poly chain id: %v", err)
		return
	}

	store, err := poly.GetStorage(utils.CrossChainManagerContractAddress.ToHexString(),
		append([]byte(btc.UTXOS), utils.GetUint64Bytes(0)...))
	if err != nil {
		log.Errorf("failed to get storage: %v", err)
		return
	}

	utxos := new(btc.Utxos)
	err = utxos.Deserialization(common.NewZeroCopySource(store))
	if err != nil {
		log.Errorf("failed to deserialize: %v", err)
		return
	}

	sum := uint64(0)
	for _, u := range utxos.Utxos {
		sum += u.Value
	}

	log.Infof("sum of utxos in poly is %d", sum)
}
