package users

import (
	"database/sql"
	"github.com/dmolina79/bookstore_users-api/datasources/mysql/users_db"
	"github.com/dmolina79/bookstore_users-api/utils/date"
	"github.com/dmolina79/bookstore_users-api/utils/errors"
	"github.com/dmolina79/bookstore_users-api/utils/mysql_utils"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?,?,?,?);"
	queryGetUser    = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?;"
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
		return mysql_utils.ParseError(err)
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

	insertResult, saveErr := stmt.Exec(
		user.FirstName,
		user.LastName,
		user.Email,
		user.DateCreated,
	)

	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	user.Id = userId

	return nil
}
