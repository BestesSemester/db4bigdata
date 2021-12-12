package model

import (
	"time"

	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	CustomerID       int
	Name             string
	FirstName        string
	Street           string
	HouseNumber      string
	ZipCode          string
	Residence        string
	PhoneNumber      string
	EmailAddress     string
	BirthDate        time.Time
	RegistrationDate time.Time
	RoleID           int
	Role             Role `gorm:"constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
}
