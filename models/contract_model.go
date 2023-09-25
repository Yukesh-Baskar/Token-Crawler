package models

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
)

type ContractDetails struct {
	Client          *ethclient.Client `json:"client" validate:"omitempty"`
	ContractAddress string            `json:"contract_address" validate:"required"`
	RPC             string            `json:"network_rpc" validate:"required"`
	Network         string            `json:"network" validate:"required"`
	DeployedBlock   int               `json:"deployed_block" validate:"required"`
}

type ContractDetails struct {
	Client                *ethclient.Client `json:"client" validate:"omitempty"`
	RPC                   string            `json:"network_rpc" validate:"required"`
	OwnerAddress          string            `json:"owner_address" validate:"required"`
	ContractAddress       string            `json:"kp_address" validate:"required"`
	TokenAddressToMigrate string            `json:"token_address" validate:"required"`
	Network               string            `json:"network" validate:"required"`
}

type RedisClient struct {
	Client *redis.Client
}
