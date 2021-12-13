package main

import (
	// "git.sys-tem.org/caos/db4bigdata/internal/importer"

	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"git.sys-tem.org/caos/db4bigdata/internal/util"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/sirupsen/logrus"
)

// This package should run to find persons
func main() {
	util.SetupLogs()
	logrus.Println("Start to find persons from MongoDB")
	mongo, err := model.ConnectStorage(model.MongoDB)
	if err != nil {
		logrus.Fatal("Connect to MongoDB failed: ", err)
	}

	var result_persons []model.Person

	var qry_persons = bson.D{{"name", "Schott"}}
	mongo.Find(&qry_persons, &result_persons)

	for i, s := range result_persons {
		logrus.Info(i, s)
	}

	var result_invoices []model.Invoice
	var qry_invoices = bson.D{{"invoiceid", "990"}}
	mongo.Find(&qry_invoices, &result_invoices)

	for i, s := range result_invoices {
		logrus.Info(i, s)
	}
}
