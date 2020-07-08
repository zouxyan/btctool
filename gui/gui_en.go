package gui

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/andlabs/ui"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/polynetwork/poly/common/log"
	"github.com/zouxyan/btctool/service"
	"strconv"
	"time"
)

func GetENBox() *ui.Box {
	paramTab := ui.NewTab()

	res := ui.NewMultilineEntry()
	res.SetReadOnly(true)
	resultBox := ui.NewVerticalBox()
	resultBox.Append(ui.NewLabel("Result:"), true)
	resultBox.Append(res, false)

	paramTab.Append("Testnet", GetBoxForTestEN(res))
	paramTab.Append("Regnet", GetBoxForRegEN(res))
	paramTab.Append("SignContract", GetBoxForSignRedeemContractEN(res))
	paramTab.Append("RegisterContract", GetBoxForRegisterRedeemEN(res))
	paramTab.Append("SignParam", GetBoxForSignTxParamEN(res))
	paramTab.Append("SetParam", GetBoxForSetTxParamEN(res))
	paramTab.Append("EncrptPrivk", GetBoxForEncryptPrivkEN(res))
	paramTab.Append("GeneratePrivk", GetBoxForGenePrivkEN(res))
	paramTab.Append("GenerateRedeem", GetBoxForGeneRedeemEN(res))
	paramTab.Append("UTXOMonitor", GetBoxForUtxoMonitorEN(res))

	div := ui.NewVerticalBox()
	div.Append(paramTab, false)
	div.Append(resultBox, true)
	div.SetPadded(false)

	return div
}

func GetBoxForRegEN(res *ui.MultilineEntry) *ui.Box {
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
	regParam.Append("BTC Amount: ", value, false)
	regParam.Append("Multisig Address: ", multiAddr, false)
	regParam.Append("Target Chain ID: ", toChainId, false)
	regParam.Append("Target Chain Address: ", targetAddr, false)
	regParam.Append("BTC Tx Fee: ", fee, false)
	regParam.Append("Private Key(WIF): ", privkb58, false)
	regParam.Append("rpcURL: ", url, false)
	regParam.Append("rpc user: ", user, false)
	regParam.Append("rpc pwd: ", pwd, false)
	buttonReg := ui.NewButton("Get Tx")
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
			RpcUrl:    url.Text(),
			Privkb58:  privkb58.Text(),
			Fee:       feeVal,
			Value:     valueVal,
			OntAddr:   targetAddr.Text(),
			Pwd:       pwd.Text(),
			User:      user.Text(),
			ToChainId: toChainIdVal,
			ToAddr:    multiAddr.Text(),
			NetParam:  &chaincfg.RegressionNetParams,
		}
		res.SetText("txid:\n" + handler.Run())
		privkb58.SetText("")
	})

	return regBox
}

