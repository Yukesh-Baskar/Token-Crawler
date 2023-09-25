package models

import (
	"math/big"
)

type User struct {
	UserAddress     string `json:"useraddress"`
	TokenAmount     string `json:"tokenamount"`
	Network         string `json:"network_name"`
	LatestBlock     uint64  `json:"latesblock"`
	IsTokenMigrated bool   `json:"istoken_migrated"`
}

type TransactionDetails struct {
	FromAddress       string
	ToAddress         string
	LatestBlock       uint64 `json:"latesblock"`
	AmountTransferred big.Int
}
