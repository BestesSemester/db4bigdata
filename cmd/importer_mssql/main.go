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

	mssql.Migrate(&model.Person{}, &model.Role{}, &model.Hierarchy{}, &model.Invoice{}, &model.Provision{})
	people := []model.Person{}
	importer.ImportPersonsFromJSON("./generators/output_data/persons.json", &people)
	invoices := []model.Invoice{}
	importer.ImportInvoiceFromJSON("./generators/output_data/invoices.json", &invoices)
	mssql.Save(&invoices)

	hierarchies := []model.Hierarchy{}
	importer.ImportHierarchyFromJSON("./generators/output_data/hierarchy.json", &hierarchies)
	people = model.MatchHirarchy(people, hierarchies)
	mssql.Save(&people)
	mssql.Save(&hierarchies)

}
