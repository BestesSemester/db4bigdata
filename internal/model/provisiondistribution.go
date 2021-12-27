package model

type ProvisionDistribution struct {
	InvoiceID     int
	Invoice       Invoice `gorm:"-"`
	AgentID       int
	Agent         Person `gorm:"-"`
	ProvisionPart float32
}
