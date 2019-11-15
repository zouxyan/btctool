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
	"github.com/ontio/ontology/common/log"
	"os"
	"strconv"
	"strings"
)

type CcTx struct {
	OntAddr        string
	Txids          string
	Indexes        string
	Privkb58       string
	Value          float64
	Fee            float64
	SpvAddr        string
	NetType        string
	Vals           []float64
	ContractAddr   string
	Redeem string
	IsSegWit int
}

func (cctx *CcTx) RunCcTx() *wire.MsgTx {
	if cctx.OntAddr == "" {
		log.Error("ont address is required")
		os.Exit(1)
	}
	if cctx.NetType != "main" && cctx.NetType != "test" {
		log.Errorf("net type is not right: %s", cctx.NetType)
		os.Exit(1)
	}
	if cctx.ContractAddr == "" {
		log.Error("contract address can't be null")
		os.Exit(1)
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
	privkey := base58.Decode(cctx.Privkb58)
	privk, pubk := btcec.PrivKeyFromBytes(btcec.S256(), privkey)
	addrPubk, err := btcutil.NewAddressPubKey(pubk.SerializeCompressed(), net)
	if err != nil {
		log.Errorf("Failed to new an address pubkey: %v", err)
		os.Exit(1)
	}
	pubkScript, err := txscript.PayToAddrScript(addrPubk.AddressPubKeyHash())
	if err != nil {
		log.Errorf("Failed to build pubk script: %v", err)
		os.Exit(1)
	}

	data, err := buildData(2, 0, cctx.OntAddr, cctx.ContractAddr)
	if err != nil {
		log.Errorf("failed to build data: %v", err)
		os.Exit(1)
	}

	txidArr := strings.Split(cctx.Txids, ",")
	if len(txidArr) == 0 {
		log.Error("You need to fill the txids of transactions containing your UTXOs in")
		os.Exit(1)
	}
	idxes := strings.Split(cctx.Indexes, ",")
	if len(txidArr) != len(idxes) {
		log.Errorf("Wrong indexes")
		os.Exit(1)
	}

	var idxesNum []uint32
	for _, idx := range idxes {
		num, err := strconv.ParseUint(idx, 10, 32)
		if err != nil {
			log.Errorf("Failed to parse index %s: %v", idx, err)
			os.Exit(1)
		}
		idxesNum = append(idxesNum, uint32(num))
	}

	if len(cctx.Vals) != len(txidArr) {
		log.Errorf("Wrong vals")
		os.Exit(1)
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
		Data:           data,
		Inputs:         ipts,
		NetParam:       net,
		PrevPkScript:   pubkScript,
		Privk:          privk,
		Locktime:       nil,
		ToMultiValue:   cctx.Value,
		Changes: func() map[string]float64 {
			if changeVal := amount - cctx.Value - cctx.Fee; changeVal > 0 {
				return map[string]float64{addrPubk.EncodeAddress(): changeVal}
			} else {
				return map[string]float64{}
			}
		}(),
		Redeem: cctx.Redeem,
		IsSegWit: cctx.IsSegWit,
	})
	if err != nil {
		log.Errorf("Failed to new an instance of Builder: %v", err)
		os.Exit(1)
	}
	var buf bytes.Buffer
	if cctx.Privkb58 == "" {
		err = b.Tx.BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
		if err != nil {
			log.Errorf("Failed to encode transaction: %v", err)
			os.Exit(1)
		}
		log.Infof("------------------------Your unsigned cross chain transaction------------------------\n%x\n", buf.Bytes())
		return nil
	}
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

	if cctx.SpvAddr == "" {
		log.Infof("spv addr not set, you need to broadcast tx by yourself")
		return nil
	}

	cli := rest.NewRestCli("", "", "", cctx.SpvAddr)
	err = cli.BroadcastTxBySpv(hex.EncodeToString(buf.Bytes()))
	if err != nil {
		log.Errorf("failed to broadcast tx: %v", err)
	}
	log.Infof("and already broadcast tx %s", b.Tx.TxHash().String())

	return b.Tx
}
