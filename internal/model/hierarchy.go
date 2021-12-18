package model

import (
	"time"

	"gorm.io/gorm"
)

type Hierarchy struct {
	gorm.Model       `bson:"-"`
	Neo4jBaseNode    `bson:"-"`
	Agent            *Person `gorm:"foreignKey:PersonID"`
	Supervisor       *Person `gorm:"foreignKey:PersonID"`
	ModificationDate time.Time
	AgentStatus      *Status
}

type Status int

const (
	inactive Status = 0
	active   Status = 1
)
