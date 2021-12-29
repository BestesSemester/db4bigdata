package main

import (
	// "git.sys-tem.org/caos/db4bigdata/internal/importer"
	"time"

	"git.sys-tem.org/caos/db4bigdata/internal/db"
	"git.sys-tem.org/caos/db4bigdata/internal/importer"
	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"git.sys-tem.org/caos/db4bigdata/internal/util"

	"github.com/sirupsen/logrus"
)

// This package should run to generate and import new data into the database systems
func main() {
	util.SetupLogs()

	logrus.Info("Start to import data ")
	startTime := time.Now()
	mongo, err := db.ConnectStorage(db.MongoDB)
	if err != nil {
		logrus.Fatal("Import to MongoDB failed: ", err)
	}

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
	elapsed := time.Since(startTime)
	logrus.Info("Finished to import data in ", elapsed)
}
