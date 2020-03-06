package gui

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/andlabs/ui"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/ontio/multi-chain/common/log"
	"github.com/zouxyan/btctool/service"
	"strconv"
	"strings"
)

func StartGui(quit chan struct{}) {
	err := ui.Main(func() {
		paramTab := ui.NewTab()
		yourTx := ui.NewMultilineEntry()
		yourTx.SetReadOnly(true)

		paramTab.Append("测试网", GetBoxForTest(yourTx))
		paramTab.Append("本地私网", GetBoxForReg(yourTx))
		paramTab.Append("为合约签名", GetBoxForSignRedeemContract(yourTx))
		paramTab.Append("注册多签合约", GetBoxForRegisterRedeem(yourTx))
		paramTab.Append("为交易参数签名", GetBoxForSignTxParam(yourTx))
		paramTab.Append("设置交易参数", GetBoxForSetTxParam(yourTx))
		paramTab.Append("加密私钥", GetBoxForEncryptPrivk(yourTx))
		paramTab.Append("生成私钥", GetBoxForGenePrivk(yourTx))
		paramTab.Append("生成多签Redeem", GetBoxForGeneRedeem(yourTx))

		resultBox := ui.NewVerticalBox()
		resultBox.Append(ui.NewLabel("结果:"), true)
		resultBox.Append(yourTx, false)

		div := ui.NewVerticalBox()
		div.Append(paramTab, false)
		div.Append(resultBox, true)
		div.SetPadded(false)

		window := ui.NewWindow("比特币跨链交易构造工具", 900, 600, false)
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


func GetBoxForReg(yourTx *ui.MultilineEntry) *ui.Box {
	regParam := ui.NewForm()
	fee := ui.NewEntry()
	privkb58 := ui.NewEntry()
	targetAddr := ui.NewEntry()
	value := ui.NewEntry()
	toChainId := ui.NewEntry()
	pwd := ui.NewEntry()
	user := ui.NewEntry()
	url := ui.NewEntry()
	multiAddr := ui.NewEntry()
	regParam.Append("跨链BTC金额: ", value, false)
	regParam.Append("BTC多签地址: ", multiAddr, false)
	regParam.Append("目标链ID: ", toChainId, false)
	regParam.Append("目标链地址: ", targetAddr, false)
	regParam.Append("BTC交易手续费: ", fee, false)
	regParam.Append("私钥(WIF): ", privkb58, false)
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
			ToChainId:    toChainIdVal,
			ToAddr:       multiAddr.Text(),
			NetParam:     &chaincfg.RegressionNetParams,
		}
		yourTx.SetText("txid:\n" + handler.Run())
		privkb58.SetText("")
	})

	return regBox
}

func GetBoxForTest(yourTx *ui.MultilineEntry) *ui.Box {
	testParam := ui.NewForm()
	feeT := ui.NewEntry()
	privkb58T := ui.NewEntry()
	targetAddrT := ui.NewEntry()
	valueT := ui.NewEntry()
	toChainIdT := ui.NewEntry()
	index := ui.NewEntry()
	utxoVals := ui.NewEntry()
	txids := ui.NewEntry()
	rpcPwd := ui.NewEntry()
	rpcUser := ui.NewEntry()
	rpcUrl := ui.NewEntry()
	multiAddrT := ui.NewEntry()

	testParam.Append("跨链BTC金额: ", valueT, false)
	testParam.Append("BTC多签地址: ", multiAddrT, false)
	testParam.Append("目标链ID: ", toChainIdT, false)
	testParam.Append("目标链地址: ", targetAddrT, false)
	testParam.Append("BTC交易手续费: ", feeT, false)
	testParam.Append("私钥(WIF): ", privkb58T, false)

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
	buttonTest.OnClicked(func(button *ui.Button) {
		feeVal, err := strconv.ParseFloat(feeT.Text(), 64)
		if err != nil {
			log.Errorf("failed to parse float %s: %v", feeT.Text(), err)
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
				OntAddr:      targetAddrT.Text(),
				ToAddr:       multiAddrT.Text(),
				Fee:          feeVal,
				Privkb58:     privkb58T.Text(),
				RpcUrl:       rpcUrl.Text(),
				User:         rpcUser.Text(),
				Pwd:          rpcPwd.Text(),
			}
			yourTx.SetText("txid:\n" + handler.Run())
		} else {
			valArr, err := GetVals(utxoVals.Text())
			if err != nil {
				log.Errorf("failed to get vals: %v", err)
			}
			handler := &service.TestTxBuilder{
				Privkb58:     privkb58T.Text(),
				Fee:          feeVal,
				Value:        valueVal,
				OntAddr:      targetAddrT.Text(),
				ToChainId:    toChainIdVal,
				ToAddr:       multiAddrT.Text(),
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
			privkb58T.SetText("")
		}
	})

	return testBox
}

