package main

import (
	// "git.sys-tem.org/caos/db4bigdata/internal/importer"
	"git.sys-tem.org/caos/db4bigdata/internal/importer"
	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"git.sys-tem.org/caos/db4bigdata/internal/util"

	"github.com/sirupsen/logrus"
)

// This package should run to generate and import new data into the database systems
func main() {
	util.SetupLogs()
	logrus.Println("hello")
	mongo, err := model.ConnectStorage(model.MongoDB)
	if err != nil {
		logrus.Fatal("Import to MongoDB failed: ", err)
	}
	// mssql, _ := model.ConnectStorage(model.MSQL)
	//
	// p := model.Person{Name: "Scheffel"}

	// Call importer
	persons := []model.Person{}
	importer.ImportPersonsFromJSON("./generators/output_data/persons.json", &persons)
	err = mongo.Save(persons)
	if err != nil {
		logrus.Fatalln(err)
	}

	invoices := []model.Invoice{}
	importer.ImportInvoiceFromJSON("./generators/output_data/invoices.json", &invoices)
	err = mongo.Save(invoices)
	if err != nil {
		logrus.Fatalln(err)
	}
	hierarchy := []model.Hierarchy{}
	importer.ImportHierarchyFromJSON("./generators/output_data/hierarchy.json", &hierarchy)
	err = mongo.Save(hierarchy)
	if err != nil {
		logrus.Fatalln(err)
	}
	// missing save method for msql, will be added if ready
	// missing save method for neo4j, will be added if ready

}
