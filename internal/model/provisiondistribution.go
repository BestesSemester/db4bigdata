package model

type ProvisionDistribution struct {
	InvoiceID     int
	Invoice       *Invoice
	AgentID       int
	Agent         *Person
	ProvisionPart float32
}
