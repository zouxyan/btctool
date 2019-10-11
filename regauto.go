package btctool

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/Zou-XueYan/btctool/rest"
	"github.com/Zou-XueYan/cctxbuilder/btc"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ontio/multi-chain/common/log"
	"github.com/ontio/ontology/common"
	"os"
)

type RegAuto struct {
	Addr           string
	Fee            float64
	Value          float64
	Privkb58       string
	OntAddr        string
	RpcUrl         string
	User           string
	Pwd            string
	AddrScriptHash string
}

func (ra *RegAuto) RunRegAuto() {
	if ra.OntAddr == "" {
		log.Error("ont address is required")
		os.Exit(1)
	}

	var data []byte
	ccflag := byte(0x66)
	chainId := make([]byte, 8)
	binary.BigEndian.PutUint64(chainId, 2)
	ccfee := make([]byte, 8)
	binary.BigEndian.PutUint64(ccfee, 0)
	addrBytes, _ := common.AddressFromBase58(ra.OntAddr)
	data = append(append(append(append(data, ccflag), chainId...), ccfee...), addrBytes[:]...)

	cli := rest.NewRestCli(ra.RpcUrl, ra.User, ra.Pwd, "")
	err := cli.ImportAddress(ra.Addr)
	if err != nil {
		log.Errorf("rpc failed: %v", err)
		os.Exit(1)
	}
	cnt, err := cli.GetBlockCount()
	if err != nil {
		log.Errorf("rpc failed: %v", err)
		os.Exit(1)
	}
	utxos, err := cli.ListUnspent(6, cnt, ra.Addr)
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

	var prevPkScripts [][]byte
	var ipts []btcjson.TransactionInput
	for _, v := range selected {
		ipts = append(ipts, btcjson.TransactionInput{
			Txid: v.Txid,
			Vout: v.Vout,
		})
		sb, err := hex.DecodeString(v.ScriptPubKey)
		if err != nil {
			log.Errorf("failed to decode hex string pubk %s: %v", err)
			os.Exit(1)
		}
		prevPkScripts = append(prevPkScripts, sb)
	}

	builder, err := btc.NewBuilder(&btc.BuildCrossChainTxParam{
		AddrScriptHash: ra.AddrScriptHash,
		Data:           data,
		Inputs:         ipts,
		NetParam:       &chaincfg.RegressionNetParams,
		PrevPkScripts:  prevPkScripts,
		Privk58:        ra.Privkb58,
		Locktime:       nil,
		ToMultiValue:   ra.Value,
		Changes: func() map[string]float64 {
			if changeVal := float64(sumVal)/btcutil.SatoshiPerBitcoin - ra.Value - ra.Fee; changeVal > 0 {
				return map[string]float64{ra.Addr: changeVal}
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
	if ra.Privkb58 == "" {
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

	txid, err := cli.SendRawTx(hex.EncodeToString(buf.Bytes()))
	if err != nil {
		log.Errorf("failed to send tx: %v", err)
		os.Exit(1)
	}
	log.Infof("send tx %s to regression net", txid)
}
