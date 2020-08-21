package service

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	sdk "github.com/polynetwork/poly-go-sdk"
	"github.com/polynetwork/poly/common"
	"github.com/polynetwork/poly/common/log"
	"github.com/polynetwork/poly/native/service/cross_chain_manager/btc"
	"github.com/polynetwork/poly/native/service/governance/side_chain_manager"
	"github.com/polynetwork/poly/native/service/utils"
	"io/ioutil"
	"sort"
	"time"
)

type UtxoStatus struct {
	Sum        uint64
	Total      uint64
	P2shNum    uint64
	P2shSum    uint64
	P2wshNum   uint64
	P2wshSum   uint64
	FeeRate    uint64
	MinChange  uint64
	Less       uint64
	UtxosInStr string
}

type UtxoMonitor struct {
	Status    *UtxoStatus
	poly      *sdk.PolySdk
	rk        []byte
	lessPoint uint64
	quit      chan struct{}
}

func NewUtxoMonitor(lp uint64, rpcAddr string, redeem []byte) *UtxoMonitor {
	poly := sdk.NewPolySdk()
	poly.NewRpcClient().SetAddress(rpcAddr)
	if err := SetPolyChainId(poly); err != nil {
		log.Fatalf("failed to set poly chain id: %v", err)
		return nil
	}

	k := btcutil.Hash160(redeem)
	return &UtxoMonitor{
		Status:    &UtxoStatus{},
		poly:      poly,
		rk:        k,
		lessPoint: lp,
		quit:      make(chan struct{}),
	}
}

func (m *UtxoMonitor) RunMonitor() {
	tick := time.NewTicker(time.Second * 5)
	utxos := &btc.Utxos{
		Utxos: make([]*btc.Utxo, 0),
	}
	for {
		select {
		case <-tick.C:
			store, err := m.poly.GetStorage(utils.CrossChainManagerContractAddress.ToHexString(),
				append(append([]byte(btc.UTXOS), utils.GetUint64Bytes(1)...), []byte(hex.EncodeToString(m.rk))...))
			if err != nil {
				log.Errorf("failed to get utxos from chain: %v", err)
				continue
			}
			if err = utxos.Deserialization(common.NewZeroCopySource(store)); err != nil {
				log.Errorf("failed to deserialize utxos: %v", err)
				continue
			}
			sort.Sort(utxos)
			content := "Here is all utxo of your multisig-address\n"
			var sum, p2shCount, p2wshCount, lpCount, p2shSum, p2wshSum uint64
			for i, v := range utxos.Utxos {
				sum += v.Value
				cls := txscript.GetScriptClass(v.ScriptPubkey)
				switch cls {
				case txscript.ScriptHashTy:
					p2shCount++
					p2shSum += v.Value
				case txscript.WitnessV0ScriptHashTy:
					p2wshCount++
					p2wshSum += v.Value
				}
				if v.Value <= m.lessPoint {
					lpCount++
				}
				content += fmt.Sprintf("No.%d (outpoint: %s, value: %d, script_type: %s)\n", i, v.Op.String(), v.Value,
					cls.String())
			}
			store, err = m.poly.GetStorage(utils.SideChainManagerContractAddress.ToHexString(),
				append(append([]byte(side_chain_manager.BTC_TX_PARAM), m.rk...), utils.GetUint64Bytes(1)...))
			if err != nil {
				log.Errorf("failed to get btc tx param from chain: %v", err)
				continue
			}
			detial := &side_chain_manager.BtcTxParamDetial{}
			if store == nil {
				if err = detial.Deserialization(common.NewZeroCopySource(store)); err != nil {
					log.Errorf("deserialize BtcTxParamDetial error: %v", err)
				}
			}

			m.Status.UtxosInStr = content
			m.Status.Sum = sum
			m.Status.Total = uint64(len(utxos.Utxos))
			m.Status.Less = lpCount
			m.Status.FeeRate = detial.FeeRate
			m.Status.MinChange = detial.MinChange
			m.Status.P2shNum = p2shCount
			m.Status.P2wshNum = p2wshCount
			m.Status.P2shSum = p2shSum
			m.Status.P2wshSum = p2wshSum

			if err = ioutil.WriteFile("your_utxo", []byte(m.Status.UtxosInStr), 0644); err != nil {
				log.Errorf("failed to write utxo into file: %v", err)
			}

			log.Infof("status: (sum: %d, total: %d, less_count: %d, p2sh_utxo: %d, p2wsh_count: %d, "+
				"fee_rate: %d, min_change: %d)\n%s", m.Status.Sum, m.Status.Total, m.Status.Less, m.Status.P2shNum,
				m.Status.P2wshNum, m.Status.FeeRate, m.Status.MinChange, m.Status.UtxosInStr)
		case <-m.quit:
			log.Info("run monitor stopping")
			return
		}
	}
}

func (m *UtxoMonitor) Close() {
	close(m.quit)
}
