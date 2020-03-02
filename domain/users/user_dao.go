package users

import (
	"database/sql"
	"fmt"
	"github.com/dmolina79/bookstore_users-api/datasources/mysql/users_db"
	"github.com/dmolina79/bookstore_users-api/logger"
	"github.com/dmolina79/bookstore_users-api/utils/errors"
)

const (
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?,?,?,?,?,?);"
	queryGetUser          = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	queryUpdateUser       = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser       = "DELETE from users where id=?;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users where status=?;"
)

var (
	usersDB = make(map[int64]*User)
)

func mapUserFromRow(row *sql.Row, user *User) error {
	return row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)
}

func (u *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServer("database error")
	}
	defer stmt.Close()

	row := stmt.QueryRow(u.Id)
	if err := mapUserFromRow(row, u); err != nil {
		logger.Error("error when trying to map get user statement", err)
		return errors.NewInternalServer("database error")
	}

	return nil
}

func (u *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return errors.NewInternalServer("database error")
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(
		u.FirstName,
		u.LastName,
		u.Email,
		u.DateCreated,
		u.Status,
		u.Password,
	)

	if saveErr != nil {
		logger.Error("error when trying to execute save user statement", err)
		return errors.NewInternalServer("database error")
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return errors.NewInternalServer("database error")
	}
	u.Id = userId

	return nil
}

func (u *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors.NewInternalServer("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.FirstName, u.LastName, u.Email, u.Id)
	if err != nil {
		logger.Error("error when trying to execute update user statement", err)
		return errors.NewInternalServer("database error")
	}
	return nil
}

func (u *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors.NewInternalServer("database error")
	}
	defer stmt.Close()

	if _, err := stmt.Exec(u.Id); err != nil {
		logger.Error("error when trying to execute delete user statement", err)
		return errors.NewInternalServer("database error")
	}

	return nil
}

func (u *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find user statement", err)
		return nil, errors.NewInternalServer("database error")
	}
	defer stmt.Close()
	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to execute find user statement", err)
		return nil, errors.NewInternalServer("database error")
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when trying to scan find user result", err)
			return nil, errors.NewInternalServer("database error")
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFound(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil

}