func GetBoxForSignRedeemContract(yourTx *ui.MultilineEntry) *ui.Box {
	signParam := ui.NewForm()
	privk := ui.NewEntry()
	contractAddr := ui.NewEntry()
	contractId1 := ui.NewEntry()
	redeem := ui.NewEntry()
	cver := ui.NewEntry()
	signParam.Append("私钥(WIF)：", privk, false)
	signParam.Append("目标合约", contractAddr, false)
	signParam.Append("目标链ID：", contractId1, false)
	signParam.Append("多签Redeem", redeem, false)
	signParam.Append("合约版本", cver, false)
	buttonSign := ui.NewButton("签名")
	bbox := ui.NewHorizontalBox()
	bbox.Append(ui.NewLabel(""), true)
	bbox.Append(buttonSign, true)
	bbox.Append(ui.NewLabel(""), true)
	signBox := ui.NewVerticalBox()
	signBox.Append(signParam, false)
	signBox.Append(bbox, false)
	buttonSign.OnClicked(func(button *ui.Button) {
		v, err := strconv.ParseUint(cver.Text(), 10, 64)
		if err != nil {
			log.Fatalf("failed to get contract version: %v", err)
			return
		}
		cId, err := strconv.ParseUint(contractId1.Text(), 10, 64)
		if err != nil {
			log.Fatalf("failed to get contract id: %v", err)
			return
		}
		yourSig := service.GetSigForRedeemContract(contractAddr.Text(), redeem.Text(), privk.Text(), v, cId)
		yourTx.SetText(fmt.Sprintf("here is your sig, please remember it:\n%s", yourSig))
		privk.SetText("")
	})
	return signBox
}

func GetBoxForRegisterRedeem(yourTx *ui.MultilineEntry) *ui.Box {
	registerParam := ui.NewForm()
	rpcAllia := ui.NewEntry()
	ca := ui.NewEntry()
	cver := ui.NewEntry()
	redeem1 := ui.NewEntry()
	sigs := ui.NewEntry()
	walletFile := ui.NewEntry()
	wpwd := ui.NewEntry()
	contractId := ui.NewEntry()
	registerParam.Append("中继链RPC地址：", rpcAllia, false)
	registerParam.Append("目标合约", ca, false)
	registerParam.Append("合约版本", cver, false)
	registerParam.Append("多签Redeem：", redeem1, false)
	registerParam.Append("签名：", sigs, false)
	registerParam.Append("中继链钱包路径：", walletFile, false)
	registerParam.Append("钱包密码：", wpwd, false)
	registerParam.Append("目标链ID：", contractId, false)
	registerButton := ui.NewButton("注册")
	rbox := ui.NewHorizontalBox()
	rbox.Append(ui.NewLabel(""), true)
	rbox.Append(registerButton, true)
	rbox.Append(ui.NewLabel(""), true)
	registerBox := ui.NewVerticalBox()
	registerBox.Append(registerParam, false)
	registerBox.Append(rbox, false)
	registerButton.OnClicked(func(button *ui.Button) {
		cid, err := strconv.ParseUint(contractId.Text(), 10, 64)
		if err != nil {
			log.Fatalf("failed to get chain-id: %v", err)
			return
		}
		cver, err := strconv.ParseUint(cver.Text(), 10, 64)
		if err != nil {
			log.Fatalf("failed to get chain-id: %v", err)
			return
		}
		yourTx.SetText(fmt.Sprintf("your register tx hash is %s",
			service.RedeemRegister(rpcAllia.Text(), ca.Text(), redeem1.Text(), sigs.Text(), walletFile.Text(),
				wpwd.Text(), cid, cver)))
	})

	return registerBox
}

