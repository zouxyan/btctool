package btctool

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/Zou-XueYan/btctool/rest"
	"github.com/Zou-XueYan/cctxbuilder/btc"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ontio/multi-chain/common"
	"github.com/ontio/ontology/common/log"
	"os"
	"strconv"
	"strings"
)

type CcTx struct {
	OntAddr        string
	PrevHexTxs     string
	Indexes        string
	AddrScriptHash string
	Privkb58       string
	Value          float64
	Fee            float64
	Addr           string
	SpvAddr        string
	NetType        string
}

func (cctx *CcTx) RunCcTx() {
	if cctx.OntAddr == "" {
		log.Error("ont address is required")
		os.Exit(1)
	}
	if cctx.NetType != "main" && cctx.NetType != "test" {
		log.Errorf("net type is not right: %s", cctx.NetType)
		os.Exit(1)
	}

	var data []byte
	ccflag := byte(0x66)
	chainId := make([]byte, 8)
	binary.BigEndian.PutUint64(chainId, 2)
	ccfee := make([]byte, 8)
	binary.BigEndian.PutUint64(ccfee, 0)
	addrBytes, _ := common.AddressFromBase58(cctx.OntAddr)
	data = append(append(append(append(data, ccflag), chainId...), ccfee...), addrBytes[:]...)

	htxs := strings.Split(cctx.PrevHexTxs, ",")
	if len(htxs) == 0 {
		log.Error("You need to fill the hex of transactions containing your UTXOs in")
		os.Exit(1)
	}
	idxes := strings.Split(cctx.Indexes, ",")
	if len(idxes) == 0 || len(htxs) != len(idxes) {
		log.Errorf("Wrong indexes")
	}

	var idxesNum []uint32
	for _, idx := range idxes {
		num, err := strconv.ParseUint(idx, 10, 32)
		if err != nil {
			log.Error("Failed to parse index %s: %v", idx, err)
			os.Exit(1)
		}
		idxesNum = append(idxesNum, uint32(num))
	}

	var txids []chainhash.Hash
	var prevPkScripts [][]byte
	var amount int64
	for i, htx := range htxs {
		btx, err := hex.DecodeString(htx)
		if err != nil {
			log.Error("Failed to decode hex for no.%d transaction %s: %v", i, htx, err)
			os.Exit(1)
		}
		mtx := wire.NewMsgTx(wire.TxVersion)
		buf := bytes.NewBuffer(btx)
		err = mtx.BtcDecode(buf, wire.ProtocolVersion, wire.LatestEncoding)
		if err != nil {
			log.Error("Failed to decode MsgTx for no.%d transaction %s: %v", i, htx, err)
			os.Exit(1)
		}

		amount += mtx.TxOut[idxesNum[i]].Value
		txids = append(txids, mtx.TxHash())
		prevPkScripts = append(prevPkScripts, mtx.TxOut[idxesNum[i]].PkScript)
	}

	var ipts []btcjson.TransactionInput
	for i, txid := range txids {
		ipts = append(ipts, btcjson.TransactionInput{
			Txid: txid.String(),
			Vout: idxesNum[i],
		})
	}

	builder, err := btc.NewBuilder(&btc.BuildCrossChainTxParam{
		AddrScriptHash: cctx.AddrScriptHash,
		Data:           data,
		Inputs:         ipts,
		NetParam: func(nt string) *chaincfg.Params {
			switch nt {
			case "main":
				return &chaincfg.MainNetParams
			case "test":
				return &chaincfg.TestNet3Params
			default:
				return nil
			}
		}(cctx.NetType),
		PrevPkScripts: prevPkScripts,
		Privk58:       cctx.Privkb58,
		Locktime:      nil,
		ToMultiValue:  cctx.Value,
		Changes: func() map[string]float64 {
			if changeVal := float64(amount)/btcutil.SatoshiPerBitcoin - cctx.Value - cctx.Fee; changeVal > 0 {
				return map[string]float64{cctx.Addr: changeVal}
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
	if cctx.Privkb58 == "" {
		err = builder.Tx.BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
		if err != nil {
			log.Errorf("Failed to encode transaction: %v", err)
			os.Exit(1)
		}
		log.Infof("------------------------Your unsigned cross chain transaction------------------------\n%x\n", buf.Bytes())
		return
	}
	err = builder.BuildSignedTx()
	if err != nil || !builder.IsSigned {
		log.Errorf("Failed to build signed transaction: %v", err)
		os.Exit(1)
	}
	log.Infof("Signed cross chain transaction with your private key")
	err = builder.Tx.BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
	if err != nil {
		log.Errorf("Failed to encode transaction: %v", err)
		os.Exit(1)
	}
	log.Infof("------------------------Your signed cross chain transaction------------------------\n%x\n", buf.Bytes())

	if cctx.SpvAddr == "" {
		log.Infof("spv addr not set, you need to broadcast tx by yourself")
		return
	}

	cli := rest.NewRestCli("", "", "", cctx.SpvAddr)
	err = cli.BroadcastTxBySpv(hex.EncodeToString(buf.Bytes()))
	if err != nil {
		log.Errorf("failed to broadcast tx: %v", err)
	}
	log.Infof("and already broadcast tx %s", builder.Tx.TxHash().String())
}
