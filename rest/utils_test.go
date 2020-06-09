package rest

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"testing"
)

func TestSelectUtxos(t *testing.T) {
	val := 101
	cli := NewRestCli("http://172.168.3.77:18443", "test", "test")
	utxos, err := cli.ListUnspent(1, 710, "mjEoyyCPsLzJ23xMX6Mti13zMyN36kzn57")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	selected, sum, err := SelectUtxos(utxos, int64(val*btcutil.SatoshiPerBitcoin))
	if err != nil {
		t.Fatalf("failed to select: %v", err)
	}

	for _, v := range selected {
		fmt.Printf("%v", v)
	}
	fmt.Printf("sum is %d, val is %d\n", sum, int64(val*btcutil.SatoshiPerBitcoin))
}

func TestSth(t *testing.T) {

	cli := NewRestCli("http://172.168.3.10:20336", "test", "test")
	utxos, err := cli.ListUnspent(1, 1723442, "tb1qy94qnjuwu5w6r2g74z2z25khjdkgs6ssk5rjnyqrvcvpds8f7x9shrfspn")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	selected, sum, err := SelectUtxos(utxos, int64(0.4*btcutil.SatoshiPerBitcoin))
	if err != nil {
		t.Fatalf("failed to select: %v", err)
	}
	mtx := wire.NewMsgTx(wire.TxVersion)
	ins := make([]*wire.TxIn, len(selected))
	for i, u := range selected {
		txid, _ := chainhash.NewHashFromStr(u.Txid)
		ins[i] = wire.NewTxIn(wire.NewOutPoint(txid, u.Vout), nil, nil)
	}
	mtx.TxIn = ins

	to, _ := btcutil.DecodeAddress("mpCNjy4QYAmw8eumHJRbVtt6bMDVQvPpFn", &chaincfg.TestNet3Params)
	p2pkh, _ := txscript.PayToAddrScript(to)
	out1 := wire.NewTxOut(int64(0.4*btcutil.SatoshiPerBitcoin) - 60000, p2pkh)
	from, _ := btcutil.DecodeAddress("tb1qy94qnjuwu5w6r2g74z2z25khjdkgs6ssk5rjnyqrvcvpds8f7x9shrfspn",
		&chaincfg.TestNet3Params)
	p2wsh, _ := txscript.PayToAddrScript(from)
	out2 := wire.NewTxOut(sum - int64(0.4*btcutil.SatoshiPerBitcoin), p2wsh)
	mtx.AddTxOut(out1)
	mtx.AddTxOut(out2)

	privks := make([]*btcutil.WIF, 5)
	privks[0], _ = btcutil.DecodeWIF("cTqbqa1YqCf4BaQTwYDGsPAB4VmWKUU67G5S1EtrHSWNRwY6QSag")
	privks[1], _ = btcutil.DecodeWIF("cT2HP4QvL8c6otn4LrzUWzgMBfTo1gzV2aobN1cTiuHPXH9Jk2ua")
	privks[2], _ = btcutil.DecodeWIF("cSQmGg6spbhd23jHQ9HAtz3XU7GYJjYaBmFLWHbyKa9mWzTxEY5A")
	privks[3], _ = btcutil.DecodeWIF("cPYAx61EjwshK5SQ6fqH7QGjc8L48xiJV7VRGpYzPSbkkZqrzQ5b")
	privks[4], _ = btcutil.DecodeWIF("cVV9UmtnnhebmSQgHhbDZWCb7zBHbiAGDB9a5M2ffe1WpqvwD5zg")

	rdm, _ := hex.DecodeString("552102dec9a415b6384ec0a9331d0cdf02020f0f1e5731c327b86e2b5a92455a289748210365b1066bcfa21987c3e207b92e309b95ca6bee5f1133cf04d6ed4ed265eafdbc21031104e387cd1a103c27fdc8a52d5c68dec25ddfb2f574fbdca405edfd8c5187de21031fdb4b44a9f20883aff505009ebc18702774c105cb04b1eecebcb294d404b1cb210387cda955196cc2b2fc0adbbbac1776f8de77b563c6d2a06a77d96457dc3d0d1f2102dd7767b6a7cc83693343ba721e0f5f4c7b4b8d85eeb7aec20d227625ec0f59d321034ad129efdab75061e8d4def08f5911495af2dae6d3e9a4b6e7aeb5186fa432fc57ae")

	sh := txscript.NewTxSigHashes(mtx)
	for i, in := range mtx.TxIn {
		data := make([][]byte, 7)
		for j, p := range privks {
			sig, err := txscript.RawTxInWitnessSignature(mtx, sh, i, int64(selected[i].Amount), rdm,
				txscript.SigHashAll, p.PrivKey)
			if err != nil {
				panic(fmt.Errorf("failed to sign no %d: %v", i, err))
			}
			data[j+1] = sig
		}
		data[6] = rdm
		in.Witness = wire.TxWitness(data)
	}

	var buf bytes.Buffer
	_ = mtx.BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
	fmt.Printf("tx: %s\n", hex.EncodeToString(buf.Bytes()))
}
