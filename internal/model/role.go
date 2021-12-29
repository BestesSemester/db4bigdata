package model

type Role struct {
	Neo4jBaseNode `bson:"-"`
	RoleID        *int64    `gorm:"primaryKey;autoIncrement:false;" gogm:"name=role_id;pk=default"`
	Description   string    `gogm:"name=description"`
	People        []*Person `gorm:"-" bson:"-" gogm:"direction=incoming;relationship=hasRole" json:"-"`
}
