package model

import (
	"time"

	"gorm.io/gorm"
)

type Hierarchy struct {
	gorm.Model
	Neo4jBaseNode    `bson:"-"`
	AgentID          int
	Agent            *Person `gogm:"startNode bson:"-""`
	SupervisorID     *int
	Supervisor       *Person   `gogm:"endNode" bson:"-"`
	ModificationDate time.Time `gogm:"name=modification_date"`
	AgentStatus      *Status   `gogm:"name=agent_status"`
}

type Status int

const (
	inactive Status = 0
	active   Status = 1
)