func GetBoxForTestEN(res *ui.MultilineEntry) *ui.Box {
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

	testParam.Append("BTC Amount: ", valueT, false)
	testParam.Append("Multisig Address: ", multiAddrT, false)
	testParam.Append("Target Chain ID: ", toChainIdT, false)
	testParam.Append("Target Chain Address: ", targetAddrT, false)
	testParam.Append("BTC Tx Fee: ", feeT, false)
	testParam.Append("Private Key(WIF): ", privkb58T, false)

	utxoTab := ui.NewTab()
	inputForm := ui.NewForm()
	inputForm.Append("UTXO index: ", index, false)
	inputForm.Append("UTXO Amount: ", utxoVals, false)
	inputForm.Append("UTXO txid: ", txids, false)
	utxoTab.Append("Set UTXO by yourself", inputForm)
	byRpcForm := ui.NewForm()
	byRpcForm.Append("Bitcoind RPC-URL: ", rpcUrl, false)
	byRpcForm.Append("RPC User: ", rpcUser, false)
	byRpcForm.Append("RPC Password: ", rpcPwd, false)
	utxoTab.Append("Auto-Set by RPC", byRpcForm)
	testParam.Append("UTXO to spend:  \n\n", utxoTab, true)

	buttonTest := ui.NewButton("Get Tx")
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
				NetParam:  &chaincfg.TestNet3Params,
				ToChainId: toChainIdVal,
				Value:     valueVal,
				OntAddr:   targetAddrT.Text(),
				ToAddr:    multiAddrT.Text(),
				Fee:       feeVal,
				Privkb58:  privkb58T.Text(),
				RpcUrl:    rpcUrl.Text(),
				User:      rpcUser.Text(),
				Pwd:       rpcPwd.Text(),
			}
			res.SetText("txid:\n" + handler.Run())
		} else {
			valArr, err := GetVals(utxoVals.Text())
			if err != nil {
				log.Errorf("failed to get vals: %v", err)
			}
			handler := &service.TestTxBuilder{
				Privkb58:  privkb58T.Text(),
				Fee:       feeVal,
				Value:     valueVal,
				OntAddr:   targetAddrT.Text(),
				ToChainId: toChainIdVal,
				ToAddr:    multiAddrT.Text(),
				NetType:   "test",
				Indexes:   index.Text(),
				Vals:      valArr,
				Txids:     txids.Text(),
			}
			tx := handler.Run()
			var buf bytes.Buffer
			err = tx.BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
			if err != nil {
				log.Errorf("%v\n", err)
			}
			res.SetText(fmt.Sprintf("you can use %s to broadcast tx: \n%s", "https://tbtc.bitaps.com/broadcast",
				hex.EncodeToString(buf.Bytes())))
			privkb58T.SetText("")
		}
	})

	return testBox
}

func GetBoxForSignRedeemContractEN(res *ui.MultilineEntry) *ui.Box {
	signParam := ui.NewForm()
	privk := ui.NewEntry()
	contractAddr := ui.NewEntry()
	contractId1 := ui.NewEntry()
	redeem := ui.NewEntry()
	cver := ui.NewEntry()
	signParam.Append("Private Key(WIF)：", privk, false)
	signParam.Append("SmartContract: ", contractAddr, false)
	signParam.Append("Target Chain ID：", contractId1, false)
	signParam.Append("Multisig Redeem: ", redeem, false)
	signParam.Append("Contract Version: ", cver, false)
	buttonSign := ui.NewButton("Sign")
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
		res.SetText(fmt.Sprintf("here is your sig, please remember it:\n%s", yourSig))
		privk.SetText("")
	})
	return signBox
}

func GetBoxForRegisterRedeemEN(res *ui.MultilineEntry) *ui.Box {
	registerParam := ui.NewForm()
	rpcPoly := ui.NewEntry()
	ca := ui.NewEntry()
	cver := ui.NewEntry()
	redeem1 := ui.NewEntry()
	sigs := ui.NewEntry()
	walletFile := ui.NewEntry()
	wpwd := ui.NewEntry()
	contractId := ui.NewEntry()
	registerParam.Append("poly RPC-URL：", rpcPoly, false)
	registerParam.Append("SmartContract: ", ca, false)
	registerParam.Append("Contract Version: ", cver, false)
	registerParam.Append("Multisig Redeem：", redeem1, false)
	registerParam.Append("Signatures：", sigs, false)
	registerParam.Append("poly Wallet Path：", walletFile, false)
	registerParam.Append("Wallet Password：", wpwd, false)
	registerParam.Append("Target Chain ID：", contractId, false)
	registerButton := ui.NewButton("Register")
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
		res.SetText(fmt.Sprintf("your register tx hash is %s",
			service.RedeemRegister(rpcPoly.Text(), ca.Text(), redeem1.Text(), sigs.Text(), walletFile.Text(),
				wpwd.Text(), cid, cver)))
	})

	return registerBox
}

func GetBoxForSignTxParamEN(res *ui.MultilineEntry) *ui.Box {
	signParam := ui.NewForm()
	privk := ui.NewEntry()
	redeem := ui.NewEntry()
	pver := ui.NewEntry()
	feeRate := ui.NewEntry()
	minChange := ui.NewEntry()
	signParam.Append("Private Key(WIF)：", privk, false)
	signParam.Append("Multisig Redeem：", redeem, false)
	signParam.Append("Param Version：", pver, false)
	signParam.Append("Fee Rate（sat/byte）：", feeRate, false)
	signParam.Append("Min Change（sat）：", minChange, false)

	buttonSign := ui.NewButton("Sign")
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
		res.SetText(fmt.Sprintf("here is your sig, please remember it:\n%s", yourSig))
		privk.SetText("")
	})
	return signBox
}

