package errors

import (
    "errors"
    "net/http"
)

var New = errors.New

var ErrorUsernameInUse = errors.New("Username already in use!")
var ErrorUserNotFound = errors.New("User not found!")
var ErrorSomethingWentWrong = errors.New("Something went wrong!")
var ErrorTournamentNotFound = errors.New("Tournament not found!")
var ErrorTournamentEnded = errors.New("Tournament already ended!")
var ErrorTournamentMaxed = errors.New("Tournament already has maximum registrations!")
var ErrorNotEnoughFunds = errors.New("Not enough funds!")

var Codes = map[error]int{
    ErrorUsernameInUse: http.StatusBadRequest,
    ErrorUserNotFound: http.StatusNotFound,
    ErrorSomethingWentWrong: http.StatusInternalServerError,
    ErrorTournamentNotFound: http.StatusNotFound,
    ErrorTournamentEnded: http.StatusBadRequest,
    ErrorTournamentMaxed: http.StatusBadRequest,
    ErrorNotEnoughFunds: http.StatusBadRequest,
}
