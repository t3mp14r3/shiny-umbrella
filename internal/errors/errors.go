package errors

import (
    "errors"
    "net/http"
)

var New = errors.New

var ErrorUsernameInUse = errors.New("Username already in use!")
var ErrorUserNotFound = errors.New("User not found!")
var ErrorSomethingWentWrong = errors.New("Something went wrong!")

var Codes = map[error]int{
    ErrorUsernameInUse: http.StatusBadRequest,
    ErrorUserNotFound: http.StatusNotFound,
    ErrorSomethingWentWrong: http.StatusInternalServerError,
}
