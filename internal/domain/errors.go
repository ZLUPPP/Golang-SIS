package domain

import "errors"

var (
	ErrTaskAlreadyExists = errors.New("task already exists")
	ErrTaskNotFound      = errors.New("task not found")
	ErrTaskInvalidStatus = errors.New("task status invalid")
	ErrEmptyInput        = errors.New("empty input")
	ErrLogFormatInvalid  = errors.New("log format invalid")
)
