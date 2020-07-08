package service

import (
	"github.com/polynetwork/poly/common/log"
	"github.com/zouxyan/btctool/rest"
	"os"
	"time"
)

type BlkGene struct {
	RpcUrl      string
	User        string
	Pwd         string
	Tsec        uint
	DefaultAddr string
}

func (bg *BlkGene) Run() {
	if bg.Tsec <= 0 {
		log.Error("tsec must be positive")
		os.Exit(1)
	}
	if bg.RpcUrl == "" || bg.User == "" || bg.Pwd == "" || bg.DefaultAddr == "" {
		log.Error("wrong parameters")
		os.Exit(1)
	}

	cli := rest.NewRestCli(bg.RpcUrl, bg.User, bg.Pwd)

	tick := time.NewTicker(time.Duration(bg.Tsec) * time.Second)
	for {
		select {
		case <-tick.C:
			if size, err := cli.GetMempoolInfo(); size <= 0 || err != nil {
				if err != nil {
					log.Errorf("failed to get mempool info %v", err)
				}
				continue
			}
			hs, err := cli.GenerateToAddr(1, bg.DefaultAddr)
			if err != nil {
				log.Errorf("failed to generate block: %v", err)
				continue
			}
			h, err := cli.GetBlockHeight(hs[0])
			if err != nil {
				log.Errorf("failed to get block height: %v", err)
				continue
			}
			log.Infof("generate block %s(height:%d) to address %s", hs[0], h, bg.DefaultAddr)
		}
	}
}
