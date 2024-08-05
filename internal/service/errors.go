package service

import "errors"

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
	ErrSignToken     = errors.New("can't sign token")
	ErrParseToken    = errors.New("can't parse token")
)
