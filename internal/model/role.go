package model

type Role struct {
	Neo4jBaseNode `bson:"-"`
	RoleID        int       `gorm:"primaryKey;autoIncrement:false;" gogm:"name=role_id"`
	Description   string    `gogm:"name=description"`
	People        []*Person `gorm:"-" bson:"-" gogm:"direction=incoming;relationship=hasRole" json:"-"`
}
