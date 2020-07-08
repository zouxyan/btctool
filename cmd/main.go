package main

import (
	"flag"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/polynetwork/poly/common/log"
	"github.com/zouxyan/btctool/gui"
	"github.com/zouxyan/btctool/service"
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
var toChainId uint64
var wit int
var runGui int
var polyRpc string
var rdm string
var sigs string
var wallet string
var walletPwd string
var contractId uint64
var multiAddr string
var btcPrivPwd string
var pubks string
var req int
var cver uint64
var pver uint64
var feeRate uint64
var minChange uint64

//only for test
func init() {
	flag.StringVar(&tool, "tool", "", "which tool to use, \"reg\", \"test\" or \"blkgene\"")
	flag.StringVar(&prevHexTxs, "hextxs", "", "Raw transaction in hex")
	flag.StringVar(&indexes, "idxes", "", "Your UTXO's index for building this transaction")
	flag.StringVar(&privkb58, "privkb58", "", "Your private key in base58")
	flag.StringVar(&addr, "addr", "", "Your btc address in base58")
	flag.Float64Var(&value, "value", 0.0001, "Amount of btc to cross chain")
	flag.Float64Var(&fee, "fee", 0.00001, "Cross chain fee, not in use for now")
	flag.StringVar(&ontAddr, "targetaddr", "", "Your target chain address")
	flag.UintVar(&tsec, "tsec", 5, "block update time(seconds), 5 sec default")
	flag.StringVar(&rpcUrl, "url", "", "the bitcoind rpc address")
	flag.StringVar(&user, "user", "", "the rpc user")
	flag.StringVar(&pwd, "pwd", "", "the rpc password")
	flag.StringVar(&defaultAddr, "defaultaddr", "", "the default bitcoin address to rececive the mining reward")
	flag.StringVar(&netType, "net", "test", "the net type of bitcoin")
	flag.StringVar(&utxoVals, "utxovals", "", "val of utxos in satoshi, eg: 1000,2000")
	flag.StringVar(&txids, "txids", "", "txid of utxos, eg: xx,yy,bb")
	flag.StringVar(&contractAddr, "contract", "", "target chain smart contract address")
	flag.Int64Var(&dura, "dura", 300, "set the seconds to send a cross-tx, default 5 min")
	flag.Int64Var(&maxVal, "maxval", 2000, "the max value of cross tx")
	flag.Uint64Var(&toChainId, "tochain", 3, "target chain id")
	flag.IntVar(&wit, "wit", 0, "use segwit for output")
	flag.IntVar(&runGui, "gui", 1, "run gui")
	flag.StringVar(&polyRpc, "poly-rpc", "", "poly chain rpc address")
	flag.StringVar(&rdm, "redeem", "", "your redeem script")
	flag.StringVar(&sigs, "sigs", "", "your sig for redeem register")
	flag.StringVar(&wallet, "wallet", "", "OR chain wallet file path")
	flag.StringVar(&walletPwd, "wallet-pwd", "", "OR chain wallet password")
	flag.Uint64Var(&contractId, "contractId", 2, "chain id of your contract")
	flag.StringVar(&multiAddr, "multiaddr", "", "multisign-addr of redeem")
	flag.StringVar(&btcPrivPwd, "btcpwd", "", "password for btc privk encryption")
	flag.StringVar(&pubks, "pubks", "", "public keys to create redeem-script")
	flag.IntVar(&req, "require", 0, "require number for redeem")
	flag.Uint64Var(&cver, "contract_ver", 0, "your smart contract version")
	flag.Uint64Var(&feeRate, "fee_rate", 1, "fee rate for unlocking btc transaction (sat/byte)")
	flag.Uint64Var(&minChange, "min_change", 2000, "min change limit for unlocking btc transaction (satoshi)")
	flag.Uint64Var(&pver, "param_ver", 0, "your btc transaction parameters version")
}

func main() {
	flag.Parse()

	quit := make(chan struct{})
	if runGui == 1 {
		log.InitLog(0, "./log/", os.Stdout)
		gui.StartGui(quit)
		select {
		case <-quit:
			return
		}
	}
	log.InitLog(0, os.Stdout)
	switch tool {
	case "reg":
		handler := &service.RegTxBuilder{
			RpcUrl:    rpcUrl,
			Privkb58:  privkb58,
			Fee:       fee,
			Value:     value,
			OntAddr:   ontAddr,
			Pwd:       pwd,
			User:      user,
			ToChainId: toChainId,
			ToAddr:    multiAddr,
			NetParam:  &chaincfg.RegressionNetParams,
		}
		handler.Run()
	case "test":
		if rpcUrl != "" {
			handler := service.RegTxBuilder{
				NetParam:  &chaincfg.TestNet3Params,
				Privkb58:  privkb58,
				Fee:       fee,
				Value:     value,
				OntAddr:   ontAddr,
				Pwd:       pwd,
				User:      user,
				ToChainId: toChainId,
				ToAddr:    multiAddr,
				RpcUrl:    rpcUrl,
			}
			handler.Run()
		} else {
			valArr, err := gui.GetVals(utxoVals)
			if err != nil {
				log.Errorf("failed to get vals: %v", err)
				os.Exit(1)
			}
			handler := service.TestTxBuilder{
				OntAddr:   ontAddr,
				Value:     value,
				Fee:       fee,
				Privkb58:  privkb58,
				Indexes:   indexes,
				NetType:   netType,
				Vals:      valArr,
				Txids:     txids,
				ToAddr:    multiAddr,
				ToChainId: toChainId,
			}
			handler.Run()
		}
	case "blkgene":
		handler := service.BlkGene{
			User:        user,
			Pwd:         pwd,
			DefaultAddr: defaultAddr,
			Tsec:        tsec,
			RpcUrl:      rpcUrl,
		}
		handler.Run()
	case "autosender":
		valArr, err := gui.GetVals(utxoVals)
		if err != nil {
			log.Errorf("failed to get vals: %v", err)
			os.Exit(1)
		}

		handler := service.AutoSender{
			CcTx: &service.TestTxBuilder{
				OntAddr:  ontAddr,
				Value:    value,
				Fee:      fee,
				Privkb58: privkb58,
				Indexes:  indexes,
				NetType:  netType,
				Vals:     valArr,
				Txids:    txids,
				ToAddr:   multiAddr,
			},
			MaxVal: maxVal,
			Dura:   dura,
		}
		handler.Run()
	case "utxocounter":
		service.CountPolyUtxo(polyRpc)
	case "register_redeem":
		service.RedeemRegister(polyRpc, contractAddr, rdm, sigs, wallet, walletPwd, contractId, cver)
	case "sign_redeem_contract":
		service.GetSigForRedeemContract(contractAddr, rdm, privkb58, cver, contractId)
	case "set_tx_param":
		service.SetBtcTxParam(polyRpc, rdm, sigs, wallet, walletPwd, feeRate, minChange, pver)
	case "sign_tx_param":
		service.GetSigForBtcTxParam(feeRate, minChange, pver, rdm, privkb58)
	case "encrypt_privk":
		service.EncryptBtcPrivk(privkb58, btcPrivPwd)
	case "getprivk":
		log.Infof("your privk in WIF is %s", service.GetPrivk(netType))
	case "getredeem":
		log.Infof(service.GetRedeemForMultiSig(pubks, netType, req))

	default:
		log.Errorf("no handler matched")
		os.Exit(1)
	}
}
