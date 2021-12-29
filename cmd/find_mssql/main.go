package main

import (
	"encoding/json"

	"git.sys-tem.org/caos/db4bigdata/internal/db"
	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"git.sys-tem.org/caos/db4bigdata/internal/util"

	"github.com/sirupsen/logrus"
)

// This package should run to generate and import new data into the database system
func main() {
	util.SetupLogs()
	logrus.Println("hello")
	mssql, err := db.ConnectStorage(db.MSQL)
	if err != nil {
		logrus.Fatalln(err)
	}
	mssql.Migrate(&model.Person{}, &model.Role{}, &model.Provision{}, &model.Invoice{})
	// person := model.Person{PersonID: 23}
	//p := model.Person{Name: "Meier"}

	invoice := model.Invoice{InvoiceID: 848}

	if err := mssql.Find(&invoice, &invoice); err != nil {
		logrus.Errorln(err)
	}
	if err := getPeopleHierarchyForInvoice(mssql, &invoice); err != nil {
		logrus.Errorln(err)
	}
	ji, err := json.MarshalIndent(&invoice, "", "	")
	if err != nil {
		logrus.Errorln(err)
	}

	logrus.Println(string(ji))

	person := model.Person{PersonID: 54376}

	if err := mssql.Find(&person, &person); err != nil {
		logrus.Errorln(err)
	}
	if err := getPeopleHierarchyForInvoice(mssql, &invoice); err != nil {
		logrus.Errorln(err)
	}
	jp, err := json.MarshalIndent(&person, "", "	")
	if err != nil {
		logrus.Errorln(err)
	}
	logrus.Println(string(jp))

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

func getPeopleHierarchyForInvoice(mssql db.Database, invoice *model.Invoice) error {
	if err := mssql.Find(invoice, invoice); err != nil {
		logrus.Errorln(err)
		return err
	}
	agent := &model.Person{PersonID: invoice.AgentID}
	if err := mssql.Find(agent, agent); err != nil {
		logrus.Errorln(err)
		return err
	}
	invoice.Agent = agent
	ji, err := json.MarshalIndent(agent, "", "	")
	if err != nil {
		logrus.Errorln(err)
	}
	logrus.Println(string(ji))
	return nil
}
