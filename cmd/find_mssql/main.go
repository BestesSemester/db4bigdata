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
	p := model.Person{PersonID: 23}
	//p := model.Person{Name: "Meier"}

	i := model.Invoice{InvoiceID: 14}
	i_t := model.Invoice{}

	if err := mssql.Find(&i, &i_t); err != nil {
		logrus.Errorln(err)
	}
	ji, err := json.MarshalIndent(&i_t, "", "	")
	if err != nil {
		logrus.Errorln(err)
	}
	logrus.Println(string(ji))
	p_target := model.Person{}

	if err := mssql.Find(&p, &p_target); err != nil {
		logrus.Errorln(err)
	}
	jp, err := json.MarshalIndent(&p_target, "", "	")
	logrus.Println(string(jp))
	if err != nil {
		logrus.Errorln(err)
	}

}