func GetBoxForSetTxParamEN(res *ui.MultilineEntry) *ui.Box {
	registerParam := ui.NewForm()
	rpcPoly := ui.NewEntry()
	redeem := ui.NewEntry()
	pver := ui.NewEntry()
	sigs := ui.NewEntry()
	feeRate := ui.NewEntry()
	minChange := ui.NewEntry()
	walletFile := ui.NewEntry()
	wpwd := ui.NewEntry()
	registerParam.Append("poly RPC-URL：", rpcPoly, false)
	registerParam.Append("Multisig Redeem：", redeem, false)
	registerParam.Append("Param Version：", pver, false)
	registerParam.Append("Signatures：", sigs, false)
	registerParam.Append("Fee Rate（sat/byte）：", feeRate, false)
	registerParam.Append("Min Change（sat）：", minChange, false)
	registerParam.Append("poly Wallet Path：", walletFile, false)
	registerParam.Append("Wallet Password：", wpwd, false)
	registerButton := ui.NewButton("Register")
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
		res.SetText(fmt.Sprintf("your tx hash is %s",
			service.SetBtcTxParam(rpcPoly.Text(), redeem.Text(), sigs.Text(), walletFile.Text(), wpwd.Text(), fr, mc, pv)))
	})

	return registerBox
}

func GetBoxForEncryptPrivkEN(res *ui.MultilineEntry) *ui.Box {
	encryptParam := ui.NewForm()
	privkForEnc := ui.NewEntry()
	pwdForEnc := ui.NewPasswordEntry()
	onemore := ui.NewPasswordEntry()
	encryptParam.Append("Private Key(WIF)：", privkForEnc, false)
	encryptParam.Append("Input Password：", pwdForEnc, false)
	encryptParam.Append("Repeat Password：", onemore, false)
	encButton := ui.NewButton("Encrypt")
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
			res.SetText("Password input twice is different")
			return
		}
		service.EncryptBtcPrivk(privkForEnc.Text(), pwd)
		res.SetText(fmt.Sprintf("you can find your wallet file at %s", service.BTCPRIVK_PATH))
		privkForEnc.SetText("")
	})
	return encBox
}

func GetBoxForGenePrivkEN(res *ui.MultilineEntry) *ui.Box {
	pParam := ui.NewForm()
	netEntry := ui.NewEntry()
	pParam.Append("Type of Bitcoin (test/main/regtest)：", netEntry, false)
	pButton := ui.NewButton("Get")
	pBox := ui.NewHorizontalBox()
	pBox.Append(ui.NewLabel(""), true)
	pBox.Append(pButton, true)
	pBox.Append(ui.NewLabel(""), true)

	privBox := ui.NewVerticalBox()
	privBox.Append(pParam, false)
	privBox.Append(pBox, false)

	pButton.OnClicked(func(button *ui.Button) {
		res.SetText(service.GetPrivk(netEntry.Text()))
	})
	return privBox
}

func GetBoxForGeneRedeemEN(res *ui.MultilineEntry) *ui.Box {
	rParam := ui.NewForm()
	rnetEntry := ui.NewEntry()
	pubksEntry := ui.NewEntry()
	reqEntry := ui.NewEntry()
	rParam.Append("Type of Bitcoin (test/main/regtest)：", rnetEntry, false)
	rParam.Append("Public Keys：", pubksEntry, false)
	rParam.Append("Require Signature Number：", reqEntry, false)
	rButton := ui.NewButton("Get")
	rBox := ui.NewHorizontalBox()
	rBox.Append(ui.NewLabel(""), true)
	rBox.Append(rButton, true)
	rBox.Append(ui.NewLabel(""), true)

	pubksBox := ui.NewVerticalBox()
	pubksBox.Append(rParam, false)
	pubksBox.Append(rBox, false)

	rButton.OnClicked(func(button *ui.Button) {
		reqNum, _ := strconv.ParseInt(reqEntry.Text(), 10, 64)
		res.SetText(service.GetRedeemForMultiSig(pubksEntry.Text(), rnetEntry.Text(), int(reqNum)))
	})
	return pubksBox
}

