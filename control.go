package crownd

import "encoding/json"

type GetInfoResponse struct {
	Version         float64
	ProtocolVersion float64
	WalletVersion   float64
	Balance         float64
	Blocks          int
	TimeOffset      int
	Proxy           string
	Difficulty      float64
	Tesnet          bool
	StakingActive   bool
	KeyPoolOldest   float64
	KeyPoolSize     float64
	UnlockedUntil   float64
	PayTxFee        float64
	RelayFee        float64
	Errors          json.RawMessage
}

func (client *Client) GetInfo() (*GetInfoResponse, error) {
	resp, err := client.Request("getinfo") 
	if resperr := parseErr(err,resp.Err); resperr != nil {
		return nil, resperr
	}
	getinforesp := &GetInfoResponse{}
	err = json.Unmarshal(resp.Result,getinforesp)
	if err != nil {
		return nil, err
	}
	return getinforesp,nil

	
}
