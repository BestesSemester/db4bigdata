package model

import (
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model    `bson:"-"`
	Neo4jBaseNode `bson:"-"`
	InvoiceID     int       `gogm:"name=invoice_id"`
	InvoiceDate   time.Time `gogm:"name=invoice_date"`
	Customer      *Person   `gogm:"direction=incoming;relationship=bought"`
	//Name          string
	//FirstName     string
	//StreetHouseno string
	//ZipCode       int
	//Residence     string
	Agent        *Person    `gogm:"direction=incoming;relationship=sold"`
	NettoSum     float32    `gogm:"name=netto_sum"`
	VAT          float32    `gogm:"name=VAT"`
	BruttoSum    float32    `gogm:"name=brutto_sum"`
	Provision    *Provision `gogm:"direction=outgoing;relationship=contains"`
	ProvisionSum float32    `gogm:"name=provision_sum"`
	PayDate      time.Time  `gogm:"name=pay_date"`
	OpenSum      float32    `gogm:"name=open_sum"`
}
