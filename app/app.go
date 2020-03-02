package app

import (
	"github.com/dmolina79/bookstore_users-api/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func Start() {
	setupRoutes()

	logger.Info("about to start the application")
	router.Run(":8080")
}
