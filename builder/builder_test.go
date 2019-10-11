package builder

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ontio/multi-chain/common"
	"os/exec"
	"testing"
)

// fee 2e-6
var prevRawTx = "020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff0502f4010101ffffffff0240be4025000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac0000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000"
var txid = "646ff38cc96e3305247c4a9934f6f8d622248ae350cc670468830670b684d00b"
var addr = "mjEoyyCPsLzJ23xMX6Mti13zMyN36kzn57"
var privkey = "cRRMYvoHPNQu1tCz4ajPxytBVc2SN6GWLAVuyjzm4MVwyqZVrAcX"
var addrScriptHash = "2N5cY8y9RtbbvQRWkX5zAwTPCxSZF9xEj2C"
var addrInOp = "AGjD4Mo25kzcStyh1stp7tXkUuMopD43NT"

func TestNewBuilder(t *testing.T) {
	txb, err := hex.DecodeString(prevRawTx)
	if err != nil {
		t.Fatalf("Failed to decode raw tx: %v", err)
	}

	mtx := wire.NewMsgTx(wire.TxVersion)

	err = mtx.BtcDecode(bytes.NewReader(txb), wire.ProtocolVersion, wire.LatestEncoding)
	if err != nil {
		t.Fatalf("decode : %v", err)
	}
	_, err = NewBuilder(&BuildCrossChainTxParam{
		Inputs: []btcjson.TransactionInput{
			{
				txid,
				1,
			},
		},
		Changes: map[string]float64{
			addr: 0.01 - 4e-5,
		},
		ToMultiValue: 2e-5,
		Locktime:     nil,
		PrevPkScripts: [][]byte{
			mtx.TxOut[0].PkScript,
		},
		Privk58:  privkey,
		NetParam: &chaincfg.TestNet3Params,
	})
	if err != nil {
		t.Fatalf("Failed to build builder: %v", err)
	}
}

func TestBuilder_BuildSignedTx(t *testing.T) {
	flag := byte(0x66)
	chainId := make([]byte, 8)
	binary.BigEndian.PutUint64(chainId, 2)
	fee := make([]byte, 8)
	binary.BigEndian.PutUint64(fee, 1000)
	addrBytes, _ := common.AddressFromBase58(addrInOp) //base58.Decode(addrInOp)

	var data []byte
	data = append(append(append(append(data, flag), chainId...), fee...), addrBytes[:]...)
	if len(data) != 37 {
		t.Fatalf("Wrong length of data")
	}

	txb, err := hex.DecodeString(prevRawTx)
	if err != nil {
		t.Fatalf("Failed to decode raw tx: %v", err)
	}

	mtx := wire.NewMsgTx(wire.TxVersion)

	err = mtx.BtcDecode(bytes.NewReader(txb), wire.ProtocolVersion, wire.LatestEncoding)
	if err != nil {
		t.Fatalf("decode : %v", err)
	}

	b, err := NewBuilder(&BuildCrossChainTxParam{
		Inputs: []btcjson.TransactionInput{
			{
				txid,
				2,
			},
		},
		Changes: map[string]float64{
			addr: 0.00988,
		},
		ToMultiValue: 0.00188,
		Locktime:     nil,
		PrevPkScripts: [][]byte{
			mtx.TxOut[2].PkScript,
		},
		AddrScriptHash: addrScriptHash,
		Privk58:        privkey,
		NetParam:       &chaincfg.TestNet3Params,
		Data:           data,
	})
	if err != nil {
		t.Fatalf("Failed to get a builder: %v", err)
	}

	fmt.Println(len(b.Privks))

	//err = b.BuildMultiSignedTx()
	err = b.BuildSignedTx()
	if err != nil || !b.IsSigned {
		t.Fatalf("Failed to build a signed transaction: %v", err)
	}

	//ss, _ := txscript.DisasmString(b.Tx.TxIn[0].SignatureScript)
	//fmt.Println(ss)
	//rr, _ := hex.DecodeString(strings.Split(ss, " ")[6])
	//rss, _ := txscript.DisasmString(rr)
	//fmt.Println(rss)

	satoshi, err := btcutil.NewAmount(0.002)
	if err != nil {
		t.Fatalf("Failed to transfer amount to satoshi")
	}

	eng, err := txscript.NewEngine(mtx.TxOut[2].PkScript, b.Tx, 0, txscript.ScriptVerifyCleanStack|txscript.ScriptBip16, nil,
		nil, int64(satoshi))
	if err != nil {
		t.Fatalf("Failed to start engine: %v", err)
	}
	ls, err := txscript.DisasmString(mtx.TxOut[2].PkScript)
	if err != nil {
		t.Fatalf("spk err: %v", err)
	}
	fmt.Printf("locking script is %s\n", ls)
	res := eng.Execute()
	if res != nil {
		t.Fatalf("Failed to excute: %v", res)
	}
	var buf bytes.Buffer
	err = b.Tx.BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
	if err != nil {
		t.Fatalf("Failed to encode tx: %v", err)
	}
	fmt.Printf("serialized tx: %x\n", buf.Bytes())
}

func TestSth(t *testing.T) {
	//addr, err := btcutil.DecodeAddress("2N5cY8y9RtbbvQRWkX5zAwTPCxSZF9xEj2C", &chaincfg.TestNet3Params)
	//if err != nil {
	//	t.Fatalf("failed to : %v", err)
	//}
	//p2shScript, err := txscript.PayToAddrScript(addr)
	//if err != nil {
	//	t.Fatalf("failed to : %v", err)
	//}
	//
	//str, _ := txscript.DisasmString(p2shScript)
	//fmt.Println(str)
	cmd := exec.Command("ls", "-l")
	cmd.Run()
	fmt.Printf("cmd is : %v\n", cmd)
}