func GetBoxForSignTxParam(yourTx *ui.MultilineEntry) *ui.Box {
	signParam := ui.NewForm()
	privk := ui.NewEntry()
	redeem := ui.NewEntry()
	pver := ui.NewEntry()
	feeRate := ui.NewEntry()
	minChange := ui.NewEntry()
	signParam.Append("私钥(WIF)：", privk, false)
	signParam.Append("多签Redeem", redeem, false)
	signParam.Append("参数版本", pver, false)
	signParam.Append("费率（sat/byte）", feeRate, false)
	signParam.Append("最小找零值（sat）", minChange, false)

	buttonSign := ui.NewButton("签名")
	bbox := ui.NewHorizontalBox()
	bbox.Append(ui.NewLabel(""), true)
	bbox.Append(buttonSign, true)
	bbox.Append(ui.NewLabel(""), true)
	signBox := ui.NewVerticalBox()
	signBox.Append(signParam, false)
	signBox.Append(bbox, false)
	buttonSign.OnClicked(func(button *ui.Button) {
		fr, err := strconv.ParseUint(feeRate.Text(), 10, 64)
		if err != nil {
			log.Fatalf("failed to get fee rate: %v", err)
			return
		}
		mc, err := strconv.ParseUint(minChange.Text(), 10, 64)
		if err != nil {
			log.Fatalf("failed to get min change: %v", err)
			return
		}
		pv, err := strconv.ParseUint(pver.Text(), 10, 64)
		if err != nil {
			log.Fatalf("failed to get param version: %v", err)
			return
		}
		yourSig := service.GetSigForBtcTxParam(fr, mc, pv, redeem.Text(), privk.Text())
		yourTx.SetText(fmt.Sprintf("here is your sig, please remember it:\n%s", yourSig))
		privk.SetText("")
	})
	return signBox
}

func GetBoxForSetTxParam(yourTx *ui.MultilineEntry) *ui.Box {
	registerParam := ui.NewForm()
	rpcAllia := ui.NewEntry()
	redeem := ui.NewEntry()
	pver := ui.NewEntry()
	sigs := ui.NewEntry()
	feeRate := ui.NewEntry()
	minChange := ui.NewEntry()
	walletFile := ui.NewEntry()
	wpwd := ui.NewEntry()
	registerParam.Append("中继链RPC地址：", rpcAllia, false)
	registerParam.Append("多签Redeem：", redeem, false)
	registerParam.Append("参数版本", pver, false)
	registerParam.Append("签名：", sigs, false)
	registerParam.Append("费率（sat/byte）", feeRate, false)
	registerParam.Append("最小找零值（sat）", minChange, false)
	registerParam.Append("中继链钱包路径：", walletFile, false)
	registerParam.Append("钱包密码：", wpwd, false)
	registerButton := ui.NewButton("注册")
	rbox := ui.NewHorizontalBox()
	rbox.Append(ui.NewLabel(""), true)
	rbox.Append(registerButton, true)
	rbox.Append(ui.NewLabel(""), true)
	registerBox := ui.NewVerticalBox()
	registerBox.Append(registerParam, false)
	registerBox.Append(rbox, false)
	registerButton.OnClicked(func(button *ui.Button) {
		fr, err := strconv.ParseUint(feeRate.Text(), 10, 64)
		if err != nil {
			log.Fatalf("failed to get fee rate: %v", err)
			return
		}
		mc, err := strconv.ParseUint(minChange.Text(), 10, 64)
		if err != nil {
			log.Fatalf("failed to get min change: %v", err)
			return
		}
		pv, err := strconv.ParseUint(pver.Text(), 10, 64)
		if err != nil {
			log.Fatalf("failed to get param version: %v", err)
			return
		}
		yourTx.SetText(fmt.Sprintf("your tx hash is %s",
			service.SetBtcTxParam(rpcAllia.Text(), redeem.Text(), sigs.Text(), walletFile.Text(), wpwd.Text(), fr, mc, pv)))
	})

	return registerBox
}

