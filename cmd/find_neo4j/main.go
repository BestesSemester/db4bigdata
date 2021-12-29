package main

import (
	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"git.sys-tem.org/caos/db4bigdata/internal/util"

	"github.com/sirupsen/logrus"
)

// This package should run to generate and import new data into the database system
func main() {
	util.SetupLogs()
	logrus.Println("hello")
	neo4j, err := model.ConnectStorage(model.Neo4J)
	if err != nil {
		logrus.Fatalln(err)
	}

	p := model.Person{PersonID: 22}
	//p := model.Person{Name: "Meier"}

	p_target := model.Person{}

	neo4j.Find(&p, &p_target)

}
