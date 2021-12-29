package main

import (
	// "git.sys-tem.org/caos/db4bigdata/internal/importer"

	"git.sys-tem.org/caos/db4bigdata/internal/importer"
	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"git.sys-tem.org/caos/db4bigdata/internal/util"

	"github.com/sirupsen/logrus"
)

// This package should run to generate and import new data into the database system
func main() {
	util.SetupLogs()
	logrus.Println("hello")
	neo4j, err := model.ConnectStorage(model.Neo4J)
	if err != nil {
		logrus.Fatalln(err)
	}
	//
	// p := model.Person{Name: "Scheffel"}

	// Call importer
	// **** Following lines just works in debug mode ****
	people := []model.Person{}
	invoices := []model.Invoice{}
	hierarchy := []model.Hierarchy{}
	importer.ImportPersonsFromJSON("./generators/output_data/persons.json", &people)
	importer.ImportHierarchyFromJSON("./generators/output_data/hierarchy.json", &hierarchy)
	model.InterconnectPersonRoles(&people)
	hpeople := model.MatchHirarchy(people, hierarchy)

	importer.ImportInvoiceFromJSON("./generators/output_data/invoices.json", &invoices)
	ipeople, invoices := model.MatchPeopleAndInvoices(hpeople, invoices)
	neo4j.Save(&invoices)
	if err := neo4j.Save(&ipeople); err != nil {
		logrus.Errorln(err)
	}
	// mssql.Find("", &p)
}
