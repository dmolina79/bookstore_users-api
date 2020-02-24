package users

import (
	"fmt"
	"github.com/dmolina79/bookstore_users-api/datasources/mysql/users_db"
	"github.com/dmolina79/bookstore_users-api/utils/date"
	"github.com/dmolina79/bookstore_users-api/utils/errors"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {
	if err := users_db.Client.Ping(); err != nil {
		panic(err)
	}
	result := usersDB[user.Id]
	if result == nil {
		return errors.NewNotFound(fmt.Sprintf("user %d not found", user.Id))
	}
	user.Id = result.Id
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreated = result.DateCreated

	return nil
}

func (user *User) Save() *errors.RestErr {
	current := usersDB[user.Id]
	if current != nil {
		if current.Email == user.Email {
			return errors.NewBadRequest(fmt.Sprintf("email %s already registered", user.Email))
		}
		return errors.NewBadRequest(fmt.Sprintf("user %d already exists", user.Id))
	}
	user.DateCreated = date.GetNowString()
	usersDB[user.Id] = user

	return nil

}
