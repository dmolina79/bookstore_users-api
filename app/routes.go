package app

import (
	"github.com/dmolina79/bookstore_users-api/controllers/ping"
	"github.com/dmolina79/bookstore_users-api/controllers/users"
)

func setupRoutes() {
	router.GET("/ping", ping.Ping)

	router.POST("/users", users.CreateUser)
	router.GET("/users/:user_id", users.GetUser)
}
