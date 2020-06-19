module github.com/zouxyan/btctool

go 1.14

require (
	github.com/ConjurTech/switcheo-chain v0.0.0-00010101000000-000000000000
	github.com/FactomProject/basen v0.0.0-20150613233007-fe3947df716e // indirect
	github.com/andlabs/ui v0.0.0-20180902183112-867a9e5a498d
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/btcsuite/golangcrypto v0.0.0-20150304025918-53f62d9b43e8
	github.com/cosmos/cosmos-sdk v0.38.4
	github.com/ontio/btc-spvclient v0.0.0-00010101000000-000000000000
	github.com/ontio/eth_tools v0.0.0-00010101000000-000000000000
	github.com/ontio/go-bip32 v0.0.0-20190520025953-d3cea6894a2b // indirect
	github.com/ontio/multi-chain v0.0.0
	github.com/ontio/multi-chain-go-sdk v0.0.0-00010101000000-000000000000
	github.com/ontio/ontology v1.10.0
	github.com/ontio/ontology-crypto v1.0.9
	github.com/skyinglyh1/cosmos-poly-module v0.0.0-20200611071020-364a883092be
)

replace (
	github.com/ConjurTech/switcheo-chain => ../../ConjurTech/switcheo-chain
	github.com/cosmos/gaia => github.com/skyinglyh1/gaia v1.0.1-0.20200608095633-f3d8a0b00305
	github.com/ontio/btc-spvclient => ../../ontio/btc-spvclient
	github.com/ontio/eth_tools => ../../ontio/eth_tools
	github.com/ontio/multi-chain => ../../ontio/multi-chain
	github.com/ontio/multi-chain-go-sdk => ../../ontio/multi-chain-go-sdk
)