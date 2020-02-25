package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"strings"
)

func GetPrivk(netType string) string {
	priv, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return err.Error()
	}
	net := func(nt string) *chaincfg.Params {
		switch nt {
		case "main":
			return &chaincfg.MainNetParams
		case "test":
			return &chaincfg.TestNet3Params
		default:
			return nil
		}
	}(netType)
	wif, err := btcutil.NewWIF(priv, net, true)
	if err != nil {
		return err.Error()
	}

	addr, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), net)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("your private key is %s\nyour compressed public key is %s\nyour legacy address is %s", wif.String(),
		hex.EncodeToString(wif.PrivKey.PubKey().SerializeCompressed()), addr.EncodeAddress())
}

func GetRedeemForMultiSig(pubks, netType string, req int) string {
	net := func(nt string) *chaincfg.Params {
		switch nt {
		case "main":
			return &chaincfg.MainNetParams
		case "test":
			return &chaincfg.TestNet3Params
		default:
			return nil
		}
	}(netType)

	sArr := strings.Split(pubks, ",")
	pArr := make([]*btcutil.AddressPubKey, len(sArr))
	for i, v := range sArr {
		p, err := hex.DecodeString(v)
		if err != nil {
			return err.Error()
		}

		pkAddr, err := btcutil.NewAddressPubKey(p, net)
		if err != nil {
			return err.Error()
		}
		pArr[i] = pkAddr
	}

	redeem, err := txscript.MultiSigScript(pArr, req)
	if err != nil {
		return fmt.Sprintf("failed to get redeem: %v", err)
	}

	p2sh, err := btcutil.NewAddressScriptHash(redeem, net)
	if err != nil {
		return err.Error()
	}
	hasher := sha256.New()
	hasher.Write(redeem)
	p2wsh, err := btcutil.NewAddressWitnessScriptHash(hasher.Sum(nil), net)
	return fmt.Sprintf("your redeem is %s\nyour P2SH address is %s\nyour P2WSH address is %s",
		hex.EncodeToString(redeem), p2sh.EncodeAddress(), p2wsh.EncodeAddress())
}