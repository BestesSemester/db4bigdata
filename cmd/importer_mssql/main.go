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
	mssql, err := model.ConnectStorage(model.MSQL)
	if err != nil {
		logrus.Fatalln(err)
	}
	//
	// p := model.Person{Name: "Scheffel"}

	// Call importer
	// **** Following lines just works in debug mode ****
	people := []model.Person{}
	importer.ImportPersonsFromJSON("./generators/output_data/persons.json", &people)
	// importer.ImportInvoiceFromJSON("./generators/output_data/invoices.json")
	// importer.ImportHierarchyFromJSON("./generators/output_data/hierarchy.json")
	mssql.Migrate(&model.Person{})
	mssql.Save(&people)

	// mssql.Find("", &p)
}
