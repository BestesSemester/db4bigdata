package main

import (
	// "git.sys-tem.org/caos/db4bigdata/internal/importer"

	"errors"
	"os"
	"strconv"
	"time"

	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"git.sys-tem.org/caos/db4bigdata/internal/util"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sirupsen/logrus"
)

// This package should run to find persons
func main() {
	// AgentID, StartPayDate, EndPayDate
	util.SetupLogs()
	argsWithoutProg := os.Args[1:]

	agentId, _ := strconv.Atoi(argsWithoutProg[0])
	startDate, _ := time.Parse("2006-01-02", argsWithoutProg[1])
	endDate, _ := time.Parse("2006-01-02", argsWithoutProg[2])

	mongo, err := model.ConnectStorage(model.MongoDB)
	if err != nil {
		logrus.Fatal("Connect to MongoDB failed: ", err)
	}

	var all_invoices []model.Invoice
	provision_map := make(map[uint]float32)
	supervisors_map := make(map[uint][]int)

	logrus.Info("Start to calculate provision for all agents")

	// pm := performancemeasurement.New(model.MongoDB, "horrorlog")
	// pm.Start("test", 1*time.Second)
	startTime := time.Now()
	var invoice_sum = float32(0)

	downline_ids := findAllDownlineAgents(mongo, agentId)
	// mongo.Find(bson.D{{"agent.personid", bson.D{{"$in", downline_ids}}}}, &all_invoices)
	mongo.Find(bson.D{
		{"invoicedate", bson.D{{"$gte", startDate}, {"$lt", endDate}}},
		{"agent.personid", bson.D{{"$in", downline_ids}}},
	}, &all_invoices)

	for _, invoice := range all_invoices {
		invoice_sum = invoice_sum + invoice.NetSum
		var invoice_agentID = invoice.Agent.PersonID
		var invoice_provision = invoice.NetSum * 0.1

		// If we don´t know supervisors, request them from db
		if _, found := supervisors_map[uint(invoice_agentID)]; !found {
			supervisors_map[uint(invoice_agentID)] = findAllSupervisorsByAgentPersonId(mongo, invoice_agentID)
		}
		var supervisorIds = supervisors_map[uint(invoice_agentID)]

		if len(supervisorIds) > 0 { // Agent has supervisors
			// 70% of provison for the agent
			var provision_for_agent = invoice_provision * 0.7
			addProvisionToProvisionMap(provision_map, invoice_agentID, provision_for_agent)

			// 30% of provision shared for all supervisors
			var provision_for_others = (invoice_provision - provision_for_agent) / float32(len(supervisorIds))
			for _, supervisorID := range supervisorIds {
				addProvisionToProvisionMap(provision_map, supervisorID, provision_for_others)
			}
		} else {
			// No supervisor whole provision for the agent
			addProvisionToProvisionMap(provision_map, invoice_agentID, invoice_provision)
		}

	}
	// pm.Stop()
	// pm.Run()
	elapsed := time.Since(startTime)

	logrus.Info("Finished to calculate provision in ", elapsed)
	logrus.Info("Provsion for agent ", agentId, " is ", provision_map[uint(agentId)])
}

func addProvisionToProvisionMap(m map[uint]float32, id int, provision float32) {
	if v, found := m[uint(id)]; found {
		m[uint(id)] = v + provision
	} else {
		m[uint(id)] = provision
	}
}

// !!!! Attention - Recursive function !!!!
func findAllSupervisorsByAgentPersonId(mongo model.Database, personID int) []int {
	var ret []int
	var supervisorId, err = findSupervisorIDByAgentPersonId(mongo, personID)
	if err != nil {
		logrus.Error(err)
		return ret
	}

	if supervisorId == 0 || supervisorId == -1 { // If supervisorID is 0 this is the big boss
		// return append(ret, supervisorId)
		return ret
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
		if agent_hierarchy[0].Supervisor == nil {
			return -1, nil
		} else {
			return agent_hierarchy[0].Supervisor.PersonID, nil
		}
	}
}

// !!!! Attention - Recursive function !!!!
func findAllDownlineAgents(mongo model.Database, personID int) []int {
	ret := []int{personID}
	var agentIds, err = findAgentsBySupervisorId(mongo, personID)
	if err != nil {
		logrus.Error(err)
		return ret
	}

	if len(agentIds) == 0 {
		return ret
	} else {
		for _, agentId := range agentIds {
			ret = append(findAllDownlineAgents(mongo, agentId), ret...)
		}
		return ret
	}
}

func findAgentsBySupervisorId(mongo model.Database, personID int) ([]int, error) {
	var agent_hierarchy []model.Hierarchy
	var agent_hierarchy_qry = bson.D{{"supervisor.personid", personID}}

	mongo.Find(agent_hierarchy_qry, &agent_hierarchy)

	ret := make([]int, len(agent_hierarchy))
	for i, agent := range agent_hierarchy {
		ret[i] = agent.Agent.PersonID
	}
	return ret, nil
}
