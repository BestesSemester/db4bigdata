package main

import (
	"git.sys-tem.org/caos/db4bigdata/internal/importer"
	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"git.sys-tem.org/caos/db4bigdata/internal/util"

	"github.com/sirupsen/logrus"
)

// This package should run to generate and import new data into the database system
func main() {
	util.SetupLogs()
	logrus.Println("hello")
	mssql, err := model.ConnectStorage(model.MSQL)
	if err != nil {
		logrus.Fatalln(err)
	}

	// Call importer
	// **** Following lines just works in debug mode ****

	invoices := []*model.Invoice{}
	people := []*model.Person{}
	hierarchy := []*model.Hierarchy{}

	mssql.Migrate(&model.Person{}, &model.Role{}, &model.Hierarchy{}, &model.Invoice{}, &model.Provision{})

	importer.ImportInterfaceFromJSON("./generators/output_data/invoices.json", &invoices)
	importer.ImportInterfaceFromJSON("./generators/output_data/persons.json", &people)
	importer.ImportInterfaceFromJSON("./generators/output_data/hierarchy.json", &hierarchy)

	mssql.Save(&people)

	model.InterconnectPersonRoles(&people)
	model.MatchHirarchy(&people, &hierarchy)
	model.MatchPeopleAndInvoices(&people, &invoices)

	mssql.Save(&invoices)
	mssql.Save(&hierarchy)
}
