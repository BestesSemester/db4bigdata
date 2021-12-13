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

	var result []model.Person
	var qry = bson.D{{"name", "Schott"}}
	mongo.Find(&qry, &result) //Actually just find by name

	for i, s := range result {
		logrus.Info(i, s)
	}
}
