package model

import (
	"time"

	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	Neo4jBaseNode
	CustomerID       int       `gogm:"name=customer_id"`
	Name             string    `gogm:"name=name"`
	FirstName        string    `gogm:"name=first_name"`
	Street           string    `gogm:"name=street"`
	HouseNumber      string    `gogm:"name=house_number"`
	ZipCode          string    `gogm:"name=zip_code"`
	Residence        string    `gogm:"name=residence"`
	PhoneNumber      string    `gogm:"name=phone_number"`
	EmailAddress     string    `gogm:"name=email_address"`
	BirthDate        time.Time `gogm:"name=birth_date"`
	RegistrationDate time.Time `gogm:"name=registration_date"`
	RoleID           int
	Role             *Role `gorm:"constraint:OnUpdate:CASCADE;OnDelete:SET NULL;" gogm:"direction=outgoing;relationship=hasRole"`
}

func InterconnectPersonRoles(people []Person) []Person {
	roles := make(map[int]*Role)
	for i := range people {
		roleid := people[i].Role.RoleID
		if roles[roleid] == nil {
			roles[roleid] = people[i].Role
		} else {
			people[i].Role = roles[roleid]
		}
		people[i].Role.People = append(people[i].Role.People, &people[i])
	}
	return people
}
