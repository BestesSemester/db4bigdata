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
	Customer      *Person `gogm:"direction=incoming;relationship=bought"`
	//Name          string
	//FirstName     string
	//StreetHouseno string
	//ZipCode       int
	//Residence     string
	Agent        *Person `gogm:"direction=incoming;relationship=sold"`
	NettoSum     float32
	VAT          float32
	BruttoSum    float32
	Provision    *Provision
	ProvisionSum float32
	PayDate      time.Time
	OpenSum      float32
}
