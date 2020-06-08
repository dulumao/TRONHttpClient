package tronhttpClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	httpClient "github.com/stdevHsequeda/TRONHttpClient/client"
	"io/ioutil"
	"net/http"
)

const testNet = "https://api.shasta.trongrid.io"
const mainNet = "https://api.trongrid.io"

type Client struct {
	client  *httpClient.Client
	network string
}

func NewClient(network string) *Client {
	httpClient.MaxRetry = 5
	return &Client{client: httpClient.NewClient(), network: network}
}

func (c *Client) createTrx(toAddr, ownerAddr string, amount int) (*Transaction, error) {
	encodeData, err := json.Marshal(
		map[string]interface{}{
			"to_address":    toAddr,
			"owner_address": ownerAddr,
			"amount":        amount,
		})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", testNet+"/wallet/createtransaction",
		bytes.NewBuffer(encodeData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

	var tx Transaction
	err = json.NewDecoder(resp).Decode(&tx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v \n", tx)
	return &tx, err
}

func (c *Client) getTxSign(tx *Transaction, privKey string) (*Transaction, error) {
	encodeData, err := json.Marshal(
		struct {
			Transaction *Transaction `json:"transaction"`
			PrivateKey  string       `json:"privateKey"`
		}{
			Transaction: tx,
			PrivateKey:  privKey,
		})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", testNet+"/wallet/gettransactionsign",
		bytes.NewBuffer(encodeData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp).Decode(tx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v \n", tx)
	return tx, err
}

func (c *Client) broadcastTx(tx *Transaction) error {
	encodeData, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", testNet+"/wallet/broadcasttransaction",
		bytes.NewBuffer(encodeData))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.CallRetryable(req)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp)
	if err != nil {
		return err
	}

	fmt.Printf("%s \n", string(b))

	return nil
}

func (c *Client) generateAddress() (*Address, error) {
	req, err := http.NewRequest("GET", testNet+"/wallet/generateaddress",
		nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

	var addr Address
	err = json.NewDecoder(resp).Decode(&addr)
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

func (c *Client) createAddress() (*Address, error) {
	req, err := http.NewRequest("GET", testNet+"/wallet/generateaddress",
		nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

	var addr Address
	err = json.NewDecoder(resp).Decode(&addr)
	if err != nil {
		return nil, err
	}

	return &addr, nil
}
