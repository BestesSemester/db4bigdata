package main

import (
	"encoding/json"

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
	mssql.Migrate(&model.Person{}, &model.Role{}, &model.Provision{}, &model.Invoice{})
	// person := model.Person{PersonID: 23}
	//p := model.Person{Name: "Meier"}

	person := model.Person{PersonID: 54376}

	if err := mssql.Find(&invoice, &invoice); err != nil {
		logrus.Errorln(err)
	}
	ji, err := json.MarshalIndent(&invoice, "", "	")
	if err != nil {
		logrus.Errorln(err)
	}
	logrus.Println(string(ji))

	// if err := mssql.Find(&person, &person); err != nil {
	// 	logrus.Errorln(err)
	// }
	// jp, err := json.MarshalIndent(&person, "", "	")
	// logrus.Println(string(jp))
	// if err != nil {
	// 	logrus.Errorln(err)
	// }

	// hierarchy := &model.Hierarchy{
	// 	AgentID: 2305,
	// }
	// if err := mssql.Find(hierarchy, hierarchy); err != nil {
	// 	logrus.Errorln(err)
	// }
	// jh, err := json.MarshalIndent(hierarchy, "", "	")
	// if err != nil {
	// 	logrus.Errorln(err)
	// }
	// logrus.Println(string(jh))

}
