package users

import (
	"database/sql"
	"fmt"
	"github.com/dmolina79/bookstore_users-api/datasources/mysql/users_db"
	"github.com/dmolina79/bookstore_users-api/utils/date"
	"github.com/dmolina79/bookstore_users-api/utils/errors"
	"strings"
)

const (
	uniqueEmailError = "email_UNIQUE"
	errorNoRows = "no rows in result set"
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?,?,?,?);"
	queryGetUser = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?;"
)

var (
	usersDB = make(map[int64]*User)
)

func mapUserFromRow(row *sql.Row, user *User) error {
	return row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated)
}

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServer(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(user.Id)
	if err := mapUserFromRow(row, user); err != nil {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFound(fmt.Sprintf("user %d not found", user.Id))
		}
		fmt.Println(err)
		return errors.NewInternalServer(fmt.Sprintf("error when trying to get user %d %s", user.Id, err.Error()))
	}

	return nil
}



func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServer(err.Error())
	}
	defer stmt.Close()

	user.DateCreated = date.GetNowString()

	insertResult, err := stmt.Exec(
		user.FirstName,
		user.LastName,
		user.Email,
		user.DateCreated,
	)

	if err != nil {
		if strings.Contains(err.Error(), uniqueEmailError) {
			return errors.NewBadRequest(fmt.Sprintf("email %s is already registered", user.Email))
		}
		return errors.NewInternalServer(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServer(fmt.Sprintf("error when trying to get last insert id: %s", err.Error()))
	}
	user.Id = userId

	return nil
}
