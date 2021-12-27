package model

import (
	"time"
)

type Invoice struct {
	Neo4jBaseNode `bson:"-"`
	InvoiceID     int `gorm:"primaryKey"`
	InvoiceDate   time.Time
	CustomerID    int
	Customer      *Person `gogm:"direction=outgoing;relationship=bought"`
	//Name          string
	//FirstName     string
	//StreetHouseno string
	//ZipCode       int
	//Residence     string
	AgentID      int
	Agent        *Person `gogm:"direction=outgoing;relationship=sold"`
	NetSum       float32
	VAT          float32
	GrossSum     float32
	ProvisionID  int        `bson:"-" gogm:"-"`
	Provision    *Provision `gorm:"-"`
	ProvisionSum float32
	PayDate      time.Time
	OpenSum      float32
}
