package builder

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type BuildCrossChainTxParam struct {
	Inputs       []btcjson.TransactionInput
	Changes      map[string]float64 //pay to pubK
	ToMultiValue float64
	Redeem       string
	Locktime     *int64
	PrevPkScript []byte
	Privk        *btcec.PrivateKey
	NetParam     *chaincfg.Params
	Data         []byte
	IsSegWit     int
}

type Builder struct {
	NetParam     *chaincfg.Params
	PrevPkScript []byte
	PrivKey      *btcec.PrivateKey
	PubKey       *btcec.PublicKey
	Tx           *wire.MsgTx
	IsSigned     bool
	Privks       map[string]*btcec.PrivateKey
	IsSegWit     int
}

func NewBuilder(param *BuildCrossChainTxParam) (b *Builder, err error) {
	b = &Builder{}
	b.IsSegWit = param.IsSegWit
	b.NetParam = param.NetParam
	b.PrivKey = param.Privk
	b.PrevPkScript = param.PrevPkScript

	mtx, err := b.getUnsignedCrossChainTx(param.Inputs, param.Changes, param.Redeem, param.ToMultiValue,
		param.Locktime, param.NetParam, param.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to get raw tx: %v", err)
	}
	b.Tx = mtx

	return b, nil
}

func (builder *Builder) LookUpKey(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
	if builder.PrivKey == nil {
		return nil, false, errors.New("Private key not ready")
	}
	return builder.PrivKey, true, nil
}

// locking
func (builder *Builder) BuildSignedTx() error {
	for i, _ := range builder.Tx.TxIn {
		sig, err := txscript.SignTxOutput(builder.NetParam, builder.Tx, i, builder.PrevPkScript,
			txscript.SigHashAll, txscript.KeyClosure(builder.LookUpKey), nil, nil)
		if err != nil {
			return fmt.Errorf("Failed to sign tx's No.%d input: %v", i, err)
		}
		if err != nil {
			return fmt.Errorf("Failed to get witness: %v", err)
		}
		builder.Tx.TxIn[i].SignatureScript = sig
	}
	builder.IsSigned = true
	return nil
}

// need to make a multisig-output tx
func (builder *Builder) getUnsignedCrossChainTx(txIns []btcjson.TransactionInput, changes map[string]float64, redeem string,
	value float64, locktime *int64, netParam *chaincfg.Params, data []byte) (*wire.MsgTx, error) {
	if locktime != nil && (*locktime < 0 || *locktime > int64(wire.MaxTxInSequenceNum)) {
		return nil, fmt.Errorf("getRawTxToMultiAddr, locktime %d out of range", *locktime)
	}
	if value < 0 || value > btcutil.MaxSatoshi {
		return nil, fmt.Errorf("getRawTxToMultiAddr, wrong value to multi-addr: %f", value)
	}

	// Add all transaction inputs to a new transaction after performing
	// some validity checks.
	mtx := wire.NewMsgTx(wire.TxVersion)
	for _, input := range txIns {
		txHash, err := chainhash.NewHashFromStr(input.Txid)
		if err != nil {
			return nil, fmt.Errorf("getRawTxToMultiAddr, decode txid fail: %v", err)
		}

		prevOut := wire.NewOutPoint(txHash, input.Vout)
		txIn := wire.NewTxIn(prevOut, []byte{}, nil)
		if locktime != nil && *locktime != 0 {
			txIn.Sequence = wire.MaxTxInSequenceNum - 1
		}
		mtx.AddTxIn(txIn)
	}

	valueInSatoshi, err := btcutil.NewAmount(value)
	if err != nil {
		return nil, fmt.Errorf("getRawTxToMultiAddr, failed to convert value: %v", err)
	}

	var addr btcutil.Address
	rb, err := hex.DecodeString(redeem)
	if err != nil {
		return nil, fmt.Errorf("getRawTxToMultiAddr, failed to decode redeem hex: %v", err)
	}
	fmt.Println("wit", builder.IsSegWit)
	if builder.IsSegWit == 1 {
		hasher := sha256.New()
		hasher.Write(rb)

		addr, err = btcutil.NewAddressWitnessScriptHash(hasher.Sum(nil), builder.NetParam)
		if err != nil {
			return nil, fmt.Errorf("failed to new witness address: %v", err)
		}
	} else {
		addr, err = btcutil.NewAddressScriptHash(rb, builder.NetParam)
		if err != nil {
			return nil, fmt.Errorf("failed to new p2sh addr: %v", err)
		}
	}

	s, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, fmt.Errorf("getRawTxToMultiAddr, failed to get p2sh script: %v", err)
	}

	mtx.AddTxOut(wire.NewTxOut(int64(valueInSatoshi), s))

	// add NullData
	nullDataScript, err := txscript.NullDataScript(data)
	if err != nil {
		return nil, fmt.Errorf("getRawTxToMultiAddr, failed to build nulldata script")
	}
	mtx.AddTxOut(wire.NewTxOut(0, nullDataScript))

	// Add all transaction outputs to the transaction after performing
	// some validity checks.
	for encodedAddr, amount := range changes {
		// Ensure amount is in the valid range for monetary amounts.
		if amount <= 0 || amount > btcutil.MaxSatoshi {
			return nil, fmt.Errorf("getRawTxToMultiAddr, wrong amount: %f", amount)
		}

		// Decode the provided address.
		addr, err := btcutil.DecodeAddress(encodedAddr, netParam)
		if err != nil {
			return nil, fmt.Errorf("getRawTxToMultiAddr, decode addr fail: %v", err)
		}

		// Ensure the address is one of the supported types and that
		// the network encoded with the address matches the network the
		// server is currently on.
		switch addr.(type) {
		case *btcutil.AddressPubKeyHash:
		case *btcutil.AddressScriptHash:
		default:
			return nil, fmt.Errorf("getRawTxToMultiAddr, type of addr is not found")
		}
		if !addr.IsForNet(netParam) {
			return nil, fmt.Errorf("getRawTxToMultiAddr, addr is not for mainnet")
		}

		// Create a new script which pays to the provided address.
		pkScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return nil, fmt.Errorf("getRawTxToMultiAddr, failed to generate pay-to-address script: %v", err)
		}

		// Convert the amount to satoshi.
		satoshi, err := btcutil.NewAmount(amount)
		if err != nil {
			return nil, fmt.Errorf("getRawTxToMultiAddr, failed to convert amount: %v", err)
		}

		txOut := wire.NewTxOut(int64(satoshi), pkScript)
		mtx.AddTxOut(txOut)
	}

	// Set the Locktime, if given.
	if locktime != nil {
		mtx.LockTime = uint32(*locktime)
	}

	return mtx, nil
}
