package service

import (
	"encoding/hex"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/core/types"
	"log"
	"os"
)

const (
	BTCPRIVK_PATH = "./btcprivk"
)

func EncryptBtcPrivk(privk, pwd string) {
	pri, err := keypair.GetP256KeyPairFromWIF([]byte(privk))
	privk = ""
	if err != nil {
		log.Fatalf("failed to get privk from wif: %v", err)
		os.Exit(1)
	}

	wallet, err := account.Open(BTCPRIVK_PATH)
	if err != nil {
		log.Fatalf("failed to open wallet: %v", err)
		os.Exit(1)
	}

	pub := pri.Public()
	addr := types.AddressFromPubKey(pub)
	b58addr := addr.ToBase58()
	k, err := keypair.EncryptPrivateKey(pri, b58addr, []byte(pwd))
	if err != nil {
		log.Fatalf("failed to encrypt private key: %v", err)
		os.Exit(1)
	}

	var accMeta account.AccountMetadata
	accMeta.Address = k.Address
	accMeta.KeyType = k.Alg
	accMeta.EncAlg = k.EncAlg
	accMeta.Hash = k.Hash
	accMeta.Key = k.Key
	accMeta.Curve = k.Param["curve"]
	accMeta.Salt = k.Salt
	accMeta.Label = ""
	accMeta.PubKey = hex.EncodeToString(keypair.SerializePublicKey(pub))
	accMeta.SigSch = signature.SHA256withECDSA.Name()

	err = wallet.ImportAccount(&accMeta)
	if err != nil {
		log.Fatalf("failed to import account: %v", err)
		os.Exit(1)
	}
	pwd = ""
}
