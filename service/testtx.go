package service

import (
	"bytes"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ontio/ontology/common/log"
	"github.com/zouxyan/btctool/builder"
	"strconv"
	"strings"
)

type TestTxBuilder struct {
	OntAddr   string
	Txids     string
	Indexes   string
	Privkb58  string
	Value     float64
	Fee       float64
	NetType   string
	Vals      []float64
	ToAddr    string
	ToChainId uint64
}

func (cctx *TestTxBuilder) Run() *wire.MsgTx {
	if cctx.OntAddr == "" {
		log.Error("ont address is required")
		return nil
	}
	if cctx.NetType != "main" && cctx.NetType != "test" {
		log.Errorf("net type is not right: %s", cctx.NetType)
		return nil
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
	}(cctx.NetType)

	privk, err := btcutil.DecodeWIF(cctx.Privkb58)
	if err != nil {
		log.Fatalf("failed to decode your wif privk %s: %v", err)
		return nil
	}
	addrPubk, err := btcutil.NewAddressPubKey(privk.PrivKey.PubKey().SerializeCompressed(), net)
	if err != nil {
		log.Errorf("Failed to new an address pubkey: %v", err)
		return nil
	}
	pubkScript, err := txscript.PayToAddrScript(addrPubk.AddressPubKeyHash())
	if err != nil {
		log.Errorf("Failed to build pubk script: %v", err)
		return nil
	}

	data, err := buildData(cctx.ToChainId, 0, cctx.OntAddr)
	if err != nil {
		log.Errorf("failed to build data: %v", err)
		return nil
	}

	txidArr := strings.Split(cctx.Txids, ",")
	if len(txidArr) == 0 {
		log.Error("You need to fill the txids of transactions containing your UTXOs in")
		return nil
	}
	idxes := strings.Split(cctx.Indexes, ",")
	if len(txidArr) != len(idxes) {
		log.Errorf("Wrong indexes")
		return nil
	}

	var idxesNum []uint32
	for _, idx := range idxes {
		num, err := strconv.ParseUint(idx, 10, 32)
		if err != nil {
			log.Errorf("Failed to parse index %s: %v", idx, err)
			return nil
		}
		idxesNum = append(idxesNum, uint32(num))
	}

	if len(cctx.Vals) != len(txidArr) {
		log.Errorf("Wrong vals")
		return nil
	}
	var amount float64
	for _, val := range cctx.Vals {
		amount += val
	}

	var ipts []btcjson.TransactionInput
	for i, txid := range txidArr {
		ipts = append(ipts, btcjson.TransactionInput{
			Txid: txid,
			Vout: idxesNum[i],
		})
	}

	b, err := builder.NewBuilder(&builder.BuildCrossChainTxParam{
		Data:         data,
		Inputs:       ipts,
		NetParam:     net,
		PrevPkScript: pubkScript,
		Privk:        privk.PrivKey,
		Locktime:     nil,
		ToMultiValue: cctx.Value,
		Changes: func() map[string]float64 {
			if changeVal := amount - cctx.Value - cctx.Fee; changeVal > 0 {
				return map[string]float64{addrPubk.EncodeAddress(): changeVal}
			} else {
				return map[string]float64{}
			}
		}(),
		ToAddr: cctx.ToAddr,
	})
	if err != nil {
		log.Errorf("Failed to new an instance of Builder: %v", err)
		return nil
	}
	var buf bytes.Buffer
	if cctx.Privkb58 == "" {
		err = b.Tx.BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
		if err != nil {
			log.Errorf("Failed to encode transaction: %v", err)
			return nil
		}
		log.Infof("------------------------Your unsigned cross chain transaction------------------------\n%x\n", buf.Bytes())
		return nil
	}
	err = b.BuildSignedTx()
	if err != nil || !b.IsSigned {
		log.Errorf("Failed to build signed transaction: %v", err)
		return nil
	}
	log.Infof("Signed cross chain transaction with your private key")
	err = b.Tx.BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
	if err != nil {
		log.Errorf("Failed to encode transaction: %v", err)
		return nil
	}
	log.Infof("------------------------Your signed cross chain transaction------------------------\n%x\n", buf.Bytes())
	log.Infof("you need to broadcast tx by yourself")
	return b.Tx
}
