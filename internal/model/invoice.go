package model

import (
	"time"
)

type Invoice struct {
	Neo4jBaseNode `bson:"-"`
	InvoiceID     int       `gorm:"primaryKey" gogm:"name=invoice_id"`
	InvoiceDate   time.Time `gogm:"name=invoice_date"`
	CustomerID    int
	Customer      *Person `gogm:"direction=incoming;relationship=bought"`
	AgentID       int
	Agent         *Person `gogm:"direction=incoming;relationship=sold"`
	NetSum        float32 `gogm:"name=netto_sum"`
	VAT           float32 `gogm:"name=VAT"`
	GrossSum      float32 `gogm:"name=brutto_sum"`
	ProvisionID   *int
	Provision     *Provision //`gogm:"direction=outgoing;relationship=contains"`
	ProvisionSum  float32    `gogm:"name=provision_sum"`
	PayDate       time.Time  `gogm:"name=pay_date"`
	OpenSum       float32    `gogm:"name=open_sum"`
}
