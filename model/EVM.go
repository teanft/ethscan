package model

import "math/big"

type RPCResult struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

type RPCBody struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      int      `json:"id"`
}

type Block struct {
	Address     string   `json:"address"`
	Contract    string   `json:"contract"`
	From        string   `json:"from"`
	Gas         uint64   `json:"gas"`
	Private     string   `json:"private"`
	To          string   `json:"to"`
	Value       *big.Int `json:"value"`
	SignedTx    string   `json:"signed_tx"`
	RawTx       string   `json:"raw_tx"`
	BlockNumber *big.Int `json:"block_number"`
	BlockHash   string   `json:"block_hash"`
	TxHash      string   `json:"tx_hash"`
}
