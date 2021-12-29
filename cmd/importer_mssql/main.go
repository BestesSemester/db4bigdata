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
	invoices := []model.Invoice{}
	importer.ImportInvoiceFromJSON("./generators/output_data/invoices.json", &invoices)

	people := []*model.Person{}
	importer.ImportInterfaceFromJSON("./generators/output_data/persons.json", &people)
	hierarchy := []*model.Hierarchy{}
	importer.ImportInterfaceFromJSON("./generators/output_data/persons.json", &hierarchy)
	model.MatchHirarchy(&people, &hierarchy)
	mssql.Save(&invoices)
	mssql.Save(&hierarchy)
	mssql.Save(&people)
}
