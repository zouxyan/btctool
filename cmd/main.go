package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/andlabs/ui"
	"github.com/btcsuite/btcd/wire"
	"github.com/zouxyan/btctool/service"
	"github.com/ontio/multi-chain/common/log"
	"os"
	"strconv"
	"strings"
	_ "github.com/andlabs/ui/winmanifest"
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

//only for test
// TODO: use redeem to get these addresses
var redeem string = "5521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57ae"

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
	flag.Uint64Var(&toChainId, "tochain", 2, "target chain id")
	flag.IntVar(&wit, "wit", 0, "use segwit for output")
	flag.IntVar(&runGui, "gui", 1, "run gui")
}

func main() {
	flag.Parse()
	log.InitLog(0, "./log/", os.Stdout)
	quit := make(chan struct{})
	if runGui == 1 {
		startGui(quit)
		select {
		case <-quit:
			return
		}
	}
	switch tool {
	case "regauto":
		handler := &service.RegAuto{
			RpcUrl:         rpcUrl,
			Privkb58:       privkb58,
			Fee:            fee,
			Value:          value,
			OntAddr:        ontAddr,
			Pwd:            pwd,
			User:           user,
			ContractAddr:   contractAddr,
			ToChainId: toChainId,
			IsSegWit: wit,
			Redeem: redeem,
		}
		handler.Run()
	case "cctx":
		valArr, err := getVals(utxoVals)
		if err != nil {
			log.Errorf("failed to get vals: %v", err)
			os.Exit(1)
		}

		handler := service.CcTx{
			OntAddr:        ontAddr,
			Value:          value,
			Fee:            fee,
			Privkb58:       privkb58,
			Indexes:        indexes,
			SpvAddr:        spvAddr,
			NetType:        netType,
			Vals:           valArr,
			Txids:          txids,
			ContractAddr:   contractAddr,
			IsSegWit: wit,
			Redeem: redeem,
		}
		handler.Run()
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
			CcTx: &service.CcTx{
				OntAddr:        ontAddr,
				Value:          value,
				Fee:            fee,
				Privkb58:       privkb58,
				Indexes:        indexes,
				SpvAddr:        spvAddr,
				NetType:        netType,
				Vals:           valArr,
				Txids:          txids,
				ContractAddr:   contractAddr,
				Redeem: redeem,
				IsSegWit: wit,
			},
			MaxVal: maxVal,
			Dura:   dura,
		}
		handler.Run()
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
		box1 := ui.NewVerticalBox()
		info := ui.NewLabel("选择工具:")
		combo := ui.NewCombobox()
		combo.Append("构造本地私网跨链交易")
		combo.Append("构造测试网跨链交易")
		//combo.Append("autosender")
		combo.SetSelected(1)
		box1.Append(info, false)
		box1.Append(combo, false)

		paramBox := ui.NewVerticalBox()
		fee := ui.NewEntry()
		privkb58 := ui.NewEntry()
		targetAddr := ui.NewEntry()
		value := ui.NewEntry()
		contract := ui.NewEntry()
		toChainId := ui.NewEntry()
		paramBox.Append(ui.NewLabel("跨链BTC金额:"), false)
		paramBox.Append(value, false)
		paramBox.Append(ui.NewLabel("目标链代币合约哈希:"), false)
		paramBox.Append(contract, false)
		paramBox.Append(ui.NewLabel("目标链ID:"), false)
		paramBox.Append(toChainId, false)
		paramBox.Append(ui.NewLabel("目标链地址:"), false)
		paramBox.Append(targetAddr, false)
		paramBox.Append(ui.NewLabel("BTC交易手续费:"), false)
		paramBox.Append(fee, false)
		paramBox.Append(ui.NewLabel("私钥(Base58):"), false)
		paramBox.Append(privkb58, false)

		pwd := ui.NewEntry()
		user := ui.NewEntry()
		url := ui.NewEntry()

		index := ui.NewEntry()
		utxoVals := ui.NewEntry()
		txids := ui.NewEntry()

		regBox := ui.NewVerticalBox()
		regBox.Append(ui.NewLabel("rpc用户:"), false)
		regBox.Append(user, false)
		regBox.Append(ui.NewLabel("rpc密码:"), false)
		regBox.Append(pwd, false)
		regBox.Append(ui.NewLabel("rpcURL:"), false)
		regBox.Append(url, false)
		paramBox.Append(regBox, false)

		cctxBox := ui.NewVerticalBox()
		cctxBox.Append(ui.NewLabel("UTXO的index:"), false)
		cctxBox.Append(index, false)
		cctxBox.Append(ui.NewLabel("UTXO的金额:"), false)
		cctxBox.Append(utxoVals, false)
		cctxBox.Append(ui.NewLabel("UTXO的交易ID:"), false)
		cctxBox.Append(txids, false)
		paramBox.Append(cctxBox, false)

		regBox.Hide()
		//cctxBox.Hide()
		var tool string = "regauto"
		combo.OnSelected(func(combobox *ui.Combobox) {
			switch combo.Selected() {
			case 0:
				tool = "regauto"
				regBox.Show()
				cctxBox.Hide()
			case 1:
				tool = "cctx"
				cctxBox.Show()
				regBox.Hide()
			case 2:
				tool = "autosender"
			default:
				// log
				fmt.Println(combo.Selected())
			}
		})

		resultBox := ui.NewVerticalBox()
		yourTx := ui.NewMultilineEntry()
		yourTx.SetReadOnly(true)
		button := ui.NewButton("获取交易")
		button.OnClicked(func(button *ui.Button) {
			feeVal, err := strconv.ParseFloat(fee.Text(), 64)
			if err != nil {
				log.Errorf("failed to ")
				os.Exit(1) //TODO: log it
			}
			valueVal, err := strconv.ParseFloat(value.Text(), 64)
			if err != nil {
				log.Errorf("failed to parse value: %v", err)
				os.Exit(1)
			}
			toChainIdVal, err := strconv.ParseUint(toChainId.Text(), 10, 64)
			if err != nil {
				log.Errorf("failed to parse tochain id: %v", err)
				os.Exit(1)
			}

			switch tool {
			case "regauto":
				handler := &service.RegAuto{
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
			case "cctx":
				valArr, err := getVals(utxoVals.Text())
				if err != nil {
					log.Errorf("failed to get vals: %v", err)
					os.Exit(1)
				}

				handler := &service.CcTx{
					Privkb58:     privkb58.Text(),
					Fee:          feeVal,
					Value:        valueVal,
					OntAddr:      targetAddr.Text(),
					ContractAddr: contract.Text(),
					ToChainId:    toChainIdVal,
					IsSegWit:     0,
					Redeem:       "5521023ac710e73e1410718530b2686ce47f12fa3c470a9eb6085976b70b01c64c9f732102c9dc4d8f419e325bbef0fe039ed6feaf2079a2ef7b27336ddb79be2ea6e334bf2102eac939f2f0873894d8bf0ef2f8bbdd32e4290cbf9632b59dee743529c0af9e802103378b4a3854c88cca8bfed2558e9875a144521df4a75ab37a206049ccef12be692103495a81957ce65e3359c114e6c2fe9f97568be491e3f24d6fa66cc542e360cd662102d43e29299971e802160a92cfcd4037e8ae83fb8f6af138684bebdc5686f3b9db21031e415c04cbc9b81fbee6e04d8c902e8f61109a2c9883a959ba528c52698c055a57ae",
					NetType: "test",
					Indexes: index.Text(),
					Vals: valArr,
					Txids: txids.Text(),
				}
				tx := handler.Run()
				var buf bytes.Buffer
				err = tx.BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
				yourTx.SetText(fmt.Sprintf("you can use %s to broadcast tx: \n%s", "https://tbtc.bitaps.com/broadcast",
					hex.EncodeToString(buf.Bytes())))
			default:
				os.Exit(1)
			}
		})
		resultBox.Append(button, false)
		resultBox.Append(ui.NewLabel("结果:"), true)
		resultBox.Append(yourTx, false)

		div := ui.NewVerticalBox()
		div.Append(ui.NewLabel("-------------------首先，选择工具-------------------"), false)
		div.Append(box1, false)
		div.Append(ui.NewLabel("\n-------------------然后，填写参数-------------------"), false)
		div.Append(paramBox, false)
		div.Append(ui.NewLabel("\n-------------------最后，点击按钮构造交易-------------------"), false)
		div.Append(resultBox, true)
		div.SetPadded(false)

		window := ui.NewWindow("比特币跨链交易构造工具", 600, 1500, false)
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