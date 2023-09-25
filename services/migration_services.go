package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/crawler-tokens/migration/database"
	token "github.com/crawler-tokens/migration/gen"
	multicall "github.com/crawler-tokens/migration/gen-multicall"
	"github.com/crawler-tokens/migration/helpers"
	"github.com/crawler-tokens/migration/models"
	errorshandler "github.com/crawler-tokens/migration/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *models.RedisClient

func ConsumerService(chn chan any, done chan<- any, cd *ethclient.Client, network string) {
	var token_collection *mongo.Collection = database.OpenCollection(network)
	for {
		res, ok := <-chn
		if !ok {
			done <- true
		} else {
			r, _ := res.(models.TransactionDetails)
			fromAddress := toCheckSummedAddress(r.FromAddress)
			toAddress := toCheckSummedAddress(r.ToAddress)
			result := helpers.IsContract([]string{fromAddress, toAddress}, cd)
			var user *models.User
			if fromAddress == toAddress || result[0] && result[1] {
				continue
			} else if !result[0] && result[1] { // check if from is EOA and to is contract. If it's true -from don't +to
				filter := primitive.M{"useraddress": fromAddress}
				count := getCount(token_collection, filter, done)
				injectOrUpdateTo(count, user, token_collection, done, r.AmountTransferred, network, fromAddress, 1, r.LatestBlock) // flag 1 -> if from is EOA -FROM
			} else if result[0] && !result[1] { // check if from is contract and to is EOA, if it's true +to don't -from because we're not gonna store contract details
				filter := primitive.M{"useraddress": toAddress}
				count := getCount(token_collection, filter, done)
				injectOrUpdateTo(count, user, token_collection, done, r.AmountTransferred, network, toAddress, 2, r.LatestBlock)
			} else {
				fromFilter := primitive.M{"useraddress": fromAddress}
				toFilter := primitive.M{"useraddress": toAddress}
				fromCount := getCount(token_collection, fromFilter, done)
				toCount := getCount(token_collection, toFilter, done)
				var fromUser *models.User
				var toUser *models.User
				injectOrUpdateTo(fromCount, fromUser, token_collection, done, r.AmountTransferred, network, fromAddress, 1, r.LatestBlock)
				injectOrUpdateTo(toCount, toUser, token_collection, done, r.AmountTransferred, network, toAddress, 2, r.LatestBlock)
			}
		}
	}
}

func GetLatestBlockService() *models.User {
	var user *models.User

	var token_collection *mongo.Collection

	token_collection.FindOne(context.Background(), bson.M{"$natural": -1}).Decode(&user)

	return user
}

func MigrateUsersService(auth *bind.TransactOpts, fromAddress common.Address, cd models.ToNetworkContractDetails, multicallContractInstance *multicall.Multicall, tokenContractInstance *token.Token, network string, tokenAddress string) *errorshandler.NewError {
	var token_collection *mongo.Collection = database.OpenCollection(network)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	totalCount, _ := token_collection.CountDocuments(ctx, bson.D{})
	var countToLimit int64 = 100
	var countToSkip int64 = 0
	var totalDataPushed int64 = 0

	if totalCount < int64(countToLimit) {
		countToLimit = int64(totalCount)
	}

	for {
		if countToSkip == totalCount {
			return nil
		}
		dataChan := make(chan any)
		go produceData(dataChan, token_collection, network, &countToLimit, &countToSkip, &totalCount)
		data := <-dataChan
		res, ok := data.([]models.User)
		if !ok {
			return &errorshandler.NewError{
				Error:  errors.New("not a user type"),
				Status: http.StatusInternalServerError,
			}
		}
		// migrate these users to kp
		trx, err := consumeData(res, auth, fromAddress, multicallContractInstance, tokenContractInstance, &totalDataPushed, tokenAddress, &totalCount, token_collection)
		if err != nil {
			return &errorshandler.NewError{
				Error:  err,
				Status: http.StatusInternalServerError,
			}
		}
		fmt.Println("trx", trx)
		nonce, nonceErr := cd.Client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			return &errorshandler.NewError{
				Error:  errors.New(nonceErr.Error()),
				Status: http.StatusInternalServerError,
			}
		}
		auth.Nonce = big.NewInt(int64(nonce))
	}
}

