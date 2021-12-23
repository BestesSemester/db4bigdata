package model

import (
	"time"

	"gorm.io/gorm"
)

type Hierarchy struct {
	gorm.Model       `bson:"-"`
	Neo4jBaseNode    `bson:"-"`
	Agent            *Person   `gorm:"foreignKey:PersonID" gogm:"startNode"`
	Supervisor       *Person   `gorm:"foreignKey:PersonID" gogm:"endNode"`
	ModificationDate time.Time `gogm:"name=modification_date"`
	AgentStatus      *Status   `gogm:"name=agent_status"`
}

type Status int

const (
	inactive Status = 0
	active   Status = 1
)