func GetBoxForEncryptPrivk(yourTx *ui.MultilineEntry) *ui.Box {
	encryptParam := ui.NewForm()
	privkForEnc := ui.NewEntry()
	pwdForEnc := ui.NewPasswordEntry()
	onemore := ui.NewPasswordEntry()
	encryptParam.Append("私钥(WIF)：", privkForEnc, false)
	encryptParam.Append("输入密码：", pwdForEnc, false)
	encryptParam.Append("重复密码：", onemore, false)
	encButton := ui.NewButton("加密")
	ebox := ui.NewHorizontalBox()
	ebox.Append(ui.NewLabel(""), true)
	ebox.Append(encButton, true)
	ebox.Append(ui.NewLabel(""), true)
	encBox := ui.NewVerticalBox()
	encBox.Append(encryptParam, false)
	encBox.Append(ebox, false)
	encButton.OnClicked(func(button *ui.Button) {
		pwd, pwd1 := pwdForEnc.Text(), onemore.Text()
		if pwd != pwd1 {
			yourTx.SetText("两次密码输入不同")
			return
		}
		service.EncryptBtcPrivk(privkForEnc.Text(), pwd)
		yourTx.SetText(fmt.Sprintf("you can find your wallet file at %s", service.BTCPRIVK_PATH))
		privkForEnc.SetText("")
	})
	return encBox
}

func GetBoxForGenePrivk(yourTx *ui.MultilineEntry) *ui.Box {
	pParam := ui.NewForm()
	netEntry := ui.NewEntry()
	pParam.Append("比特币网络类型：", netEntry, false)
	pButton := ui.NewButton("获取")
	pBox := ui.NewHorizontalBox()
	pBox.Append(ui.NewLabel(""), true)
	pBox.Append(pButton, true)
	pBox.Append(ui.NewLabel(""), true)
	privBox := ui.NewVerticalBox()
	privBox.Append(pParam, false)
	privBox.Append(pBox, false)
	pButton.OnClicked(func(button *ui.Button) {
		yourTx.SetText(service.GetPrivk(netEntry.Text()))
	})
	return privBox
}

func GetBoxForGeneRedeem(yourTx *ui.MultilineEntry) *ui.Box {
	rParam := ui.NewForm()
	rnetEntry := ui.NewEntry()
	pubksEntry := ui.NewEntry()
	reqEntry := ui.NewEntry()
	rParam.Append("比特币网络类型：", rnetEntry, false)
	rParam.Append("公钥：", pubksEntry, false)
	rParam.Append("要求签名数目：", reqEntry, false)
	rButton := ui.NewButton("获取")
	rBox := ui.NewHorizontalBox()
	rBox.Append(ui.NewLabel(""), true)
	rBox.Append(rButton, true)
	rBox.Append(ui.NewLabel(""), true)
	pubksBox := ui.NewVerticalBox()
	pubksBox.Append(rParam, false)
	pubksBox.Append(rBox, false)
	rButton.OnClicked(func(button *ui.Button) {
		reqNum, _ := strconv.ParseInt(reqEntry.Text(), 10, 64)
		yourTx.SetText(service.GetRedeemForMultiSig(pubksEntry.Text(), rnetEntry.Text(), int(reqNum)))
	})
	return pubksBox
}

func GetVals(val string) ([]float64, error) {
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