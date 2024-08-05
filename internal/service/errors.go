package service

import "errors"

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
	ErrSignToken     = errors.New("can't sign token")
	ErrParseToken    = errors.New("can't parse token")

	ErrHouseExists   = errors.New("house already exists")
	ErrHouseNotFound = errors.New("house not found")

	ErrFlatExists   = errors.New("flat already exists")
	ErrFlatNotFound = errors.New("flat not found")

	ErrHouseSubscriptionExists = errors.New("house subscription already exists")
)