func consumeData(data []models.User, auths *bind.TransactOpts, fromAddress common.Address, multiCallInstance *multicall.Multicall, tokenContractInstance *token.Token, totalDataPused *int64, tokenAddress string, totalCount *int64, token_collection *mongo.Collection) (*types.Transaction, []error) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	var errorsList = new([]error)
	var err error
	var trx *types.Transaction
	go func(wg *sync.WaitGroup) {
		addressess, amounts := getData(data, fromAddress)
		token := common.HexToAddress(tokenAddress)
		trx, err = multiCallInstance.MultiSigcall(auths, token, fromAddress, addressess, amounts)
		if err != nil {
			*errorsList = append(*errorsList, err)
		}
		time.Sleep(5 * time.Second)
		for _, userData := range data {
			convertedAddress := common.HexToAddress(userData.UserAddress)
			balance, _ := tokenContractInstance.BalanceOf(&bind.CallOpts{}, convertedAddress)
			userDatabaseBalance, _ := new(big.Int).SetString(userData.TokenAmount, 0)
			fmt.Println("balance, userDatabaseBalance", convertedAddress, balance, userDatabaseBalance, userDatabaseBalance.String() == balance.String())
			if userDatabaseBalance.String() != balance.String() {
				*errorsList = append(*errorsList, fmt.Errorf("user balance mismatched or trx may be failed for address and amount: %v %v %v", userData.UserAddress, userData.TokenAmount, userDatabaseBalance))
			}
			Client.Client.Set(context.Background(), userData.UserAddress, true, time.Minute*5)
		}
		*totalDataPused += int64(len(data))
		wg.Done()
	}(&wg)
	wg.Wait()
	if len(*errorsList) > 0 {
		return nil, *errorsList
	}
	size, _ := Client.Client.DBSize(context.Background()).Result()
	if size >= int64(300) || *totalDataPused == *totalCount {
		keysCmd := Client.Client.Keys(context.Background(), "*")
		keys, err := keysCmd.Result()
		if err != nil {
			*errorsList = append(*errorsList, err)
			return nil, *errorsList
		}

		filter := bson.D{{Key: "useraddress", Value: bson.M{"$in": keys}}}
		update := bson.M{"$set": bson.M{"istokenmigrated": true}}

		res, err := token_collection.UpdateMany(context.Background(), filter, update)

		if err != nil {
			*errorsList = append(*errorsList, err)
			return nil, *errorsList
		}

		fmt.Printf("%+v \n", res)
		fmt.Println(Client.Client.FlushAll(context.Background()))
	}
	return trx, nil
}

func getData(data []models.User, fromAddress common.Address) ([]common.Address, []*big.Int) {
	var toAddressess []common.Address
	var amounts []*big.Int
	for i := 0; i < len(data); i++ {
		convertedAddress := common.HexToAddress(data[i].UserAddress)
		toAddressess = append(toAddressess, convertedAddress)
		userTokenAmount, _ := new(big.Int).SetString(data[i].TokenAmount, 0)
		amounts = append(amounts, userTokenAmount)
	}
	return toAddressess, amounts
}

