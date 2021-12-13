package model

type Role struct {
	Neo4jBaseNode
	RoleID      int       `gorm:"primaryKey;autoIncrement:false;" gogm:"name=id"`
	Description string    `gogm:"name=description"`
	People      []*Person `gorm:"-" bson:"-" gogm:"direction=incoming;relationship=hasRole" json:"-"`
}
