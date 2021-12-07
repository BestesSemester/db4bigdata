package model

import "time"

type Person struct {
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
	Role             Role
}
