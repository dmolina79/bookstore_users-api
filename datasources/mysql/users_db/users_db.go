package users_db

import (
	"database/sql"
	"github.com/dmolina79/bookstore_users-api/config"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	Client *sql.DB
)

func init() {
	var err error
	Client, err := sql.Open("mysql", config.GetUsersDBDataSourceName())
	if err != nil {
		panic(err)
	}
	if err := Client.Ping(); err != nil {
		panic(err)
	}
	log.Println("connected to database")

}
