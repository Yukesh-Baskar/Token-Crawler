package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/crawler-tokens/migration/database"
	"github.com/crawler-tokens/migration/helpers"
	"github.com/crawler-tokens/migration/models"
	errorshandler "github.com/crawler-tokens/migration/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckGetUsersInput(ctx *gin.Context) {
	var contractDetails models.ContractDetails
	if err := ctx.BindJSON(&contractDetails); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var e *errorshandler.NewError

	if strings.TrimSpace(contractDetails.ContractAddress) == "" || strings.TrimSpace(contractDetails.RPC) == "" || strings.TrimSpace(contractDetails.Network) == "" || contractDetails.DeployedBlock == 0 {
		e = errorshandler.GetErrors(errors.New("contract address, network, rpc or deployed block cant be empty"), http.StatusPartialContent)
		ctx.AbortWithStatusJSON(e.Status, gin.H{
			"error": e.HandleError().Error,
		})
		return
	}
	err := helpers.CheckIsHexAddress(contractDetails, ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.HandleError().Error,
		})
		return
	}
	ctx.Set("contractDetails", &contractDetails)
	ctx.Next()
}

func CheckIsCollectionExist(ctx *gin.Context) {
	var contractDetails models.ContractDetails

	if err := ctx.BindJSON(&contractDetails); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	fmt.Println("contractDetails", contractDetails)
	err := helpers.CheckIsHexAddress(contractDetails, ctx)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("%v", err.Error),
		})
		return
	}

	var token_collection *mongo.Collection = database.OpenCollection(contractDetails.Network)
	cNames, cErr := token_collection.Database().ListCollectionNames(context.Background(), bson.D{})

	if cErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	isCollectionExist := false
	for _, cname := range cNames {
		if cname == contractDetails.Network {
			isCollectionExist = true
			break
		}
	}

	if !isCollectionExist {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "collection not exist in db.",
		})
		return
	}

	ctx.Set("migrationContractDetails", &contractDetails)
	ctx.Next()
}

func CheckIsDataExist(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	network := ctx.Query("network")

	if page < 1 || network == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "page can't be < 1 or network can't be empty!",
		})
		return
	}

	ctx.Set("details", []interface{}{page, network})

	ctx.Next()
}

func CheckIsContractCollectionExist(ctx *gin.Context) {
	contractAddress := ctx.Query("contract")
	network := ctx.Query("network")

	if strings.TrimSpace(contractAddress) == "" || strings.TrimSpace(network) == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "contract and network can't be empty!",
		})
	}

	if !common.IsHexAddress(contractAddress) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "not an contract address!",
		})
	}

	values := []string{contractAddress, network}

	ctx.Set("contractAddress", values)

	ctx.Next()
}
