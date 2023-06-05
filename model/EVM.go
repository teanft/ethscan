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

type EVM struct {
	Address     string   `json:"address"`
	Contract    string   `json:"contract"`
	From        string   `json:"from"`
	Gas         uint64   `json:"gas"`
	Private     string   `json:"private"`
	To          string   `json:"to"`
	Value       *big.Int `json:"value"`
	SignedTx    string   `json:"signedTx"`
	RawTx       string   `json:"rawTx"`
	BlockNumber *big.Int `json:"blockNumber"`
	BlockHash   string   `json:"blockHash"`
	TxHash      string   `json:"tx_hash"`
}

type Transaction struct {
	TransactionIndex     uint     `json:"transactionIndex"`
	TxHash               string   `json:"hash"`
	BlockNumber          *big.Int `json:"blockNumber"`
	Timestamp            uint64   `json:"timestamp"`
	From                 string   `json:"from"`
	To                   string   `json:"to"`
	Value                *big.Int `json:"value"`
	GasLimit             string   `json:"gasLimit"` // todo `json:"gas"`
	GasUsed              uint64   `json:"gasUsed"`
	GasPrice             *big.Int `json:"gasPrice"`
	BaseFeePerGas        *big.Int `json:"baseFeePerGas"`
	MaxFeePerGas         *big.Int `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *big.Int `json:"maxPriorityFeePerGas"`
	Burnt                string   `json:"burnt"`          // todo
	TxnSavingsFees       string   `json:"txnSavingsFees"` // todo
	BlockHash            string   `json:"blockHash"`
	InputData            string   `json:"input"`
	Nonce                uint64   `json:"nonce"`
	TxnType              uint8    `json:"type"`
	Size                 uint64   `json:"size"`
	Cost                 *big.Int `json:"cost"`
	//ChainId              string        `json:"chainId"`
	//AccessList           []interface{} `json:"accessList"`
	//V string `json:"v"`
	//R string `json:"r"`
	//S string `json:"s"`
}

type Block struct {
	Number    *big.Int `json:"number"`
	Timestamp uint64   `json:"timestamp"`
	Miner     string   `json:"miner"`
	//TotalDifficulty string     `json:"totalDifficulty"`
	Size            uint64   `json:"size"`
	GasUsed         uint64   `json:"gasUsed"`
	GasLimit        uint64   `json:"gasLimit"`
	BaseFeePerGas   *big.Int `json:"baseFeePerGas"`
	BurntFees       *big.Int `json:"burntFees"`
	ExtraData       string   `json:"extraData"`
	Hash            string   `json:"hash"`
	ParentHash      string   `json:"parentHash"`
	StateRoot       string   `json:"stateRoot"`
	WithdrawalsRoot string   `json:"withdrawalsRoot"`
	Nonce           uint64   `json:"nonce"`
	Difficulty      *big.Int `json:"difficulty"`

	//LogsBloom    string `json:"logsBloom"`
	//MixHash      string `json:"mixHash"`
	//ReceiptsRoot string `json:"receiptsRoot"`
	//Sha3Uncles   string `json:"sha3Uncles"`
	//Transactions []struct {} `json:"transactions"`
	//TransactionsRoot string        `json:"transactionsRoot"`
	//Uncles           []interface{} `json:"uncles"`
	//Withdrawals      []struct {} `json:"withdrawals"`
}
