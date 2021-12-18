package model

import (
	"time"
)

type Hierarchy struct {
	Neo4jBaseNode    `bson:"-"`
	Agent            *Person `gorm:"constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
	Supervisor       *Person `gorm:"constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
	ModificationDate time.Time
	AgentStatus      *Status
}

type Status int

const (
	inactive Status = 0
	active   Status = 1
)
