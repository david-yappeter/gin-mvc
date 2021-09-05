package service

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrRecordFound    = errors.New("record found")
	ErrRecordNotFound = gorm.ErrRecordNotFound
	ErrInvalidJWT     = errors.New("invalid jwt token")
)
