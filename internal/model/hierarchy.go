package model

import (
	"time"
)

type Hierarchy struct {
	Neo4jBaseNode    `bson:"-"`
	Agent            *Person `gorm:"-" bson:"-" gogm:"direction=incoming;relationship=hasRole" json:"-"`
	Supervisor       *Person `gorm:"-" bson:"-" gogm:"direction=incoming;relationship=hasRole" json:"-"`
	ModificationDate time.Time
	AgentStatus      *Status
}

type Status int

const (
	inactive Status = 0
	active   Status = 1
)
