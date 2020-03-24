package service

import (
	"encoding/hex"
	"github.com/btcsuite/btcutil"
	"github.com/ontio/btc-spvclient/utils"
	sdk "github.com/ontio/multi-chain-go-sdk"
	"github.com/ontio/multi-chain/common/log"
	"github.com/ontio/multi-chain/native/service/governance/side_chain_manager"
	mutils "github.com/ontio/multi-chain/native/service/utils"
	"github.com/ontio/ontology/common"
	"strings"
)

func RedeemRegister(alliaRpc, contractAddr, redeem, sigs, walletFile, pwd string, contractId, cver uint64) string {
	ss := strings.Split(sigs, ",")
	sigArr := make([][]byte, len(ss))
	for i, s := range ss {
		sb, err := hex.DecodeString(s)
		if err != nil {
			log.Fatalf("failed to decode no.%d sig %s: %v", i, s, err)
			return "" 
		}
		sigArr[i] = sb
	}
	allia := sdk.NewMultiChainSdk()
	allia.NewRpcClient().SetAddress(alliaRpc)
	r, err := hex.DecodeString(redeem)
	if err != nil {
		log.Fatalf("failed to decode redeem: %v", err)
		return ""
	}
	contractAddr = strings.Replace(contractAddr, "0x", "", 1)
	var c []byte
	if contractId == 3 {
		addr, err := common.AddressFromHexString(contractAddr)
		if err != nil {
			log.Fatalf("")
			return ""
		}
		c = addr[:]
	} else {
		c, err = hex.DecodeString(contractAddr)
		if err != nil {
			log.Fatalf("failed to decode contract: %v", err)
			return ""
		}
	}

	acct, err := utils.GetAccountByPassword(allia, walletFile, []byte(pwd))
	if err != nil {
		log.Fatalf("failed to get account: %v", err)
		return ""
	}

	txHash, err := allia.Native.Scm.RegisterRedeem(1, contractId, r, c, cver, sigArr, acct)
	if err != nil {
		log.Fatalf("failed to register: %v", err)
		return ""
	}
	log.Infof("successful to register call, tx hash is %s", txHash.ToHexString())
	return txHash.ToHexString()
}

func GetSigForRedeemContract(contract, redeem, privk string, cver, toChainId uint64) string {
	r, err := hex.DecodeString(redeem)
	if err != nil {
		log.Fatalf("failed to decode redeem: %v", err)
		return ""
	}
	contract = strings.Replace(contract, "0x", "", 1)
	var c []byte
	if toChainId == 3 {
		addr, err := common.AddressFromHexString(contract)
		if err != nil {
			log.Fatalf("")
			return ""
		}
		c = addr[:]
	} else {
		c, err = hex.DecodeString(contract)
		if err != nil {
			log.Fatalf("failed to decode contract: %v", err)
			return ""
		}
	}

	fromChainId := mutils.GetUint64Bytes(side_chain_manager.BTC_CHAIN_ID)
	toChainIdBytes := mutils.GetUint64Bytes(toChainId)
	cverBytes := mutils.GetUint64Bytes(cver)
	hash := btcutil.Hash160(append(append(append(append(r, fromChainId...), c...), toChainIdBytes...), cverBytes...))
	wif, err := btcutil.DecodeWIF(privk)
	if err != nil {
		log.Fatalf("failed to decode wif: %v", err)
		return ""
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
			return ""
		}
		sigArr[i] = sb
	}
	allia := sdk.NewMultiChainSdk()
	allia.NewRpcClient().SetAddress(alliaRpc)
	r, err := hex.DecodeString(redeem)
	if err != nil {
		log.Fatalf("failed to decode redeem: %v", err)
		return ""
	}

	acct, err := utils.GetAccountByPassword(allia, walletFile, []byte(pwd))
	if err != nil {
		log.Fatalf("failed to get account: %v", err)
		return ""
	}

	txHash, err := allia.Native.Scm.SetBtcTxParam(r, side_chain_manager.BTC_CHAIN_ID, fr, mc, pver, sigArr, acct)
	if err != nil {
		log.Fatalf("failed to set btc tx param: %v", err)
		return ""
	}
	log.Infof("successful to set btc tx param, tx hash is %s", txHash.ToHexString())
	return txHash.ToHexString()
}

func GetSigForBtcTxParam(fr, mc, pver uint64, redeem, privk string) string {
	r, err := hex.DecodeString(redeem)
	if err != nil {
		log.Fatalf("failed to decode redeem: %v", err)
		return ""
	}
	fromChainId := mutils.GetUint64Bytes(side_chain_manager.BTC_CHAIN_ID)
	frBytes := mutils.GetUint64Bytes(fr)
	mcBytes := mutils.GetUint64Bytes(mc)
	verBytes := mutils.GetUint64Bytes(pver)
	hash := btcutil.Hash160(append(append(append(append(r, fromChainId...), frBytes...), mcBytes...), verBytes...))
	wif, err := btcutil.DecodeWIF(privk)
	if err != nil {
		log.Fatalf("failed to decode wif: %v", err)
		return ""
	}
	sig, err := wif.PrivKey.Sign(hash)
	res := hex.EncodeToString(sig.Serialize())
	log.Infof("your sig is %s", res)
	return res
}
