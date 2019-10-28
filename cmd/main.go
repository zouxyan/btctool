package main

import (
	"flag"
	"fmt"
	"github.com/Zou-XueYan/btctool/service"
	"github.com/ontio/multi-chain/common/log"
	"os"
	"strconv"
	"strings"
)

var prevHexTxs string
var indexes string
var privkb58 string
var addr string
var value float64
var fee float64
var ontAddr string
var tsec uint
var spvAddr string
var rpcUrl string
var user string
var pwd string
var defaultAddr string
var tool string
var netType string
var utxoVals string
var txids string
var contractAddr string
var dura int64
var maxVal int64

var addrScriptHash string = "2N5cY8y9RtbbvQRWkX5zAwTPCxSZF9xEj2C"

func init() {
	flag.StringVar(&tool, "tool", "", "which tool to use, \"regauto\", \"cctx\" or \"blkgene\"")
	flag.StringVar(&prevHexTxs, "hextxs", "", "Raw transaction in hex")
	flag.StringVar(&indexes, "idxes", "", "Your UTXO's index for building this transaction")
	flag.StringVar(&privkb58, "privkb58", "", "Your private key in base58")
	flag.StringVar(&addr, "addr", "", "Your btc address in base58")
	flag.Float64Var(&value, "value", 0.0001, "Amount of btc to cross chain")
	flag.Float64Var(&fee, "fee", 0.00001, "Cross chain fee, not in use for now")
	flag.StringVar(&ontAddr, "targetaddr", "", "Your target chain address")
	flag.UintVar(&tsec, "tsec", 5, "block update time(seconds), 5 sec default")
	flag.StringVar(&rpcUrl, "url", "", "the bitcoind rpc address")
	flag.StringVar(&spvAddr, "spvaddr", "", "the spv client address for broadcasting tx")
	flag.StringVar(&user, "user", "", "the rpc user")
	flag.StringVar(&pwd, "pwd", "", "the rpc password")
	flag.StringVar(&defaultAddr, "defaultaddr", "", "the default bitcoin address to rececive the mining reward")
	flag.StringVar(&netType, "net", "test", "the net type of bitcoin")
	flag.StringVar(&utxoVals, "utxovals", "", "val of utxos in satoshi, eg: 1000,2000")
	flag.StringVar(&txids, "txids", "", "txid of utxos, eg: xx,yy,bb")
	flag.StringVar(&contractAddr, "contract", "", "target chain smart contract address")
	flag.Int64Var(&dura, "dura", 300, "set the seconds to send a cross-tx, default 5 min")
	flag.Int64Var(&maxVal, "maxval", 2000, "the max value of cross tx")
}

func main() {
	flag.Parse()

	switch tool {
	case "regauto":
		handler := &service.RegAuto{
			RpcUrl:         rpcUrl,
			Privkb58:       privkb58,
			Addr:           addr,
			Fee:            fee,
			Value:          value,
			AddrScriptHash: addrScriptHash,
			OntAddr:        ontAddr,
			Pwd:            pwd,
			User:           user,
			ContractAddr:   contractAddr,
		}
		handler.RunRegAuto()
	case "cctx":
		valArr, err := getVals(utxoVals)
		if err != nil {
			log.Errorf("failed to get vals: %v", err)
			os.Exit(1)
		}

		handler := service.CcTx{
			OntAddr:        ontAddr,
			AddrScriptHash: addrScriptHash,
			Value:          value,
			Fee:            fee,
			Privkb58:       privkb58,
			Indexes:        indexes,
			SpvAddr:        spvAddr,
			NetType:        netType,
			Vals:           valArr,
			Txids:          txids,
			ContractAddr:   contractAddr,
		}
		handler.RunCcTx()
	case "blkgene":
		handler := service.BlkGene{
			User:        user,
			Pwd:         pwd,
			DefaultAddr: defaultAddr,
			Tsec:        tsec,
			RpcUrl:      rpcUrl,
		}
		handler.RunBlkGene()
	case "autosender":
		valArr, err := getVals(utxoVals)
		if err != nil {
			log.Errorf("failed to get vals: %v", err)
			os.Exit(1)
		}

		handler := service.AutoSender{
			CcTx: &service.CcTx{
				OntAddr:        ontAddr,
				AddrScriptHash: addrScriptHash,
				Value:          value,
				Fee:            fee,
				Privkb58:       privkb58,
				Indexes:        indexes,
				SpvAddr:        spvAddr,
				NetType:        netType,
				Vals:           valArr,
				Txids:          txids,
				ContractAddr:   contractAddr,
			},
			MaxVal: maxVal,
			Dura:   dura,
		}
		handler.Sending()
	default:
		log.Errorf("no handler matched")
		os.Exit(1)
	}
}

func getVals(val string) ([]float64, error) {
	var valArr []float64
	vals := strings.Split(val, ",")
	for i, val := range vals {
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse val %s: %v", val, err)
		}
		if num <= 0 {
			return nil, fmt.Errorf("no.%d value %d can not less than or equal to zero", i, val)
		}
		valArr = append(valArr, num)
	}

	return valArr, nil
}
