package app

import "github.com/gin-gonic/gin"

var (
	router = gin.Default()
)


func Start() {
	setupRoutes()
	router.Run(":8080")
}
