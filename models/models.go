package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name           string
	Email          string
	HashedPassword string `json:"-"`
}

type SignupRequest struct {
	Name     string
	Email    string
	Password string
}

type LoginRequest struct {
	Email    string
	Password string
}
