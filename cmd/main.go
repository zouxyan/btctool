package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/ontio/multi-chain/common/log"
	"github.com/zouxyan/btctool/service"
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
var toChainId uint64
var wit int
var runGui int
var alliaRpc string

//only for test
// TODO: use redeem to get these addresses
var redeem string = "5521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57ae"

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
	flag.Uint64Var(&toChainId, "tochain", 3, "target chain id")
	flag.IntVar(&wit, "wit", 0, "use segwit for output")
	flag.IntVar(&runGui, "gui", 1, "run gui")
	flag.StringVar(&alliaRpc, "allia-rpc", "", "alliance chain rpc address")
}

func main() {
	flag.Parse()
	quit := make(chan struct{})
	if runGui == 1 {
		log.InitLog(0, "./log/", os.Stdout)
		startGui(quit)
		select {
		case <-quit:
			return
		}
	}
	log.InitLog(0, os.Stdout)
	switch tool {
	case "reg":
		handler := &service.RegTxBuilder{
			RpcUrl:       rpcUrl,
			Privkb58:     privkb58,
			Fee:          fee,
			Value:        value,
			OntAddr:      ontAddr,
			Pwd:          pwd,
			User:         user,
			ContractAddr: contractAddr,
			ToChainId:    toChainId,
			IsSegWit:     wit,
			Redeem:       redeem,
			NetParam:     &chaincfg.RegressionNetParams,
		}
		handler.Run()
	case "test":
		if rpcUrl != "" {
			handler := service.RegTxBuilder{
				NetParam:     &chaincfg.TestNet3Params,
				Privkb58:     privkb58,
				Fee:          fee,
				Value:        value,
				OntAddr:      ontAddr,
				Pwd:          pwd,
				User:         user,
				ContractAddr: contractAddr,
				ToChainId:    toChainId,
				IsSegWit:     wit,
				Redeem:       redeem,
				RpcUrl:       rpcUrl,
			}
			handler.Run()
		} else {
			valArr, err := getVals(utxoVals)
			if err != nil {
				log.Errorf("failed to get vals: %v", err)
				os.Exit(1)
			}
			handler := service.TestTxBuilder{
				OntAddr:      ontAddr,
				Value:        value,
				Fee:          fee,
				Privkb58:     privkb58,
				Indexes:      indexes,
				SpvAddr:      spvAddr,
				NetType:      netType,
				Vals:         valArr,
				Txids:        txids,
				ContractAddr: contractAddr,
				IsSegWit:     wit,
				Redeem:       redeem,
				ToChainId:    toChainId,
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
		valArr, err := getVals(utxoVals)
		if err != nil {
			log.Errorf("failed to get vals: %v", err)
			os.Exit(1)
		}

		handler := service.AutoSender{
			CcTx: &service.TestTxBuilder{
				OntAddr:      ontAddr,
				Value:        value,
				Fee:          fee,
				Privkb58:     privkb58,
				Indexes:      indexes,
				SpvAddr:      spvAddr,
				NetType:      netType,
				Vals:         valArr,
				Txids:        txids,
				ContractAddr: contractAddr,
				Redeem:       redeem,
				IsSegWit:     wit,
			},
			MaxVal: maxVal,
			Dura:   dura,
		}
		handler.Run()
	case "utxocounter":
		service.CountAlliaUtxo(alliaRpc)
	case "syncgh":

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

func startGui(quit chan struct{}) {
	err := ui.Main(func() {
		paramTab := ui.NewTab()
		yourTx := ui.NewMultilineEntry()
		yourTx.SetReadOnly(true)

		regParam := ui.NewForm()
		fee := ui.NewEntry()
		privkb58 := ui.NewEntry()
		targetAddr := ui.NewEntry()
		value := ui.NewEntry()
		contract := ui.NewEntry()
		toChainId := ui.NewEntry()
		pwd := ui.NewEntry()
		user := ui.NewEntry()
		url := ui.NewEntry()
		regParam.Append("跨链BTC金额: ", value, false)
		regParam.Append("目标链代币合约哈希: ", contract, false)
		regParam.Append("目标链ID: ", toChainId, false)
		regParam.Append("目标链地址: ", targetAddr, false)
		regParam.Append("BTC交易手续费: ", fee, false)
		regParam.Append("私钥(Base58): ", privkb58, false)
		regParam.Append("rpcURL: ", url, false)
		regParam.Append("rpc用户: ", user, false)
		regParam.Append("rpc密码: ", pwd, false)
		buttonReg := ui.NewButton("获取交易")
		rbbox := ui.NewHorizontalBox()
		rbbox.Append(ui.NewLabel(""), true)
		rbbox.Append(buttonReg, true)
		rbbox.Append(ui.NewLabel(""), true)
		regBox := ui.NewVerticalBox()
		regBox.Append(regParam, false)
		regBox.Append(rbbox, false)
		buttonReg.OnClicked(func(button *ui.Button) {
			feeVal, err := strconv.ParseFloat(fee.Text(), 64)
			if err != nil {
				log.Errorf("failed to parse float %s: %v", fee.Text(), err)
				//TODO: log it
			}
			valueVal, err := strconv.ParseFloat(value.Text(), 64)
			if err != nil {
				log.Errorf("failed to parse value: %v", err)
			}
			toChainIdVal, err := strconv.ParseUint(toChainId.Text(), 10, 64)
			if err != nil {
				log.Errorf("failed to parse tochain id: %v", err)
			}

			handler := &service.RegTxBuilder{
				RpcUrl:       url.Text(),
				Privkb58:     privkb58.Text(),
				Fee:          feeVal,
				Value:        valueVal,
				OntAddr:      targetAddr.Text(),
				Pwd:          pwd.Text(),
				User:         user.Text(),
				ContractAddr: contract.Text(),
				ToChainId:    toChainIdVal,
				IsSegWit:     0,
				Redeem:       "5521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57ae",
			}
			yourTx.SetText("txid:\n" + handler.Run())
		})

		testParam := ui.NewForm()
		feeT := ui.NewEntry()
		privkb58T := ui.NewEntry()
		targetAddrT := ui.NewEntry()
		valueT := ui.NewEntry()
		contractT := ui.NewEntry()
		toChainIdT := ui.NewEntry()
		index := ui.NewEntry()
		utxoVals := ui.NewEntry()
		txids := ui.NewEntry()
		rpcPwd := ui.NewEntry()
		rpcUser := ui.NewEntry()
		rpcUrl := ui.NewEntry()

		testParam.Append("跨链BTC金额: ", valueT, false)
		testParam.Append("目标链代币合约哈希: ", contractT, false)
		testParam.Append("目标链ID: ", toChainIdT, false)
		testParam.Append("目标链地址: ", targetAddrT, false)
		testParam.Append("BTC交易手续费: ", feeT, false)
		testParam.Append("私钥(Base58): ", privkb58T, false)

		utxoTab := ui.NewTab()
		inputForm := ui.NewForm()
		inputForm.Append("UTXO的index: ", index, false)
		inputForm.Append("UTXO的金额: ", utxoVals, false)
		inputForm.Append("UTXO的交易ID: ", txids, false)
		utxoTab.Append("自行填写", inputForm)
		byRpcForm := ui.NewForm()
		byRpcForm.Append("全节点的URL: ", rpcUrl, false)
		byRpcForm.Append("RPC用户名: ", rpcUser, false)
		byRpcForm.Append("RPC密码: ", rpcPwd, false)
		utxoTab.Append("RPC自动获取", byRpcForm)
		testParam.Append("作为输入的UTXO:  \n\n", utxoTab, true)

		buttonTest := ui.NewButton("获取交易")
		hbox := ui.NewHorizontalBox()
		hbox.Append(ui.NewLabel(""), true)
		hbox.Append(buttonTest, true)
		hbox.Append(ui.NewLabel(""), true)
		testBox := ui.NewVerticalBox()
		testBox.Append(testParam, false)
		testBox.Append(hbox, false)
		paramTab.Append("构造测试网跨链交易", testBox)
		paramTab.Append("构造本地私网跨链交易", regBox)
		buttonTest.OnClicked(func(button *ui.Button) {
			feeVal, err := strconv.ParseFloat(feeT.Text(), 64)
			if err != nil {
				log.Errorf("failed to parse float %s: %v", fee.Text(), err)
			}
			valueVal, err := strconv.ParseFloat(valueT.Text(), 64)
			if err != nil {
				log.Errorf("failed to parse value: %v", err)
			}
			toChainIdVal, err := strconv.ParseUint(toChainIdT.Text(), 10, 64)
			if err != nil {
				log.Errorf("failed to parse tochain id: %v", err)
			}
			if rpcUrl.Text() != "" {
				handler := &service.RegTxBuilder{
					NetParam:     &chaincfg.TestNet3Params,
					ToChainId:    toChainIdVal,
					Value:        valueVal,
					ContractAddr: contractT.Text(),
					OntAddr:      targetAddrT.Text(),
					Redeem:       "5521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57ae",
					IsSegWit:     0,
					Fee:          feeVal,
					Privkb58:     privkb58T.Text(),
					RpcUrl:       rpcUrl.Text(),
					User:         rpcUser.Text(),
					Pwd:          rpcPwd.Text(),
				}
				yourTx.SetText("txid:\n" + handler.Run())
			} else {
				valArr, err := getVals(utxoVals.Text())
				if err != nil {
					log.Errorf("failed to get vals: %v", err)
				}
				handler := &service.TestTxBuilder{
					Privkb58:     privkb58T.Text(),
					Fee:          feeVal,
					Value:        valueVal,
					OntAddr:      targetAddrT.Text(),
					ContractAddr: contractT.Text(),
					ToChainId:    toChainIdVal,
					IsSegWit:     0,
					Redeem:       "5521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57ae",
					NetType:      "test",
					Indexes:      index.Text(),
					Vals:         valArr,
					Txids:        txids.Text(),
				}
				tx := handler.Run()
				var buf bytes.Buffer
				err = tx.BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
				if err != nil {
					log.Errorf("%v\n", err)
				}
				yourTx.SetText(fmt.Sprintf("you can use %s to broadcast tx: \n%s", "https://tbtc.bitaps.com/broadcast",
					hex.EncodeToString(buf.Bytes())))
			}
		})

		resultBox := ui.NewVerticalBox()
		resultBox.Append(ui.NewLabel("结果:"), true)
		resultBox.Append(yourTx, false)

		div := ui.NewVerticalBox()
		div.Append(paramTab, false)
		div.Append(resultBox, true)
		div.SetPadded(false)

		window := ui.NewWindow("比特币跨链交易构造工具", 600, 600, false)
		window.SetChild(div)
		window.SetMargined(true)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			close(quit)
			log.Info("stopping gui")
			return true
		})

		window.Show()
	})
	if err != nil {
		log.Errorf("gui error: %v", err)
		close(quit)
	}
}
