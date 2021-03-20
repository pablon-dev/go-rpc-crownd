package crownd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	*http.Client
	reqAddr string
	rpcUser string
	rpcPass string
}

type request struct {
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int64         `json:"id"`
	JsonRpc string        `json:"jsonrpc"`
}

type response struct {
	Id     int64           `json:"id"`
	Result json.RawMessage `json:"result"`
	Err    *responseError  `json:"error"`
}
type responseError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewClientWithSSL(host string, port int, user, pass string, timeout int) (*Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cssl := &http.Client{
		Timeout:   time.Duration(timeout) * time.Millisecond,
		Transport: transport,
	}
	return newClient(cssl, "https//:", host, port, user, pass)

}

func NewClient(host string, port int, user, pass string, timeout int) (*Client, error) {
	c := &http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	return newClient(c, "http//:", host, port, user, pass)
}

func newClient(client *http.Client, prefix, host string, port int, user, pass string) (*Client, error) {
	if len(host) == 0 {
		return nil, errors.New("missing host")
	}
	addr := fmt.Sprintf("http://%s:%d", host, port)
	return &Client{
		Client:  client,
		reqAddr: addr,
		rpcUser: user,
		rpcPass: pass,
	}, nil

}

func (client *Client) Request(method string, params ...interface{}) (*response, error) {
	req := &request{
		Method:  method,
		Params:  params,
		ID:      time.Now().UnixNano(),
		JsonRpc: "1.0",
	}
	return client.doRequest(req)

}
func (client *Client) doRequest(reqRPC *request) (*response, error) {
	reqBodyMar, err := json.Marshal(reqRPC)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", client.reqAddr, bytes.NewBuffer(reqBodyMar))
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")
	if len(client.rpcUser) > 0 {
		req.SetBasicAuth(client.rpcUser, client.rpcPass)
	}
	httpResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	resp := &response{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func parseErr(reqerr error, resperr *responseError) (err error) {
	if reqerr != nil {
		return err
	}
	if resperr != nil{
		err = errors.New(fmt.Sprintf("Error code: %d\nError message: %s", resperr.Code, resperr.Message))
	}
	return err
}
