package model

import (
	"time"
)

type Hierarchy struct {
	Neo4jBaseNode
	Agent            *Person
	Supervisor       *Person
	ModificationDate time.Time
	AgentStatus      *Status
}

type Status int

const (
	inactive Status = 0
	active   Status = 1
)
