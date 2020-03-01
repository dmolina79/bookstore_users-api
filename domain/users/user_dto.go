package users

import (
	"github.com/dmolina79/bookstore_users-api/utils/errors"
	"strings"
)

const (
	StatusActive = "active"
)

type Users []User

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"-"`
}

func (u *User) Validate() *errors.RestErr {
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	if u.Email == "" {
		return errors.NewBadRequest("invalid email address")
	}

	u.Password = strings.TrimSpace(u.Password)

	if u.Password == "" {
		return errors.NewBadRequest("invalid password for user")
	}

	return nil
}
