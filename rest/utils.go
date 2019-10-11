package rest

import (
	"fmt"
	"sort"
	"strconv"
)

type Utxo struct {
	Txid         string
	Vout         uint32
	ScriptPubKey string
	Amount       int64
	Confs        int64
}

func (u *Utxo) String() string {
	return "\nutxo: " + u.Txid + ":" + strconv.FormatUint(uint64(u.Vout), 10) + "\n\tScriptPubKey: " +
		u.ScriptPubKey + "\n\tAmount: " + strconv.FormatInt(u.Amount, 10) + "\n\tConfs: " +
		strconv.FormatInt(u.Confs, 10) + "\n"
}

func SelectUtxos(utxos []*Utxo, value int64) ([]*Utxo, int64, error) {
	if value <= 0 {
		return nil, -1, fmt.Errorf("value must be positive")
	}

	ul := make(UtxoList, 0)
	for _, v := range utxos {
		ul = append(ul, &UtxoItem{
			key: v,
			val: v.Confs * v.Amount,
		})
	}

	sort.Sort(ul)
	if !sort.IsSorted(ul) {
		return nil, -1, fmt.Errorf("ItemList not sorted")
	}

	selected := make([]*Utxo, 0)
	selectedVal := int64(0)
	for i := len(ul) - 1; i >= 0; i-- {
		selected = append(selected, ul[i].key)
		if selectedVal += ul[i].key.Amount; selectedVal >= value {
			return selected, selectedVal, nil
		}
	}
	return nil, selectedVal, fmt.Errorf("not enough utxo for %d, all we have is %d", value, selectedVal)
}

type UtxoItem struct {
	key *Utxo
	val int64
}

type UtxoList []*UtxoItem

func (ul UtxoList) Swap(i, j int) {
	ul[i], ul[j] = ul[j], ul[i]
}

func (ul UtxoList) Len() int {
	return len(ul)
}

func (ul UtxoList) Less(i, j int) bool {
	return ul[i].val < ul[j].val
}
