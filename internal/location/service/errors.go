package service

import "errors"

var (
	ErrDriverNotFound     = errors.New("driver is not found")
	ErrRequestIsIncorrect = errors.New("request is incorrect")
	ErrTest               = errors.New("test error")
)
