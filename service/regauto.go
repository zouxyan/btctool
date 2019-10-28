package service

import (
	"bytes"
	"encoding/hex"
	"github.com/Zou-XueYan/btctool/builder"
	"github.com/Zou-XueYan/btctool/rest"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ontio/multi-chain/common/log"
	"os"
)

type RegAuto struct {
	Fee            float64
	Value          float64
	Privkb58       string
	OntAddr        string
	RpcUrl         string
	User           string
	Pwd            string
	AddrScriptHash string
	ContractAddr   string
}

func (ra *RegAuto) RunRegAuto() {
	if ra.OntAddr == "" {
		log.Error("ont address is required")
		os.Exit(1)
	}
	if ra.Privkb58 == "" {
		log.Error("privk can't be null")
		os.Exit(1)
	}

	privkey := base58.Decode(ra.Privkb58)
	privk, pubk := btcec.PrivKeyFromBytes(btcec.S256(), privkey)
	addrPubk, err := btcutil.NewAddressPubKey(pubk.SerializeCompressed(), &chaincfg.RegressionNetParams)
	if err != nil {
		log.Errorf("Failed to new an address pubkey: %v", err)
		os.Exit(1)
	}
	pubkScript, err := txscript.PayToAddrScript(addrPubk.AddressPubKeyHash())
	if err != nil {
		log.Errorf("Failed to build pubk script: %v", err)
		os.Exit(1)
	}

	data, err := buildData(2, 0, ra.OntAddr, ra.ContractAddr)
	if err != nil {
		log.Errorf("Failed to ge data: %v", err)
		os.Exit(1)
	}

	cli := rest.NewRestCli(ra.RpcUrl, ra.User, ra.Pwd, "")
	addr := addrPubk.EncodeAddress()
	err = cli.ImportAddress(addr)
	if err != nil {
		log.Errorf("rpc failed: %v", err)
		os.Exit(1)
	}
	cnt, err := cli.GetBlockCount()
	if err != nil {
		log.Errorf("rpc failed: %v", err)
		os.Exit(1)
	}
	utxos, err := cli.ListUnspent(6, cnt, addr)
	if err != nil {
		log.Errorf("rpc failed: %v", err)
		os.Exit(1)
	}
	total, err := btcutil.NewAmount(ra.Value + ra.Fee)
	if err != nil {
		log.Errorf("failed to new amount: %v", err)
		os.Exit(1)
	}
	selected, sumVal, err := rest.SelectUtxos(utxos, int64(total))
	if err != nil {
		log.Errorf("failed to select utxo: %v", err)
		os.Exit(1)
	}

	//var prevPkScripts [][]byte
	var ipts []btcjson.TransactionInput
	for _, v := range selected {
		ipts = append(ipts, btcjson.TransactionInput{
			Txid: v.Txid,
			Vout: v.Vout,
		})
		//sb, err := hex.DecodeString(v.ScriptPubKey)
		//if err != nil {
		//	log.Errorf("failed to decode hex string pubk %s: %v", err)
		//	os.Exit(1)
		//}
		//prevPkScripts = append(prevPkScripts, sb)
	}

	b, err := builder.NewBuilder(&builder.BuildCrossChainTxParam{
		AddrScriptHash: ra.AddrScriptHash,
		Data:           data,
		Inputs:         ipts,
		NetParam:       &chaincfg.RegressionNetParams,
		PrevPkScript:   pubkScript,
		Privk:          privk,
		Locktime:       nil,
		ToMultiValue:   ra.Value,
		Changes: func() map[string]float64 {
			if changeVal := float64(sumVal)/btcutil.SatoshiPerBitcoin - ra.Value - ra.Fee; changeVal > 0 {
				return map[string]float64{addrPubk.EncodeAddress(): changeVal}
			} else {
				return map[string]float64{}
			}
		}(),
	})
	if err != nil {
		log.Errorf("Failed to new an instance of Builder: %v", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	//if ra.Privkb58 == "" {
	//	err = b.Tx.BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
	//	if err != nil {
	//		log.Errorf("Failed to encode transaction: %v", err)
	//		os.Exit(1)
	//	}
	//	log.Infof("------------------------Your unsigned cross chain transaction------------------------\n%x\n", buf.Bytes())
	//	return
	//}
	err = b.BuildSignedTx()
	if err != nil || !b.IsSigned {
		log.Errorf("Failed to build signed transaction: %v", err)
		os.Exit(1)
	}
	log.Infof("Signed cross chain transaction with your private key")
	err = b.Tx.BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
	if err != nil {
		log.Errorf("Failed to encode transaction: %v", err)
		os.Exit(1)
	}
	log.Infof("------------------------Your signed cross chain transaction------------------------\n%x\n", buf.Bytes())

	txid, err := cli.SendRawTx(hex.EncodeToString(buf.Bytes()))
	if err != nil {
		log.Errorf("failed to send tx: %v", err)
		os.Exit(1)
	}
	log.Infof("send tx %s to regression net", txid)
}
