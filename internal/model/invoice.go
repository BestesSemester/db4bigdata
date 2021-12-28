package model

import (
	"time"
)

type Invoice struct {
	Neo4jBaseNode `bson:"-"`
	InvoiceID     int       `gogm:"name=invoice_id"`
	InvoiceDate   time.Time `gogm:"name=invoice_date"`
	CustomerID    int
	Customer      *Person `gorm:"foreignKey:PersonID" gogm:"direction=incoming;relationship=bought"`
	AgentID       int
	Agent         *Person    `gorm:"foreignKey:PersonID" gogm:"direction=incoming;relationship=sold"`
	NetSum        float32    `gogm:"name=netto_sum"`
	VAT           float32    `gogm:"name=VAT"`
	GrossSum      float32    `gogm:"name=brutto_sum"`
	Provision     *Provision `gorm:"foreignKey:ProvisionID"` //`gogm:"direction=outgoing;relationship=contains"`
	ProvisionSum  float32    `gogm:"name=provision_sum"`
	PayDate       time.Time  `gogm:"name=pay_date"`
	OpenSum       float32    `gogm:"name=open_sum"`
}
