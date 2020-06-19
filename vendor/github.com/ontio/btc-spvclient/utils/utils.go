package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	sdk "github.com/ontio/multi-chain-go-sdk"
	"golang.org/x/crypto/ripemd160"
	"time"
)

const (
	MIN_FEE        = 100
	SIGNED_TX_KEY  = "btcTxToRelay"
	TO_SIGN_TX_KEY = "makeBtcTx"
)

type ToSignItem struct {
	Mtx  *wire.MsgTx
	Amts []uint64
}

func (item *ToSignItem) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	var txb bytes.Buffer
	if err := item.Mtx.BtcEncode(&txb, wire.ProtocolVersion, wire.LatestEncoding); err != nil {
		return nil, err
	}
	tx := txb.Bytes()
	if err := binary.Write(&buf, binary.BigEndian, uint32(len(tx))); err != nil {
		return nil, err
	}
	buf.Write(tx)

	if err := binary.Write(&buf, binary.BigEndian, uint32(len(item.Amts))); err != nil {
		return nil, err
	}
	for _, v := range item.Amts {
		if err := binary.Write(&buf, binary.BigEndian, v); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (item *ToSignItem) Deserialize(buf []byte) error {
	r := bytes.NewReader(buf)
	var lenTx uint32
	if err := binary.Read(r, binary.BigEndian, &lenTx); err != nil {
		return err
	}
	rawTx := make([]byte, lenTx)
	if _, err := r.Read(rawTx); err != nil {
		return err
	}
	tx := wire.NewMsgTx(wire.TxVersion)
	if err := tx.BtcDecode(bytes.NewBuffer(rawTx), wire.ProtocolVersion, wire.LatestEncoding); err != nil {
		return err
	}
	item.Mtx = tx

	var lenAmts uint32
	if err := binary.Read(r, binary.BigEndian, &lenAmts); err != nil {
		return err
	}
	amts := make([]uint64, lenAmts)
	var val uint64
	for i := uint32(0); i < lenAmts; i++ {
		if err := binary.Read(r, binary.BigEndian, &val); err != nil {
			return err
		}
		amts[i] = val
	}
	item.Amts = amts

	return nil
}

type SavedItem struct {
	Item         *ToSignItem
	TimeReceived time.Time
	Done         bool
}

func (saved *SavedItem) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	raw, err := saved.Item.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize ToSignItem: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, uint32(len(raw))); err != nil {
		return nil, err
	}
	buf.Write(raw)

	t, err := saved.TimeReceived.GobEncode()
	if err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, uint32(len(t))); err != nil {
		return nil, err
	}
	buf.Write(t)

	if err = binary.Write(&buf, binary.BigEndian, saved.Done); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (saved *SavedItem) Deserialize(buf []byte) error {
	r := bytes.NewReader(buf)
	var lenItem uint32
	if err := binary.Read(r, binary.BigEndian, &lenItem); err != nil {
		return err
	}
	raw := make([]byte, lenItem)
	if _, err := r.Read(raw); err != nil {
		return err
	}
	item := &ToSignItem{}
	if err := item.Deserialize(raw); err != nil {
		return err
	}
	saved.Item = item

	var tr time.Time
	var lenTr uint32
	if err := binary.Read(r, binary.BigEndian, &lenTr); err != nil {
		return err
	}
	raw = make([]byte, lenTr)
	if _, err := r.Read(raw); err != nil {
		return err
	}
	if err := tr.GobDecode(raw); err != nil {
		return err
	}
	saved.TimeReceived = tr

	var done bool
	if err := binary.Read(r, binary.BigEndian, &done); err != nil {
		return err
	}
	saved.Done = done

	return nil
}

type SavedItemArr []*SavedItem

func (arr SavedItemArr) Less(i, j int) bool {
	if arr[i].TimeReceived.Before(arr[j].TimeReceived) {
		return true
	}
	return false
}

func (arr SavedItemArr) Swap(i, j int) {
	temp := arr[i]
	arr[i] = arr[j]
	arr[j] = temp
}

func (arr SavedItemArr) Len() int {
	return len(arr)
}

func GetAccountByPassword(sdk *sdk.MultiChainSdk, path string, pwd []byte) (*sdk.Account, error) {
	wallet, err := sdk.OpenWallet(path)
	if err != nil {
		return nil, fmt.Errorf("open wallet %s error: %v", path, err)
	}
	user, err := wallet.GetDefaultAccount(pwd)
	if err != nil {
		return nil, fmt.Errorf("getDefaultAccount error: %v", err)
	}
	return user, nil
}

func Wait(dura time.Duration) {
	t := time.NewTimer(dura)
	<-t.C
	t.Stop()
}

func GetUtxoKey(scriptPk []byte) string {
	switch txscript.GetScriptClass(scriptPk) {
	case txscript.MultiSigTy:
		return hex.EncodeToString(btcutil.Hash160(scriptPk))
	case txscript.ScriptHashTy:
		return hex.EncodeToString(scriptPk[2:22])
	case txscript.WitnessV0ScriptHashTy:
		hasher := ripemd160.New()
		hasher.Write(scriptPk[2:34])
		return hex.EncodeToString(hasher.Sum(nil))
	default:
		return ""
	}
}
