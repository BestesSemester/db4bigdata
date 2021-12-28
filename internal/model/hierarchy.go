package model

import (
	"time"

	"gorm.io/gorm"
)

type Hierarchy struct {
	gorm.Model
	Neo4jBaseNode    `bson:"-"`
	AgentID          int
	Agent            *Person `gorm:"foreignKey:PersonID" gogm:"startNode"`
	SupervisorID     *int
	Supervisor       *Person   `gorm:"foreignKey:PersonID" gogm:"endNode"`
	ModificationDate time.Time `gogm:"name=modification_date"`
	AgentStatus      *Status   `gogm:"name=agent_status"`
}

type Status int

const (
	inactive Status = 0
	active   Status = 1
)
