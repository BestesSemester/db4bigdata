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

	mssql.Migrate(&model.Provision{}, &model.Invoice{}, &model.ProvisionDistribution{})

	i := model.ProvisionDistribution{AgentID: 1080, InvoiceID: 114}

	i_t := model.ProvisionDistribution{}

	if err := mssql.Find(&i, &i_t); err != nil {
		logrus.Errorln(err)
	}
	ji, err := json.MarshalIndent(&i_t, "", "	")
	if err != nil {
		logrus.Errorln(err)
	}
	logrus.Println(string(ji))

}
