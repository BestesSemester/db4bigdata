package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	RoleID      int `gorm:"index"`
	Description string
}

func (r *Role) Save() {

}
