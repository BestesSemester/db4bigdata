package model

import (
	"time"
)

type Provision struct {
	Neo4jBaseNode
	ProvisionID      int
	ProvisionAmount  float32
	MainAgentAmount  float32
	ModificationDate time.Time
}
