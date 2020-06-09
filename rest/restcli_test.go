package rest

import (
	"fmt"
	"testing"
)

func TestRestCli_GenerateToAddr(t *testing.T) {
	cli := NewRestCli("http://172.168.3.77:18443", "test", "test")
	addrs, err := cli.GenerateToAddr(1, "mjEoyyCPsLzJ23xMX6Mti13zMyN36kzn57")
	if err != nil {
		t.Fatalf("failed to generate blk: %v", err)
	}
	if len(addrs) != 1 {
		t.Fatalf("wrong length of addrs: should be 1 not %d", len(addrs))
	}
	h, err := cli.GetBlockHeight(addrs[0])
	if err != nil {
		t.Fatalf("failed to get block height: %v", err)
	}
	fmt.Println("height is", h)
}

func TestRestCli_GetMempoolInfo(t *testing.T) {
	cli := NewRestCli("http://172.168.3.77:18443", "test", "test")
	size, err := cli.GetMempoolInfo()
	if err != nil {
		t.Fatalf("failed to get mempool info: %v", err)
	}

	fmt.Println("size:", size)
}

func TestRestCli_ListUnspent(t *testing.T) {
	cli := NewRestCli("http://172.168.3.77:18443", "test", "test")

	utxos, err := cli.ListUnspent(1, 999, "2N5cY8y9RtbbvQRWkX5zAwTPCxSZF9xEj2C")
	if err != nil {
		t.Fatalf("failed to list: %v", err)
	}
	if len(utxos) == 0 {
		t.Fatalf("shouldn't be zero")
	}
	for _, v := range utxos {
		fmt.Println(v.ScriptPubKey, v.Txid, v.Amount)
	}
}

func TestRestCli_ImportAddress(t *testing.T) {
	cli := NewRestCli("http://172.168.3.77:18443", "test", "test")

	err := cli.ImportAddress("2N5cY8y9RtbbvQRWkX5zAwTPCxSZF9xEj2C")
	if err != nil {
		t.Fatalf("failed to import address: %v", err)
	}
}

func TestRestCli_GetBlockCount(t *testing.T) {
	cli := NewRestCli("http://172.168.3.77:18443", "test", "test")

	cnt, err := cli.GetBlockCount()
	if err != nil {
		t.Fatalf("failed to get count: %v", err)
	}
	fmt.Println(cnt)
}

func TestRestCli_BroadcastTxBySpv(t *testing.T) {
	cli := NewRestCli("http://172.168.3.77:18443", "test", "test")
	err := cli.BroadcastTxBySpv("0100000001ee1e6d9a24cd7d66fa3641e7d44fe7cfcbe3026d523c9b2327dab6056cee838a020000006b483045022100be224e258a1900488df22d2016f03e178bf4286d00723eafcf8dd4a8ff7c9d2202207765c8a26e55982ee17152a1eedcc8b6eabc70a474315c8444ec97ed5044fa71012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03102700000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d3870000000000000000276a256600000000000000020000000000000000dab47e816313a79c9459b544720c90a725264e0da0bb0d00000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000")
	if err != nil {
		t.Fatalf("failed to broadcast: %v", err)
	}
}
