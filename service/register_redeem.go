package service

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/btcsuite/btcutil"
	"github.com/ontio/btc-spvclient/utils"
	sdk "github.com/ontio/multi-chain-go-sdk"
	"github.com/ontio/multi-chain/common/log"
	"os"
	"strings"
)

func RedeemRegister(alliaRpc, contractAddr, redeem, sigs, walletFile, pwd string, contractId uint64) string {
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

	txHash, err := allia.Native.Scm.RegisterRedeem(1, contractId, r, c, "", sigArr, acct)
	if err != nil {
		log.Fatalf("failed to register: %v", err)
		os.Exit(1)
	}
	log.Infof("successful to register call, tx hash is %s", txHash.ToHexString())
	return txHash.ToHexString()
}

func GetSigForRedeemContract(contract, redeem, privk string, cver uint64) string {
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
	cverBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(cverBytes, cver)

	hash := btcutil.Hash160(append(append(r, c...), cverBytes...))
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

func Test(alliaRpc, walletFile, pwd string) string {
	allia := sdk.NewMultiChainSdk()
	allia.NewRpcClient().SetAddress(alliaRpc)
	acct, err := utils.GetAccountByPassword(allia, walletFile, pwd)
	if err != nil {
		log.Fatalf("failed to get account: %v", err)
		os.Exit(1)
	}
	txHash, err := allia.Native.Scm.RegisterSideChain("123", 997, 997, "TEST", 100, acct)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Infof("successful to register call, tx hash is %s", txHash.ToHexString())
	return txHash.ToHexString()
}
