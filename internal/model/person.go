package model

import (
	"time"
)

type Person struct {
	Neo4jBaseNode    `bson:"-"`
	PersonID         int       `gorm:"primaryKey" gogm:"name=person_id"`
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
	Role             *Role      `gogm:"direction=outgoing;relationship=hasRole"`
	SupervisorID     int        `gogm:"-" bson:"-"`
	Supervisor       *Person    `gogm:"direction=outgoing;relationship=supervised_by" bson:"-"`
	AgentInvoices    []*Invoice `gorm:"-" bson:"-" gogm:"direction=outgoing;relationship=sold"`
	CustomerInvoices []*Invoice `gorm:"-" bson:"-" gogm:"direction=outgoing;relationship=bought"`
	Employees        []*Person  `gorm:"-" bson:"-" gogm:"direction=incoming;relationship=supervised_by"`
}

func InterconnectPersonRoles(people *[]*Person) {
	roles := make(map[int]*Role)
	for _, person := range *people {
		roleid := person.Role.RoleID
		if roles[roleid] == nil {
			roles[roleid] = person.Role
		} else {
			person.Role = roles[roleid]
		}
		roles[roleid].People = append(roles[roleid].People, person)
	}

}

func MatchPeopleAndInvoices(people *[]*Person, invoices *[]*Invoice) {
	p := make(map[int]*Person)
	for k := range *people {
		p[(*people)[k].PersonID] = (*people)[k]
	}
	for _, invoice := range *invoices {

		for _, person := range *people {

			if person.PersonID == invoice.Agent.PersonID {
				person.AgentInvoices = append(person.AgentInvoices, invoice)
				invoice.Agent = person
			}
			if person.PersonID == invoice.Customer.PersonID {
				person.CustomerInvoices = append(person.CustomerInvoices, invoice)
				invoice.Customer = person
			}
		}
	}
}

func MatchHirarchy(people *[]*Person, hierarchy *[]*Hierarchy) {
	p := make(map[int]*Person)
	for k, per := range *people {
		pe := *people
		p[pe[k].PersonID] = per
	}
	for _, set := range *hierarchy {
		if set.Supervisor != nil {
			for _, ag := range *people {
				if ag.PersonID == set.Agent.PersonID {
					for _, sup := range *people {
						if sup.PersonID == set.Supervisor.PersonID {
							sup.Employees = append(sup.Employees, ag)
							ag.Supervisor = sup
							ag.SupervisorID = sup.PersonID
						}
					}
				}
			}
		}
	}
}
