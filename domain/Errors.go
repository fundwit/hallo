package domain

import "errors"

var (
	IdGenerateFailure = errors.New("failed to generate a new id")
)

type AccountNameIsOccupied struct {
}

func (e *AccountNameIsOccupied) Error() string {
	return "account email is occupied"
}

type AccountEmailIsOccupied struct {
}

func (e *AccountEmailIsOccupied) Error() string {
	return "account email is occupied"
}

type AccountAuthenticationFailure struct {
}

func (e *AccountAuthenticationFailure) Error() string {
	return "user is not exist or credential is invalid"
}

type ErrUnauthorized struct {
}

func (e *ErrUnauthorized) Error() string {
	return "unauthorized"
}

type ErrRegisterTokenInvalid struct {
}

func (e *ErrRegisterTokenInvalid) Error() string {
	return "register.token.is.invalid"
}
