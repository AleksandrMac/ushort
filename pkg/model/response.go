package model

import "fmt"

type SQLResponse uint8

const (
	SQLNoResult SQLResponse = iota
)

var SQLResult = map[SQLResponse]error{
	SQLNoResult: fmt.Errorf("sql: no rows in result set"),
}

type ErrorResponse uint8
