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
	//hierarchy := []model.Hierarchy{}
	importer.ImportPersonsFromJSON("./generators/output_data/persons.json", &people)
	model.InterconnectPersonRoles(&people)
	importer.ImportInvoiceFromJSON("./generators/output_data/invoices.json", &invoices)
	people, invoices = model.MatchPeopleAndInvoices(people, invoices)
	//importer.ImportHierarchyFromJSON("./generators/output_data/hierarchy.json", &hierarchy)
	neo4j.Save(&invoices)
	//neo4j.Save(&hierarchy)

	// mssql.Find("", &p)
}
