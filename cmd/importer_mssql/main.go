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

	people := []model.Person{}
	importer.ImportPersonsFromJSON("./generators/output_data/persons.json", &people)
	mssql.Migrate(&model.Person{})
	mssql.Save(&people)

	hierarchies := []model.Hierarchy{}
	importer.ImportHierarchyFromJSON("./generators/output_data/hierarchy.json", &hierarchies)
	mssql.Migrate(&model.Hierarchy{})
	mssql.Save(&hierarchies)

	invoices := []model.Invoice{}
	importer.ImportInvoiceFromJSON("./generators/output_data/invoices.json", &invoices)
	mssql.Migrate(&model.Invoice{})
	mssql.Save(&invoices)

}
