package rest

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/golangcrypto/ripemd160"
	"testing"
)

func TestSelectUtxos(t *testing.T) {
	val := 101
	cli := NewRestCli("http://172.168.3.77:18443", "test", "test", "138.91.6.125:50071")
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
	bb, _ := hex.DecodeString("44978a77e4e983136bf1cca277c45e5bd4eff6a7848e900416daf86fd32c2743")
	r := ripemd160.New()
	r.Write(bb)
	fmt.Println(hex.EncodeToString(r.Sum(nil)))
}
