package routes

import (
	"github.com/crawler-tokens/migration/controllers"
	"github.com/crawler-tokens/migration/middlewares"
	"github.com/gin-gonic/gin"
)

// no authentication are maded, just migration process.
func HandleRoutes(handler *gin.Engine) {
	handler.POST("/migrate", middlewares.CheckGetUsersInput, controllers.GetUsers)
	handler.GET("/latestblock", controllers.LatestBlockInfo)
	handler.POST("/migratekp", middlewares.CheckIsCollectionExist, controllers.Migrate)
	handler.GET("/users/", middlewares.CheckIsDataExist, controllers.GetUsersByPage)
	handler.GET("/users/export-data/", middlewares.CheckIsContractCollectionExist)
}
