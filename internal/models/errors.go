package models

import "errors"

// Put every error when communicating with DB in here
var ErrNoRecord  =  errors.New("models:  no  matching record found")

var ErrInvalidCredentials = errors.New("models: Invalid credentials")

var ErrDuplicateEmail error = errors.New("models: Duplicate email")
