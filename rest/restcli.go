package rest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Request struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}

type Response struct {
	Result interface{}       `json:"result"`
	Error  *btcjson.RPCError `json:"error"` //maybe wrong
	Id     int               `json:"id"`
}

// Get tx in block; Get proof;
type RestCli struct {
	Addr    string
	Cli     *http.Client
}

func NewRestCli(addr, user, pwd string) *RestCli {
	return &RestCli{
		Cli: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   5,
				DisableKeepAlives:     false,
				IdleConnTimeout:       time.Second * 300,
				ResponseHeaderTimeout: time.Second * 300,
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				Proxy: func(req *http.Request) (*url.URL, error) {
					req.SetBasicAuth(user, pwd)
					return nil, nil
				},
			},
			Timeout: time.Second * 300,
		},
		Addr:    addr,
	}
}

func (cli *RestCli) sendPostReq(req []byte) (*Response, error) {
	resp, err := cli.sendPostReqWithAddr(cli.Addr, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (cli *RestCli) sendPostReqWithAddr(addr string, req []byte) (*Response, error) {
	resp, err := cli.Cli.Post(addr, "application/json;charset=UTF-8",
		bytes.NewReader(req))
	if err != nil {
		return nil, fmt.Errorf("failed to post: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error:%s", err)
	}

	response := new(Response)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return response, nil
}

func (cli *RestCli) SendRestRequestToSpv(addr string, data []byte) ([]byte, error) {
	resp, err := cli.Cli.Post(addr, "application/json;charset=UTF-8",
		bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("rest post request:%s error:%s", data, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read rest response body error:%s", err)
	}
	return body, nil
}

func (cli *RestCli) GenerateToAddr(n int, addr string) ([]string, error) {
	req, err := json.Marshal(Request{
		Jsonrpc: "1.0",
		Method:  "generatetoaddress",
		Params:  []interface{}{n, addr},
		Id:      1,
	})
	if err != nil {
		return nil, fmt.Errorf("[GenerateToAddr] failed to marshal request: %v", err)
	}

	resp, err := cli.sendPostReq(req)
	if err != nil {
		return nil, fmt.Errorf("[GenerateToAddr] failed to send post: %v", err)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("[GenerateToAddr] response shows failure: %v", resp.Error.Message)
	}

	hashes := make([]string, 0)
	for _, v := range resp.Result.([]interface{}) {
		hashes = append(hashes, v.(string))
	}

	return hashes, nil
}

func (cli *RestCli) GetBlockHeight(hash string) (int32, error) {
	req, err := json.Marshal(Request{
		Jsonrpc: "1.0",
		Method:  "getblock",
		Params:  []interface{}{hash},
		Id:      1,
	})
	if err != nil {
		return -1, fmt.Errorf("[GetBlockHeight] failed to marshal request: %v", err)
	}

	resp, err := cli.sendPostReq(req)
	if err != nil {
		return -1, fmt.Errorf("[GetBlockHeight] failed to send post: %v", err)
	}
	if resp.Error != nil {
		return -1, fmt.Errorf("[GetBlockHeight] response shows failure: %v", resp.Error.Message)
	}

	return int32(resp.Result.(map[string]interface{})["height"].(float64)), nil
}

func (cli *RestCli) GetMempoolInfo() (int32, error) {
	req, err := json.Marshal(Request{
		Jsonrpc: "1.0",
		Method:  "getmempoolinfo",
		Params:  []interface{}{},
		Id:      1,
	})
	if err != nil {
		return -1, fmt.Errorf("[GetMempoolInfo] failed to marshal request: %v", err)
	}

	resp, err := cli.sendPostReq(req)
	if err != nil {
		return -1, fmt.Errorf("[GetMempoolInfo] failed to send post: %v", err)
	}
	if resp.Error != nil {
		return -1, fmt.Errorf("[GetMempoolInfo] response shows failure: %v", resp.Error.Message)
	}

	return int32(resp.Result.(map[string]interface{})["size"].(float64)), nil
}

func (cli *RestCli) ListUnspent(minConfs, maxConfs int64, addr string) ([]*Utxo, error) {
	req, err := json.Marshal(Request{
		Jsonrpc: "1.0",
		Method:  "listunspent",
		Params:  []interface{}{minConfs, maxConfs, []string{addr}},
		Id:      1,
	})
	if err != nil {
		return nil, fmt.Errorf("[ListUnspent] failed to marshal request: %v", err)
	}

	resp, err := cli.sendPostReq(req)
	if err != nil {
		return nil, fmt.Errorf("[ListUnspent] failed to send post: %v", err)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("[ListUnspent] response shows failure: %v", resp.Error.Message)
	}

	utxos := make([]*Utxo, 0)
	arr, ok := resp.Result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	for _, v := range arr {
		item := v.(map[string]interface{})

		amount := item["amount"].(float64)
		txid := item["txid"].(string)
		a, err := btcutil.NewAmount(amount)
		if err != nil {
			return nil, fmt.Errorf("failed to get amount for %f(txid:%s)", amount, txid)
		}
		utxos = append(utxos, &Utxo{
			Txid:         txid,
			Vout:         uint32(item["vout"].(float64)),
			ScriptPubKey: item["scriptPubKey"].(string),
			Amount:       int64(a),
			Confs:        int64(item["confirmations"].(float64)),
		})
	}

	return utxos, err
}

func (cli *RestCli) ImportAddress(addr string) error {
	info, err := cli.GetAddressInfo(addr)
	if err != nil {
		return err
	}
	if len(info["labels"].([]interface{})) != 0 {
		return nil
	}
	req, err := json.Marshal(Request{
		Jsonrpc: "1.0",
		Method:  "importaddress",
		Params:  []interface{}{addr},
		Id:      1,
	})
	if err != nil {
		return fmt.Errorf("[ImportAddress] failed to marshal request: %v", err)
	}

	resp, err := cli.sendPostReq(req)
	if err != nil {
		return fmt.Errorf("[ImportAddress] failed to send post: %v", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("[ImportAddress] response shows failure: %v", resp.Error.Message)
	}

	return nil
}

func (cli *RestCli) GetAddressInfo(addr string) (map[string]interface{}, error) {
	req, err := json.Marshal(Request{
		Jsonrpc: "1.0",
		Method:  "getaddressinfo",
		Params:  []interface{}{addr},
		Id:      1,
	})
	if err != nil {
		return nil, fmt.Errorf("[GetAddressInfo] failed to marshal request: %v", err)
	}

	resp, err := cli.sendPostReq(req)
	if err != nil {
		return nil, fmt.Errorf("[GetAddressInfo] failed to send post: %v", err)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("[GetAddressInfo] response shows failure: %v", resp.Error.Message)
	}

	return resp.Result.(map[string]interface{}), nil
}

func (cli *RestCli) GetBlockCount() (int64, error) {
	req, err := json.Marshal(Request{
		Jsonrpc: "1.0",
		Method:  "getblockcount",
		Params:  []interface{}{},
		Id:      1,
	})
	if err != nil {
		return -1, fmt.Errorf("[GetBlockCount] failed to marshal request: %v", err)
	}

	resp, err := cli.sendPostReq(req)
	if err != nil {
		return -1, fmt.Errorf("[GetBlockCount] failed to send post: %v", err)
	}
	if resp.Error != nil {
		return -1, fmt.Errorf("[GetBlockCount] response shows failure: %v", resp.Error.Message)
	}

	return int64(resp.Result.(float64)), nil
}

func (cli *RestCli) SendRawTx(rawTx string) (string, error) {
	req, err := json.Marshal(Request{
		Jsonrpc: "1.0",
		Method:  "sendrawtransaction",
		Params:  []interface{}{rawTx},
		Id:      1,
	})
	if err != nil {
		return "", fmt.Errorf("[GetMempoolInfo] failed to marshal request: %v", err)
	}

	resp, err := cli.sendPostReq(req)
	if err != nil {
		return "", fmt.Errorf("[GetMempoolInfo] failed to send post: %v", err)
	}
	if resp.Error != nil {
		return "", fmt.Errorf("[GetMempoolInfo] response shows failure: %v", resp.Error.Message)
	}

	return resp.Result.(string), nil
}

func (cli *RestCli) BroadcastTxBySpv(tx string) error {
	//if cli.SpvAddr == "" {
	//	return fmt.Errorf("spv addr not set")
	//}
	//req, err := json.Marshal(common.BroadcastReq{
	//	Tx: tx,
	//})
	//if err != nil {
	//	return fmt.Errorf("[BroadcastTxBySpv] failed to marshal request: %v", err)
	//}
	//
	//data, err := cli.SendRestRequestToSpv("http://"+cli.SpvAddr+"/api/v1/broadcasttx", req)
	//if err != nil {
	//	return fmt.Errorf("Failed to send request: %v", err)
	//}
	//
	//var resp common.Response
	//err = json.Unmarshal(data, &resp)
	//if err != nil {
	//	return fmt.Errorf("Failed to unmarshal resp to json: %v", err)
	//}
	//if resp.Error != 0 || resp.Desc != "SUCCESS" {
	//	return fmt.Errorf("Response shows failure: %s", resp.Desc)
	//}

	return nil
}
