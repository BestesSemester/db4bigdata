package model

import (
	"time"

	"gorm.io/gorm"
)

type Person struct {
	gorm.Model       `bson:"-"`
	Neo4jBaseNode    `bson:"-"`
	PersonID         int       `gorm:"primaryKey;" gogm:"name=person_id"`
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
	Role             *Role      `gorm:"constraint:OnUpdate:CASCADE;OnDelete:SET NULL;" gogm:"direction=outgoing;relationship=hasRole"`
	Supervisor       *Person    `gorm:"-" bson:"-" gogm:"direction=outgoing;relationship=supervised_by"`
	AgentInvoices    []*Invoice `gorm:"-" bson:"-" gogm:"direction=outgoing;relationship=sold"`
	CustomerInvoices []*Invoice `gorm:"-" bson:"-" gogm:"direction=outgoing;relationship=bought"`
	Employees        []*Person  `gorm:"-" bson:"-" gogm:"direction=incoming;relationship=supervised_by"`
}

func InterconnectPersonRoles(pe *[]Person) {
	roles := make(map[int]*Role)
	people := *pe
	for i := range people {
		roleid := people[i].Role.RoleID
		if roles[roleid] == nil {
			roles[roleid] = people[i].Role
		} else {
			people[i].Role = roles[roleid]
		}
		role := roles[roleid]
		role.People = append(role.People, &people[i])
	}
	pe = &people
}

func MatchPeopleAndInvoices(people []Person, in []Invoice) ([]Person, []Invoice) {
	p := make(map[int]*Person)
	for k := range people {
		p[people[k].PersonID] = &people[k]
	}
	for i, invoice := range in {
		save_invoice := invoice
		agent := p[invoice.Agent.PersonID]
		customer := p[invoice.Customer.PersonID]
		in[i].Agent = agent
		in[i].Customer = customer
		customer.CustomerInvoices = append(p[invoice.Customer.PersonID].CustomerInvoices, &save_invoice)
		p[invoice.Customer.PersonID] = customer
		agent.AgentInvoices = append(p[invoice.Agent.PersonID].AgentInvoices, &save_invoice)
		p[invoice.Agent.PersonID] = agent
	}
	outPeople := []Person{}
	for k := range p {
		outPeople = append(outPeople, *p[k])
	}
	return outPeople, in
}

func MatchHirarchy(people []Person, hierarchy []Hierarchy) []Person {
	p := make(map[int]*Person)
	for k := range people {
		p[people[k].PersonID] = &people[k]
	}
	for _, hi := range hierarchy {
		agentID := hi.Agent.PersonID

		if hi.Supervisor != nil {
			supervisor := p[hi.Supervisor.PersonID]
			agent := p[agentID]
			agent.Supervisor = supervisor
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
