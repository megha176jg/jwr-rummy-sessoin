package errors

import "errors"

var (
	ErrInvalidUserId      = errors.New("Invalid User Id")
	ErrCreatingToken      = errors.New("Error in creating token")
	ErrValidation         = errors.New("Error in validating")
	ErrIncorrectAuthtoken = errors.New("Incorrect Auth Token")
	ErrInvalidation       = errors.New("Error in logging out")
)
