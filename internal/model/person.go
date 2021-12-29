package model

import (
	"time"
)

type Person struct {
	Neo4jBaseNode    `bson:"-"`
	PersonID         *int64    `gorm:"primaryKey" gogm:"name=person_id"`
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
	SupervisorID     *int64     `gogm:"-" bson:"-"`
	Supervisor       *Person    `gogm:"direction=outgoing;relationship=supervised_by" bson:"-"`
	AgentInvoices    []*Invoice `gorm:"-" bson:"-" gogm:"direction=outgoing;relationship=sold"`
	CustomerInvoices []*Invoice `gorm:"-" bson:"-" gogm:"direction=outgoing;relationship=bought"`
	Employees        []*Person  `gorm:"-" bson:"-" gogm:"direction=incoming;relationship=supervised_by"`
}

func InterconnectPersonRoles(people *[]*Person) {
	roles := make(map[int64]*Role)
	for _, person := range *people {
		roleid := *person.Role.RoleID
		if roles[roleid] == nil {
			roles[roleid] = person.Role
		} else {
			person.Role = roles[roleid]
		}
		roles[roleid].People = append(roles[roleid].People, person)
	}
}

func MatchPeopleAndInvoices(people *[]*Person, invoices *[]*Invoice) {
	p := make(map[int64]*Person)
	for k := range *people {
		p[*(*people)[k].PersonID] = (*people)[k]
	}
	for _, invoice := range *invoices {
		save_invoice := invoice
		// agent := p[*invoice.Agent.PersonID]
		// customer := p[*invoice.Customer.PersonID]
		// (*invoices)[i].Agent = agent
		// (*invoices)[i].Customer = customer
		// customer.CustomerInvoices = append(p[*invoice.Customer.PersonID].CustomerInvoices, &save_invoice)
		// p[*invoice.Customer.PersonID] = customer
		// agent.AgentInvoices = append(p[*invoice.Agent.PersonID].AgentInvoices, &save_invoice)
		// p[*invoice.Agent.PersonID] = agent
		for _, person := range *people {
			if person.PersonID == invoice.Agent.PersonID {
				person.AgentInvoices = append(person.AgentInvoices, save_invoice)
				invoice.Agent = person
			}
			if person.PersonID == invoice.Customer.PersonID {
				person.CustomerInvoices = append(person.CustomerInvoices, save_invoice)
				invoice.Customer = person
			}
		}
	}
}

func MatchHirarchy(people *[]*Person, hierarchy *[]*Hierarchy) []Person {
	p := make(map[int64]*Person)
	for k, per := range *people {
		pe := *people
		p[*pe[k].PersonID] = per
	}
	for _, hi := range *hierarchy {
		agentID := hi.Agent.PersonID

		if hi.Supervisor != nil {
			supervisor := p[*hi.Supervisor.PersonID]
			agent := p[*agentID]
			agent.Supervisor = supervisor
			agent.SupervisorID = supervisor.PersonID
			supervisor.Employees = append(supervisor.Employees, agent)
		}
		// if agentID == 1078 {
		// 	logrus.Println(p[agentID])
		// }
	}
	outPeople := []Person{}
	for k := range p {
		outPeople = append(outPeople, *p[k])
	}
	return outPeople
}
