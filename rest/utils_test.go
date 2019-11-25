package rest

import (
	"fmt"
	"github.com/btcsuite/btcutil"
	"math"
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
	aa := make([]int, 0)
	fmt.Printf("%v, cap(%d)\n", aa, cap(aa))

	aa = append(aa, 1, 2, 3)
	fmt.Printf("%v, cap(%d)\n", aa, cap(aa))

	bb := make([]int, 10241)
	aa = append(aa, bb...)
	fmt.Printf("%v, cap(%d)\n", nil, cap(aa))
	fmt.Println(math.Log2(float64(11264)))
}
