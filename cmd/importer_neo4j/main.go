package main

import (
	// "git.sys-tem.org/caos/db4bigdata/internal/importer"

	"encoding/json"

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
	people := []*model.Person{}
	invoices := []*model.Invoice{}
	hierarchy := []*model.Hierarchy{}
	importer.ImportInterfaceFromJSON("./generators/output_data/persons.json", &people)
	importer.ImportInterfaceFromJSON("./generators/output_data/hierarchy.json", &hierarchy)
	importer.ImportInterfaceFromJSON("./generators/output_data/invoices.json", &invoices)
	model.InterconnectPersonRoles(&people)
	// hpeople := model.MatchHirarchy(&people, &hierarchy)

	model.MatchPeopleAndInvoices(&people, &invoices)
	str, err := json.MarshalIndent(&invoices, "", "	")
	if err != nil {
		logrus.Println(err)
	}
	logrus.Println(string(str))
	neo4j.Save(&invoices)
	// if err := neo4j.Save(&hpeople); err != nil {
	// 	logrus.Errorln(err)
	// }
	// mssql.Find("", &p)
}
