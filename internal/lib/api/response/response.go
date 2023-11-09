package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errors validator.ValidationErrors) Response {
	var errorMessages []string

	for _, err := range errors {
		switch err.ActualTag() {
		case "required":
			errorMessages = append(errorMessages, fmt.Sprintf("Field %s is required", err.Field()))
		case "url":
			errorMessages = append(errorMessages, fmt.Sprintf("Field %s is not valid URL", err.Field()))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("Field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errorMessages, ", "),
	}
}