func GetBoxForUtxoMonitorEN(res *ui.MultilineEntry) *ui.Box {
	param := ui.NewForm()
	rpc := ui.NewEntry()
	less := ui.NewEntry()
	redeem := ui.NewEntry()
	param.Append("poly RPC-URL：", rpc, false)
	param.Append("Small Value UTXO Limit：", less, false)
	param.Append("Multisig Redeem：", redeem, false)

	status := ui.NewForm()
	sum := ui.NewEntry()
	sum.SetReadOnly(true)
	total := ui.NewEntry()
	total.SetReadOnly(true)
	p2sh := ui.NewEntry()
	p2sh.SetReadOnly(true)
	p2wsh := ui.NewEntry()
	p2wsh.SetReadOnly(true)
	lc := ui.NewEntry()
	lc.SetReadOnly(true)
	fr := ui.NewEntry()
	fr.SetReadOnly(true)
	mc := ui.NewEntry()
	mc.SetReadOnly(true)

	status.Append("UTXO Sum：", sum, false)
	status.Append("UTXO Total Num：", total, false)
	status.Append("P2SH UTXO Num：", p2sh, false)
	status.Append("P2WSH UTXO Num：", p2wsh, false)
	status.Append("Small Value UTXO Num：", lc, false)
	status.Append("Fee Rate：", fr, false)
	status.Append("Min Change：", mc, false)

	start := ui.NewButton("Start Monitor")
	stop := ui.NewButton("Stop Monitor")
	rBox := ui.NewHorizontalBox()
	rBox.Append(ui.NewLabel(""), true)
	rBox.Append(start, true)
	rBox.Append(stop, true)
	rBox.Append(ui.NewLabel(""), true)

	quit := make(chan struct{})
	start.OnClicked(func(button *ui.Button) {
		if !button.Enabled() {
			return
		}
		button.Disable()
		lessPoint, err := strconv.ParseUint(less.Text(), 10, 64)
		if err != nil {
			res.SetText(fmt.Sprintf("wrong small value utxo limit：%v", err))
			button.Enable()
			return
		}
		r, err := hex.DecodeString(redeem.Text())
		if err != nil {
			res.SetText(fmt.Sprintf("wrong multisig redeem：%v", err))
			button.Enable()
			return
		}
		m := service.NewUtxoMonitor(lessPoint, rpc.Text(), r)
		go m.RunMonitor()
		go func() {
			tick := time.NewTicker(time.Second * 5)
			for {
				select {
				case <-tick.C:
					sum.SetText(fmt.Sprintf("%d sat", m.Status.Sum))
					total.SetText(strconv.FormatUint(m.Status.Total, 10))
					p2sh.SetText(fmt.Sprintf("%d(%.3f%%) sum: %d sat", m.Status.P2shNum,
						100*float64(m.Status.P2shNum)/float64(m.Status.Total), m.Status.P2shSum))
					p2wsh.SetText(fmt.Sprintf("%d(%.3f%%) sum: %d sat", m.Status.P2wshNum,
						100*float64(m.Status.P2wshNum)/float64(m.Status.Total), m.Status.P2wshSum))
					lc.SetText(strconv.FormatUint(m.Status.Less, 10))
					fr.SetText(fmt.Sprintf("%d sat/byte", m.Status.FeeRate))
					mc.SetText(fmt.Sprintf("%d sat", m.Status.MinChange))
				case <-quit:
					m.Close()
					button.Enable()
					return
				}
			}
		}()
		res.SetText("you can find all utxo in file ./your_utxo")
	})
	stop.OnClicked(func(button *ui.Button) {
		if start.Enabled() {
			return
		}
		quit <- struct{}{}
	})

	mBox := ui.NewVerticalBox()
	mBox.Append(param, false)
	mBox.Append(rBox, true)
	mBox.Append(status, false)

	return mBox
}
