package model

import (
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model    `bson:"-"`
	Neo4jBaseNode `bson:"-"`
	InvoiceID     int
	InvoiceDate   time.Time
	Customer      *Person `gorm:"foreignKey:PersonID"; gogm:"direction=outgoing;relationship=bought"`
	//Name          string
	//FirstName     string
	//StreetHouseno string
	//ZipCode       int
	//Residence     string
	AgentID      int
	Agent        *Person `gorm:"foreignKey:PersonID"; gogm:"direction=outgoing;relationship=sold"`
	NetSum       float32
	VAT          float32
	GrossSum     float32
	Provision    *Provision `gorm:"foreignKey:ProvisionID"`
	ProvisionSum float32
	PayDate      time.Time
	OpenSum      float32
}
