module github.com/ontio/multi-chain-go-sdk

go 1.14

require (
	github.com/FactomProject/basen v0.0.0-20150613233007-fe3947df716e // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/itchyny/base58-go v0.1.0
	github.com/ontio/go-bip32 v0.0.0-20190520025953-d3cea6894a2b
	github.com/ontio/multi-chain v0.0.0
	github.com/ontio/ontology-crypto v1.0.9-0.20200427073653-7eba89bbe29b
	github.com/stretchr/testify v1.4.0
	github.com/tyler-smith/go-bip39 v1.0.2
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79
)

replace github.com/ontio/multi-chain v0.0.0 => ./../multi-chain
