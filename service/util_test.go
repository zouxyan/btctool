package service

import (
	"fmt"
	"github.com/zouxyan/btctool/rest"
	"testing"
)

func TestBuildData(t *testing.T) {
	addr := "0x7dD16c0c71F71A123c4BDAF0a468aBC60Db41C0C"
	_, err := buildData(1, 0, addr, addr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBlkGene_Run(t *testing.T) {
	cli := rest.NewRestCli("http://172.168.3.77:18443", "test", "test", "")4

	for cnt := 10; cnt > 0; cnt-- {
		hs, err := cli.GenerateToAddr(1, "mjEoyyCPsLzJ23xMX6Mti13zMyN36kzn57")
		if err != nil {
			t.Fatalf("failed to generate block: %v", err)
			continue
		}
		h, err := cli.GetBlockHeight(hs[0])
		if err != nil {
			t.Fatalf("failed to get block height: %v", err)
			continue
		}
		fmt.Printf("generate block %s(height:%d) to address %s", hs[0], h, "mjEoyyCPsLzJ23xMX6Mti13zMyN36kzn57")
	}
}
