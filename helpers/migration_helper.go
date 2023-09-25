package helpers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/crawler-tokens/migration/models"
	errorshandler "github.com/crawler-tokens/migration/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

func InitializeContract(contractDetails *models.ContractDetails, contractChannel chan interface{}) {
	client, err := ethclient.Dial(contractDetails.RPC)

	var e *errorshandler.NewError

	if err != nil {
		e = errorshandler.GetErrors(err, http.StatusInternalServerError)
		contractChannel <- e.HandleError()
		return
	}
	contractChannel <- models.ContractDetails{
		Client:          client,
		ContractAddress: contractDetails.ContractAddress,
		RPC:             contractDetails.RPC,
		Network:         contractDetails.Network,
		DeployedBlock:   contractDetails.DeployedBlock,
	}
	close(contractChannel)
}

func InitializeMigrationContractDetails(contractDetails *models.ContractDetails, contractChannel chan interface{}) {
	client, err := ethclient.Dial(contractDetails.RPC)

	var e *errorshandler.NewError

	if err != nil {
		e = errorshandler.GetErrors(err, http.StatusInternalServerError)
		fmt.Println("e", e)
		contractChannel <- e.HandleError()
		return
	}

	contractChannel <- models.ContractDetails{
		Client:                client,
		RPC:                   contractDetails.RPC,
		ContractAddress:       contractDetails.ContractAddress,
		TokenAddressToMigrate: contractDetails.TokenAddressToMigrate,
	}
	close(contractChannel)
}

func IsContract(address []string, cd *ethclient.Client) (res [2]bool) {
	for i := 0; i <= 1; i++ {
		byteCode, _ := cd.CodeAt(context.Background(), common.HexToAddress(address[i]), nil)
		if address[i] == "0x0000000000000000000000000000000000000000" {
			res[i] = true
			continue
		}
		res[i] = len(byteCode) > 0
	}
	return res
}

func CheckIsHexAddress(contractDetails interface{}, c *gin.Context) *errorshandler.NewError {
	cd, ok := contractDetails.(models.ContractDetails)

	if ok {
		if !common.IsHexAddress(cd.ContractAddress) {
			e := errorshandler.GetErrors(errors.New("not a hex address"), http.StatusBadRequest)
			return e
		}
	} else {
		migrateDataContractData := contractDetails.(models.ContractDetails)
		var caDetails = []string{migrateDataContractData.ContractAddress, migrateDataContractData.TokenAddressToMigrate, migrateDataContractData.OwnerAddress}
		for i := 0; i < len(caDetails); i++ {
			if !common.IsHexAddress(caDetails[i]) {
				e := errorshandler.GetErrors(errors.New("not a hex address"), http.StatusBadRequest)
				return e
			}
		}
	}
	return nil
}
