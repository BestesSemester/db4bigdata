package model

type Role struct {
	RoleID      int `gorm:"primaryKey;autoIncrement:false;"`
	Description string
}
