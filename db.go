package merrors

import (
	"regexp"
)

const (
	ErrorDB             ErrorType = "db_error"
	ErrorDuplicateEntry ErrorType = "duplicate_entry"
)

var duplicateEntryRegExp = regexp.MustCompile("Error 1062")

func DBError(e error) CommonError {
	switch {
	case duplicateEntryRegExp.MatchString(e.Error()):
		return ErrorBadRequest(e.Error(), ErrorDuplicateEntry)
	default:
		return ErrorInternalServerError(e.Error(), ErrorDB)
	}
}
