package entity

import "gorm.io/gorm"

type UserEntity struct {
	gorm.Model

	ID       int
	Account  string
	Name     string
	Password string
	Level    int
	Status   int
}
