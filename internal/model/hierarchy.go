package model

import (
	"time"

	"gorm.io/gorm"
)

type Hierarchy struct {
	gorm.Model
	Neo4jBaseNode    `bson:"-"`
	AgentID          int
	Agent            *Person
	SupervisorID     *int
	Supervisor       *Person
	ModificationDate time.Time
	AgentStatus      *Status
}

type Status int

const (
	inactive Status = 0
	active   Status = 1
)
