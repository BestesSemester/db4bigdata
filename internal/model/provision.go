package model

import "time"

type Provision struct {
	ProvisionID      int
	ProvisionAmount  float32
	MainAgentAmount  float32
	ModificationDate time.Time
}
