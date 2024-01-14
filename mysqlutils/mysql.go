package mysqlutils

import (
	"fmt"
	"strings"

	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/go-sql-driver/mysql"
)

const (
	errorNoRows       = "no rows in result set"
	duplicatedKeyCode = 1062
)

// HandleMySQLError - handle mysql errors
func HandleMySQLError(err error) resterrors.RestErr {

	//if exists the error on mysql.MySQLError
	sqlErr, exists := err.(*mysql.MySQLError)
	if !exists {
		if strings.Contains(err.Error(), errorNoRows) {
			return resterrors.NewNotFoundError("No records find with the parameters")
		}
		return resterrors.NewInternalServerError("Error database response", err.Error())
	}

	switch sqlErr.Number {
	case duplicatedKeyCode:
		// example: Error 1062: Duplicate entry 'test@gmail' for key 'users_email_uindex'
		// will return: The email test@gmail already exists
		duplicatedKey := between(sqlErr.Message, "key '", "_UNIQUE")
		duplicatedKeyValue := between(sqlErr.Message, "entry '", "' for key")
		return resterrors.NewConflictError(fmt.Sprintf("The %s %s already exists", duplicatedKey, duplicatedKeyValue))
	}

	return resterrors.NewInternalServerError("Error trying to processing database request", err.Error())
}

func between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}

	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}

	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}

	return value[posFirstAdjusted:posLast]
}

func before(value string, a string) string {
	// Get substring before a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:]
}
