package main

import (
	// "git.sys-tem.org/caos/db4bigdata/internal/importer"

	"errors"

	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"git.sys-tem.org/caos/db4bigdata/internal/util"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sirupsen/logrus"
)

// This package should run to find persons
func main() {
	util.SetupLogs()
	logrus.Println("Start to calculate provision for all agents")
	mongo, err := model.ConnectStorage(model.MongoDB)
	if err != nil {
		logrus.Fatal("Connect to MongoDB failed: ", err)
	}

	var all_invoices []model.Invoice
	var result_persons []model.Person
	provision_map := make(map[uint]float32)

	mongo.Find(bson.D{{}}, &all_invoices)
	for _, invoice := range all_invoices {

		var agent = invoice.Agent
		var invoice_provision = invoice.GrossSum * 0.1

		var supervisorIds = findAllSupervisorsByAgentPersonId(mongo, agent.PersonID)
		logrus.Debug("Supervisors for agent: ", agent.PersonID, " < ", supervisorIds, " > ")

		if len(supervisorIds) > 0 { // Agent has supervisors
			// 70% of provison for the agent
			var provision_for_agent = invoice_provision * 0.7
			addProvisionToProvisionMap(provision_map, agent.PersonID, provision_for_agent)

			// 30% of provision shared for all supervisors
			var provision_for_others = (invoice_provision - provision_for_agent) / float32(len(supervisorIds))
			for _, supervisorID := range supervisorIds {
				addProvisionToProvisionMap(provision_map, supervisorID, provision_for_others)
			}
		} else {
			// No supervisor whole provision for the agent
			addProvisionToProvisionMap(provision_map, agent.PersonID, invoice_provision)
		}

	}
	for key, value := range provision_map {
		logrus.Info("Provsion for agent ", key, " is ", value)
	}

	// Vertreter ID suchen und als Ausgabewert die Summe der Provisionen.
	// Rechnung ID eingeben und als Ausgabewert die Provision der Rechnung
	// Für alle Vertreter die Provision berechnen - Das ist die Summe der Provisionen aller Rechnungen.
	// Die Provision einer Rechnung ist 10% -> Vertreter 70% Provision (100% wenn kein Parent). Die restlichen 30% werden gleichmäßig auf alle Parents verteilt
	// Hilfstabelle erzeugen (evtl. in Memory)

	for i, s := range result_persons {
		logrus.Info(i, s)
	}
}

func addProvisionToProvisionMap(m map[uint]float32, id int, provision float32) {
	if v, found := m[uint(id)]; found {
		m[uint(id)] = v + provision
	} else {
		m[uint(id)] = provision
	}
}

func findAllSupervisorsByAgentPersonId(mongo model.Database, personID int) []int {
	var ret []int
	var supervisorId, err = findSupervisorIDByAgentPersonId(mongo, personID)
	if err != nil {
		logrus.Error(err)
		return ret
	}

	if supervisorId == 0 { // If supervisorID is 0 this is the big boss
		return append(ret, supervisorId)
	} else {
		return append(findAllSupervisorsByAgentPersonId(mongo, supervisorId), supervisorId)
	}
}

func findSupervisorIDByAgentPersonId(mongo model.Database, personID int) (int, error) {
	var agent_hierarchy []model.Hierarchy
	var agent_hierarchy_qry = bson.D{{"agent.personid", personID}}

	mongo.Find(agent_hierarchy_qry, &agent_hierarchy)
	if len(agent_hierarchy) < 1 {
		return 0, errors.New("No hierarchy object found")
	} else if len(agent_hierarchy) > 1 {
		return 0, errors.New("Invalid data structure: More than 1 hierarchy object found.")
	} else {
		return agent_hierarchy[0].Supervisor.PersonID, nil
	}
}
