package mysql_utils

import (
	"github.com/dmolina79/bookstore_users-api/utils/errors"
	"github.com/go-sql-driver/mysql"
	"strings"
)

const (
	DuplicateKeyError = 1062
	errorNoRows       = "no rows in result set"
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFound("no records matching given id")
		}
		return errors.NewInternalServer("error parsing database response")
	}

	switch sqlErr.Number {
	case DuplicateKeyError:
		return errors.NewBadRequest("duplicated key")
	}

	return errors.NewInternalServer("error processing request")
}