func produceData(dataChan chan<- any, token_collection *mongo.Collection, network string, countToLimit *int64, countToSkip *int64, totalCount *int64) {
	var users []models.User
	// var length int64
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var toLimit int64 = (*totalCount - *countToSkip)
	if *countToLimit > *totalCount || toLimit > 100 {
		toLimit = *countToLimit
	}

	skipStage := bson.D{{Key: "$skip", Value: *countToSkip}}
	limitStage := bson.D{{Key: "$limit", Value: toLimit}}

	cursor, err := token_collection.Aggregate(ctx, mongo.Pipeline{skipStage, limitStage})

	if err != nil {
		dataChan <- errorshandler.NewError{
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		dataChan <- errorshandler.NewError{
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	*countToSkip += *countToLimit
	if *countToSkip > *totalCount || countToSkip == totalCount {
		*countToSkip = *totalCount
	}

	defer cursor.Close(ctx)

	dataChan <- (users)
}

func toCheckSummedAddress(addr string) string {
	hexAddress := common.HexToAddress(addr)
	userAddress := []common.Address{hexAddress}
	return userAddress[0].String()
}

func getCount(token_collection *mongo.Collection, filter primitive.M, done chan<- any) int64 {
	count, err := token_collection.CountDocuments(context.Background(), filter)
	if err != nil {
		e := errorshandler.GetErrors(err, http.StatusInternalServerError)
		done <- e.HandleError()
	}
	return count
}

func injectOrUpdateTo(count int64, user *models.User, token_collection *mongo.Collection, done chan<- any, amountT big.Int, network, addr string, flag int32, block uint64) {
	filter := bson.M{"useraddress": addr}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	var e *errorshandler.NewError
	var tokenAmount *big.Int
	if amountT.String() == "0" {
		fmt.Printf("eat fivestart do nothing %v \n", amountT.String())
	} else if count > 0 {
		if err := token_collection.FindOne(ctx, filter).Decode(&user); err != nil {
			e = errorshandler.GetErrors(err, http.StatusNotFound)
			done <- e.HandleError()
		}

		userTokenAmount, _ := new(big.Int).SetString(user.TokenAmount, 0)
		if flag == 1 {
			tokenAmount = userTokenAmount.Sub(userTokenAmount, &amountT)
		} else {
			tokenAmount = userTokenAmount.Add(userTokenAmount, &amountT)
		}

		if tokenAmount.String() == "0" {
			// deleting the user because he's not the holder any more once the token balance drains to 0
			_, err := token_collection.DeleteOne(ctx, filter)
			if err != nil {
				e := errorshandler.GetErrors(err, http.StatusInternalServerError)
				done <- e.HandleError()
			}
		} else {
			upsert := true
			opts := options.UpdateOptions{
				Upsert: &upsert,
			}
			_, err := token_collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.M{"tokenamount": tokenAmount.String(), "latestblock": block}}}, &opts)
			if err != nil {
				e := errorshandler.GetErrors(err, http.StatusInternalServerError)
				done <- e.HandleError()
			}
		}
	} else {
		user := &models.User{
			UserAddress: addr,
			TokenAmount: amountT.String(),
			LatestBlock: block,
			Network:     network,
		}
		_, err := token_collection.InsertOne(ctx, user)
		if err != nil {
			e = errorshandler.GetErrors(err, http.StatusInternalServerError)
			done <- e.HandleError()
		}
	}
}

func GetUsersService(page interface{}, network string) (interface{}, error) {

	var users []models.User
	var token_collection *mongo.Collection = database.OpenCollection(network)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)

	defer cancel()

	limitPerPage := 100
	skip := 0
	pageNumber := 1

	if page != 1 {
		pageNumber = page.(int)
		skip = (pageNumber - 1) * limitPerPage
	}

	totalRows, err := token_collection.CountDocuments(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	totalPages := math.Ceil(float64(totalRows/int64(limitPerPage)) + 1)

	if pageNumber > int(totalPages) {
		return nil, fmt.Errorf("page %v is > totalPages %v", pageNumber, totalPages)
	}

	skipStage := bson.D{{Key: "$skip", Value: skip}}
	limitStage := bson.D{{Key: "$limit", Value: limitPerPage}}

	cursor, err := token_collection.Aggregate(ctx, mongo.Pipeline{skipStage, limitStage})

	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	return gin.H{
		"users":                      users,
		"totalPages":                 totalPages,
		"pageNumber":                 pageNumber,
		"totalCountInTheCurrentPage": len(users),
	}, nil
}
