package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Address            string `json:"address" gorm:"column:address;unique"`
	Email            string `json:"email" gorm:"column:email;"`
}

func (User) TableName() string {
	return "users"
}
