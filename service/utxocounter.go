package service

import (
	"github.com/ontio/eth_tools/log"
	sdk "github.com/ontio/multi-chain-go-sdk"
	"github.com/ontio/multi-chain/common"
	"github.com/ontio/multi-chain/native/service/cross_chain_manager/btc"
	"github.com/ontio/multi-chain/native/service/utils"
)

func CountAlliaUtxo(rpcAddr string) {
	allia := sdk.NewMultiChainSdk()
	allia.NewRpcClient().SetAddress(rpcAddr)

	store, err := allia.GetStorage(utils.CrossChainManagerContractAddress.ToHexString(),
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

	log.Infof("sum of utxos in alliance is %d", sum)
}
