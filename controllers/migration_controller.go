package controllers

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	token "github.com/crawler-tokens/migration/gen"
	multicall "github.com/crawler-tokens/migration/gen-multicall"
	"github.com/crawler-tokens/migration/helpers"
	"github.com/crawler-tokens/migration/models"
	"github.com/crawler-tokens/migration/services"
	errorshandler "github.com/crawler-tokens/migration/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var LogTransferSigHash = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

func GetUsers(ctx *gin.Context) {

	contractChannel := make(chan interface{}, 1)
	contractDetails, _ := ctx.Value("contractDetails").(*models.ContractDetails)
	go helpers.InitializeContract(contractDetails, contractChannel)
	contractDetail := <-contractChannel
	defer ctx.Request.Body.Close()

	cd, ok := contractDetail.(models.ContractDetails)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": contractDetail,
		})
		return
	}

	contractHexAddr := common.HexToAddress(cd.ContractAddress)

	blocksToIncreasePerSpin := 5000
	fromBlock := cd.DeployedBlock
	toBlock := fromBlock + blocksToIncreasePerSpin

	currentBlockNumber, err := cd.Client.BlockNumber(context.TODO())

	if err != nil {
		e := errorshandler.GetErrors(err, http.StatusInternalServerError)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, e)
		return
	}

	for {
		if toBlock > int(currentBlockNumber) {
			toBlock = int(currentBlockNumber)
		}
		filterQuery := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(fromBlock)),
			ToBlock:   big.NewInt(int64(toBlock)),
			Addresses: []common.Address{
				contractHexAddr,
			},
		}
		logs, err := cd.Client.FilterLogs(context.Background(), filterQuery)
		if err != nil {
			e := errorshandler.GetErrors(err, http.StatusInternalServerError)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, e)
			return
		}
		if len(logs) > 0 {
			chn := make(chan any)
			done := make(chan any)
			go producerController(chn, logs)
			go services.ConsumerService(chn, done, cd.Client, contractDetails.Network)

			response := <-done
			_, ok := response.(errorshandler.NewError)
			if ok {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
				return
			}
		}
		if toBlock == int(currentBlockNumber) {
			// wait until all the trx in the block to update.
			time.Sleep(time.Second * 60)
			ctx.JSON(http.StatusOK, gin.H{
				"message": "success",
			})
			break
		}
		fromBlock = toBlock + 1
		toBlock += blocksToIncreasePerSpin
	}
}

func LatestBlockInfo(ctx *gin.Context) {
	defer ctx.Request.Body.Close()
	latestUser := services.GetLatestBlockService()
	ctx.JSON(http.StatusOK, gin.H{
		"message": latestUser,
	})
}

func Migrate(ctx *gin.Context) {
	contractChannel := make(chan interface{})
	contractDetailsToMigrate, _ := ctx.Value("migrationContractDetails").(*models.ContractDetails)

	go helpers.InitializeMigrationContractDetails(contractDetailsToMigrate, contractChannel)
	migratingContractDetails := <-contractChannel
	defer ctx.Request.Body.Close()

	cd, ok := migratingContractDetails.(models.ContractDetails)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": migratingContractDetails,
		})
		return
	}

	tokenContractAddress := common.HexToAddress(cd.ContractAddress)
	multiCallContractAddress := common.HexToAddress(cd.ContractAddress)

	tokenContractInstance, err := token.NewToken(tokenContractAddress, cd.Client)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	multiCallContractInstance, err := multicall.NewMulticall(multiCallContractAddress, cd.Client)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	pkHex := viper.Get("DUMMY_PK").(string)

	privateKey, err := crypto.HexToECDSA(pkHex)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		e := errorshandler.GetErrors(errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey"), http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, e)
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := cd.Client.PendingNonceAt(context.Background(), fromAddress)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	gasPrice, err := cd.Client.SuggestGasPrice(context.Background())

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	chainId, err := cd.Client.ChainID(context.Background())

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	auths, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	auths.Nonce = big.NewInt(int64(nonce)) // number only once should be incremented for every trx so this will be updated on the upcoming trx.
	auths.Value = big.NewInt(0)
	auths.GasLimit = uint64(100 * 36000)
	auths.GasPrice = gasPrice

	serviceErr := services.MigrateUsersService(auths, fromAddress, cd, multiCallContractInstance, tokenContractInstance, contractDetailsToMigrate.Network, contractDetailsToMigrate.ContractAddress)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%+v", serviceErr.Error),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully migrated to kp",
	})
}

func GetUsersByPage(ctx *gin.Context) {
	queryParams := ctx.Value("details").([]interface{})

	data, err := services.GetUsersService(queryParams[0].(int), queryParams[1].(string))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func producerController(chn chan any, logs []types.Log) {
	for key, log := range logs {
		if log.Topics[0].Hex() == LogTransferSigHash.Hex() {
			fromAddress, toAddress, amount := extractAddressAndAmount(log.Topics[1].String(), log.Topics[2].String(), log.Data)
			user := models.TransactionDetails{
				FromAddress:       fromAddress,
				ToAddress:         toAddress,
				LatestBlock:       log.BlockNumber,
				AmountTransferred: *amount,
			}
			chn <- user
		}
		if key+1 == len(logs) {
			close(chn)
		}
	}
}

func extractAddressAndAmount(fromAddress, toAddress string, amounts []byte) (f, t string, amount *big.Int) {
	f = "0x" + fromAddress[26:]
	t = "0x" + toAddress[26:]
	res := hex.EncodeToString(amounts)
	tokenAmount, _ := new(big.Int).SetString("0x"+res, 0)
	return f, t, tokenAmount
}
