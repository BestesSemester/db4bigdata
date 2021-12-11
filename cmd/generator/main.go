package main

import (
	// "git.sys-tem.org/caos/db4bigdata/internal/importer"

	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"git.sys-tem.org/caos/db4bigdata/internal/util"

	"github.com/sirupsen/logrus"
)

// This package should run to generate and import new data into the database system
func main() {
	util.SetupLogs()
	logrus.Println("hello")
	// mssql, _ := model.ConnectStorage(model.MSQL)
	//
	// p := model.Person{Name: "Scheffel"}
	mongo, err := model.ConnectStorage(model.MongoDB)
	if err != nil {
		logrus.Fatal("Import to MongoDB failed: ", err)
	}
	var persons []model.Person
	// Call importer
	// **** Following lines just works in debug mode ****
	// importer.ImportPersonsFromJSON("./../../generators/output_data/persons.json")
	// importer.ImportInvoiceFromJSON("./../../generators/output_data/invoices.json")
	// importer.ImportHierarchyFromJSON("./../../generators/output_data/hierarchy.json")
	mongo.Find("", &persons)
	logrus.Info(persons)
	// mongo.Save()

	// mssql.Find("", &p)
}
