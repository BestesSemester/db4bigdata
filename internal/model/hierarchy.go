package model

import (
	"time"
)

type Hierarchy struct {
	Neo4jBaseNode    `bson:"-"`
	Agent            *Person   `gogm:"startNode"`
	Supervisor       *Person   `gogm:"endNode"`
	ModificationDate time.Time `gogm:"name=modification_date"`
	AgentStatus      *Status   `gogm:"name=agent_status"`
}

type Status int

const (
	inactive Status = 0
	active   Status = 1
)
