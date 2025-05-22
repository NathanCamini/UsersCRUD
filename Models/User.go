package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ID uuid.UUID

type User struct {
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname"  validate:"required"`
	Bio       string `json:"bio"       validate:"required"`
}

var Validate = validator.New()
