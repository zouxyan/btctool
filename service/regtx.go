package service

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/polynetwork/poly/common/log"
	"github.com/zouxyan/btctool/builder"
	"github.com/zouxyan/btctool/rest"
)

type RegTxBuilder struct {
	Fee       float64
	Value     float64
	Privkb58  string
	OntAddr   string
	RpcUrl    string
	User      string
	Pwd       string
	ToChainId uint64
	ToAddr    string
	NetParam  *chaincfg.Params
}

func (ra *RegTxBuilder) Run() string {
	if ra.OntAddr == "" {
		log.Error("ont address is required")
		return ""
	}
	if ra.Privkb58 == "" {
		log.Error("privk can't be null")
		return ""
	}

	privk, err := btcutil.DecodeWIF(ra.Privkb58)
	if err != nil {
		log.Fatalf("failed to decode your wif privk %s: %v", err)
		return err.Error()
	}
	addrPubk, err := btcutil.NewAddressPubKey(privk.PrivKey.PubKey().SerializeCompressed(), ra.NetParam)
	if err != nil {
		log.Errorf("Failed to new an address pubkey: %v", err)
		return err.Error()
	}
	pubkScript, err := txscript.PayToAddrScript(addrPubk.AddressPubKeyHash())
	if err != nil {
		log.Errorf("Failed to build pubk script: %v", err)
		return err.Error()
	}

	data, err := buildData(ra.ToChainId, 0, ra.OntAddr)
	if err != nil {
		log.Errorf("Failed to ge data: %v", err)
		return err.Error()
	}

	cli := rest.NewRestCli(ra.RpcUrl, ra.User, ra.Pwd)
	addr := addrPubk.EncodeAddress()
	err = cli.ImportAddress(addr)
	if err != nil {
		log.Errorf("rpc failed: %v", err)
		return err.Error()
	}
	cnt, err := cli.GetBlockCount()
	if err != nil {
		log.Errorf("rpc failed: %v", err)
		return err.Error()
	}
	utxos, err := cli.ListUnspent(6, cnt, addr)
	if err != nil {
		log.Errorf("rpc failed: %v", err)
		return err.Error()
	}
	total, err := btcutil.NewAmount(ra.Value + ra.Fee)
	if err != nil {
		log.Errorf("failed to new amount: %v", err)
		return err.Error()
	}
	selected, sumVal, err := rest.SelectUtxos(utxos, int64(total))
	if err != nil {
		log.Errorf("failed to select utxo: %v", err)
		return err.Error()
	}

	//var prevPkScripts [][]byte
	var ipts []btcjson.TransactionInput
	for _, v := range selected {
		ipts = append(ipts, btcjson.TransactionInput{
			Txid: v.Txid,
			Vout: v.Vout,
		})
	}

	b, err := builder.NewBuilder(&builder.BuildCrossChainTxParam{
		Data:         data,
		Inputs:       ipts,
		NetParam:     ra.NetParam,
		PrevPkScript: pubkScript,
		Privk:        privk.PrivKey,
		Locktime:     nil,
		ToMultiValue: ra.Value,
		Changes: func() map[string]float64 {
			if changeVal := float64(sumVal)/btcutil.SatoshiPerBitcoin - ra.Value - ra.Fee; changeVal > 0 {
				return map[string]float64{addrPubk.EncodeAddress(): changeVal}
			} else {
				return map[string]float64{}
			}
		}(),
		ToAddr: ra.ToAddr,
	})
	if err != nil {
		log.Errorf("Failed to new an instance of Builder: %v", err)
		return err.Error()
	}

	var buf bytes.Buffer
	err = b.BuildSignedTx()
	if err != nil || !b.IsSigned {
		log.Errorf("Failed to build signed transaction: %v", err)
		return err.Error()
	}
	log.Infof("Signed cross chain transaction with your private key")
	err = b.Tx.BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
	if err != nil {
		log.Errorf("Failed to encode transaction: %v", err)
		return err.Error()
	}
	log.Infof("------------------------Your signed cross chain transaction------------------------\n%x\n", buf.Bytes())

	txid, err := cli.SendRawTx(hex.EncodeToString(buf.Bytes()))
	if err != nil {
		log.Errorf("failed to send tx: %v", err)
		return err.Error()
	}
	log.Infof("send tx %s to %s", txid, ra.NetParam.Name)

	return txid
}
