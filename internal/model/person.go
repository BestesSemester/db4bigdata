package model

import "time"

type Person struct {
	CustomerID       int       `bson: "customerId"`
	Name             string    `bson: "name"`
	FirstName        string    `bson: "firstname"`
	Street           string    `bson: "street"`
	HouseNumber      string    `bson: "housenumber"`
	ZipCode          string    `bson: "zipcode"`
	Residence        string    `bson: "residence"`
	PhoneNumber      string    `bson: "phonenumber"`
	EmailAddress     string    `bson: "emailaddress"`
	BirthDate        time.Time `bson: "birthdate"`
	RegistrationDate time.Time `bson: "registrationdate"`
	Role             Role      `bson: "role"`
}
