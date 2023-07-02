package models

import (
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Name  		string		`json:"name"`
	Completed	bool		`json:"completed"`
}