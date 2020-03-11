package service

import (
	"encoding/hex"
	"github.com/btcsuite/btcutil"
	"github.com/ontio/btc-spvclient/utils"
	sdk "github.com/ontio/multi-chain-go-sdk"
	"github.com/ontio/multi-chain/common/log"
	"github.com/ontio/multi-chain/native/service/cross_chain_manager/btc"
	mutils "github.com/ontio/multi-chain/native/service/utils"
	"os"
	"strings"
)

func RedeemRegister(alliaRpc, contractAddr, redeem, sigs, walletFile, pwd string, contractId, cver uint64) string {
	ss := strings.Split(sigs, ",")
	sigArr := make([][]byte, len(ss))
	for i, s := range ss {
		sb, err := hex.DecodeString(s)
		if err != nil {
			log.Fatalf("failed to decode no.%d sig %s: %v", i, s, err)
			os.Exit(1)
		}
		sigArr[i] = sb
	}
	allia := sdk.NewMultiChainSdk()
	allia.NewRpcClient().SetAddress(alliaRpc)
	r, err := hex.DecodeString(redeem)
	if err != nil {
		log.Fatalf("failed to decode redeem: %v", err)
		os.Exit(1)
	}
	contractAddr = strings.Replace(contractAddr, "0x", "", 1)
	c, err := hex.DecodeString(contractAddr)
	if err != nil {
		log.Fatalf("failed to decode contract: %v", err)
		os.Exit(1)
	}
	acct, err := utils.GetAccountByPassword(allia, walletFile, pwd)
	if err != nil {
		log.Fatalf("failed to get account: %v", err)
		os.Exit(1)
	}

	txHash, err := allia.Native.Scm.RegisterRedeem(1, contractId, cver, r, c, sigArr, acct)
	if err != nil {
		log.Fatalf("failed to register: %v", err)
		os.Exit(1)
	}
	log.Infof("successful to register call, tx hash is %s", txHash.ToHexString())
	return txHash.ToHexString()
}

func GetSigForRedeemContract(contract, redeem, privk string, cver, toChainId uint64) string {
	r, err := hex.DecodeString(redeem)
	if err != nil {
		log.Fatalf("failed to decode redeem: %v", err)
		os.Exit(1)
	}
	contract = strings.Replace(contract, "0x", "", 1)
	c, err := hex.DecodeString(contract)
	if err != nil {
		log.Fatalf("failed to decode contract: %v", err)
		os.Exit(1)
	}

	fromChainId := mutils.GetUint64Bytes(btc.BTC_CHAIN_ID)
	toChainIdBytes := mutils.GetUint64Bytes(toChainId)
	cverBytes := mutils.GetUint64Bytes(cver)
	hash := btcutil.Hash160(append(append(append(append(r, fromChainId...), c...), toChainIdBytes...), cverBytes...))
	wif, err := btcutil.DecodeWIF(privk)
	if err != nil {
		log.Fatalf("failed to decode wif: %v", err)
		os.Exit(1)
	}

	sig, err := wif.PrivKey.Sign(hash)
	res := hex.EncodeToString(sig.Serialize())
	log.Infof("your sig is %s", res)
	return res
}

func SetBtcTxParam(alliaRpc, redeem, sigs, walletFile, pwd string, fr, mc, pver uint64) string {
	ss := strings.Split(sigs, ",")
	sigArr := make([][]byte, len(ss))
	for i, s := range ss {
		sb, err := hex.DecodeString(s)
		if err != nil {
			log.Fatalf("failed to decode no.%d sig %s: %v", i, s, err)
			os.Exit(1)
		}
		sigArr[i] = sb
	}
	allia := sdk.NewMultiChainSdk()
	allia.NewRpcClient().SetAddress(alliaRpc)
	r, err := hex.DecodeString(redeem)
	if err != nil {
		log.Fatalf("failed to decode redeem: %v", err)
		os.Exit(1)
	}

	acct, err := utils.GetAccountByPassword(allia, walletFile, pwd)
	if err != nil {
		log.Fatalf("failed to get account: %v", err)
		os.Exit(1)
	}

	txHash, err := allia.Native.Scm.SetBtcTxParam(r, btc.BTC_CHAIN_ID, fr, mc, pver, sigArr, acct)
	if err != nil {
		log.Fatalf("failed to set btc tx param: %v", err)
		os.Exit(1)
	}
	log.Infof("successful to set btc tx param, tx hash is %s", txHash.ToHexString())
	return txHash.ToHexString()
}

func GetSigForBtcTxParam(fr, mc, pver uint64, redeem, privk string) string {
	r, err := hex.DecodeString(redeem)
	if err != nil {
		log.Fatalf("failed to decode redeem: %v", err)
		os.Exit(1)
	}
	fromChainId := mutils.GetUint64Bytes(btc.BTC_CHAIN_ID)
	frBytes := mutils.GetUint64Bytes(fr)
	mcBytes := mutils.GetUint64Bytes(mc)
	verBytes := mutils.GetUint64Bytes(pver)
	hash := btcutil.Hash160(append(append(append(append(r, fromChainId...), frBytes...), mcBytes...), verBytes...))
	wif, err := btcutil.DecodeWIF(privk)
	if err != nil {
		log.Fatalf("failed to decode wif: %v", err)
		os.Exit(1)
	}
	sig, err := wif.PrivKey.Sign(hash)
	res := hex.EncodeToString(sig.Serialize())
	log.Infof("your sig is %s", res)
	return res
}
