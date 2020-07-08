package service

import (
	"github.com/btcsuite/btcutil"
	"github.com/polynetwork/poly/common/log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type AutoSender struct {
	CcTx    *TestTxBuilder
	RegA    *RegTxBuilder
	Dura    int64
	MaxVal  int64
	NetType string
}

func (sender *AutoSender) sendingNotReg() {
	if sender.CcTx == nil {
		log.Error("CcTx service need to be init")
		os.Exit(1)
	}
	tick := time.NewTicker(time.Duration(sender.Dura) * time.Second)
	for {
		select {
		case t := <-tick.C:
			log.Infof("\ntry to build a tx (%s)", t.String())
			ptx := sender.CcTx.Run()
			sender.CcTx.Txids = ptx.TxHash().String()
			nextVal := rand.Int63n(sender.MaxVal-500) + 500
			if len(ptx.TxOut) > 2 && ptx.TxOut[2].Value < nextVal {
				log.Errorf("value not enough to pay: next is %d but only %d left", nextVal, ptx.TxOut[2].Value)
				os.Exit(1)
			}
			sender.CcTx.Value = float64(nextVal) / btcutil.SatoshiPerBitcoin
			sender.CcTx.Vals = []float64{float64(ptx.TxOut[2].Value) / btcutil.SatoshiPerBitcoin}
			sender.CcTx.Indexes = strconv.Itoa(2)
		}
	}
}

func (sender *AutoSender) sendingReg() {
	if sender.RegA == nil {
		log.Error("RegAuto service need to be init")
		os.Exit(1)
	}
	tick := time.NewTicker(time.Duration(sender.Dura) * time.Second)
	for {
		select {
		case t := <-tick.C:
			log.Infof("\ntry to build a tx (%s)", t.String())
			sender.RegA.Run()
			nextVal := rand.Int63n(sender.MaxVal-500) + 500
			sender.RegA.Value = float64(nextVal) / btcutil.SatoshiPerBitcoin
		}
	}
}

func (sender *AutoSender) Run() {
	switch sender.NetType {
	case "test":
		sender.sendingNotReg()
	case "main":
		sender.sendingNotReg()
	case "reg":
		sender.sendingReg()
	default:
		log.Errorf("WTF you input: %s", sender.NetType)
	}
}
