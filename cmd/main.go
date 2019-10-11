package main

import (
	"flag"
	"github.com/Zou-XueYan/btctool"
	"github.com/ontio/multi-chain/common/log"
	"os"
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
}

func main() {
	flag.Parse()

	switch tool {
	case "regauto":
		handler := &btctool.RegAuto{
			RpcUrl:         rpcUrl,
			Privkb58:       privkb58,
			Addr:           addr,
			Fee:            fee,
			Value:          value,
			AddrScriptHash: addrScriptHash,
			OntAddr:        ontAddr,
			Pwd:            pwd,
			User:           user,
		}
		handler.RunRegAuto()
	case "cctx":
		handler := btctool.CcTx{
			OntAddr:        ontAddr,
			AddrScriptHash: addrScriptHash,
			Value:          value,
			Fee:            fee,
			Addr:           addr,
			Privkb58:       privkb58,
			Indexes:        indexes,
			PrevHexTxs:     prevHexTxs,
			SpvAddr:        spvAddr,
			NetType:        netType,
		}
		handler.RunCcTx()
	case "blkgene":
		handler := btctool.BlkGene{
			User:        user,
			Pwd:         pwd,
			DefaultAddr: defaultAddr,
			Tsec:        tsec,
			RpcUrl:      rpcUrl,
		}
		handler.RunBlkGene()
	default:
		log.Errorf("no handler matched")
		os.Exit(1)
	}
}
