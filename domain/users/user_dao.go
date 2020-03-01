package users

import (
	"database/sql"
	"fmt"
	"github.com/dmolina79/bookstore_users-api/datasources/mysql/users_db"
	"github.com/dmolina79/bookstore_users-api/utils/errors"
	"github.com/dmolina79/bookstore_users-api/utils/mysql_utils"
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
		return errors.NewInternalServer(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(u.Id)
	if err := mapUserFromRow(row, u); err != nil {
		return mysql_utils.ParseError(err)
	}

	return nil
}

func (u *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServer(err.Error())
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
		return mysql_utils.ParseError(saveErr)
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	u.Id = userId

	return nil
}

func (u *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServer(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.FirstName, u.LastName, u.Email, u.Id)
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (u *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return errors.NewInternalServer(err.Error())
	}
	defer stmt.Close()

	if _, err := stmt.Exec(u.Id); err != nil {
		return mysql_utils.ParseError(err)
	}

	return nil
}

func (u *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		return nil, errors.NewInternalServer(err.Error())
	}
	defer stmt.Close()
	rows, err := stmt.Query(status)
	if err != nil {
		return nil, errors.NewInternalServer(err.Error())
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFound(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil

}
