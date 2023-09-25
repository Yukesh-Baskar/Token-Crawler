package configurations

import (
	"context"
	"log"
	"time"

	"github.com/crawler-tokens/migration/database"
	"github.com/crawler-tokens/migration/models"
	routes "github.com/crawler-tokens/migration/routes"
	"github.com/crawler-tokens/migration/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var RPC_URL string

func StartApp() {
	viper.SetConfigName("config")

	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	PORT, ok := viper.Get("PORT").(string)

	if !ok {
		log.Fatalf("invalid type: %v, %T \n", ok, PORT)
	}

	err = NewRedisClient()

	if err != nil {
		log.Fatalf("error occured while connecting redis: %v", err)
	}

	app := gin.Default()
	app.Use(cors.Default())

	routes.HandleRoutes(app)
	database.ConnectToDataBase()
	app.Run(PORT)
}

func NewRedisClient() error {
	redisPort, _ := viper.Get("REDIS_PORT").(string)

	client := redis.NewClient(&redis.Options{
		Addr:        redisPort,
		DB:          0,
		Username:    "user",
		DialTimeout: 100 * time.Millisecond,
		ReadTimeout: 100 * time.Millisecond,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return err
	}
	services.Client = &models.RedisClient{
		Client: client,
	}

	return nil
}
