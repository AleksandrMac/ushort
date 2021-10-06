package model

import (
	"encoding/json"
	"fmt"
)

type SQLResponse uint8

const (
	SQLNoResult SQLResponse = iota
)

var SQLResult = map[SQLResponse]error{
	SQLNoResult: fmt.Errorf("sql: no rows in result set"),
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var ErrorResponseMap = map[int]*ErrorResponse{
	400: {
		Code:    "400",
		Message: "The server could not understand your request",
	},
	401: {
		Code:    "401",
		Message: "Unauthenticated request",
	},
	403: {
		Code:    "403",
		Message: "Unauthorized request",
	},
	404: {
		Code:    "404",
		Message: "Not Found",
	},
	500: {
		Code:    "500",
		Message: "Server error, try again later",
	},
}

func (er *ErrorResponse) JSON() ([]byte, error) {
	return json.Marshal(er)
}
