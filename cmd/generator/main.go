package main

import (
	"git.sys-tem.org/caos/db4bigdata/internal/importer"
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

	importer.ImportPersonsFromJSON("./generators/output_data/persons.json")

	// mongo.Save()

	// mssql.Find("", &p)
}
