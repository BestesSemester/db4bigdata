package model

import (
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	Neo4jBaseNode
	InvoiceID   int
	InvoiceDate time.Time
	Customer    *Person
	//Name          string
	//FirstName     string
	//StreetHouseno string
	//ZipCode       int
	//Residence     string
	Agent        *Person
	NettoSum     float32
	VAT          float32
	BruttoSum    float32
	Provision    *Provision
	ProvisionSum float32
	PayDate      time.Time
	OpenSum      float32
}
