package main

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"git.sys-tem.org/caos/db4bigdata/internal/db"
	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"git.sys-tem.org/caos/db4bigdata/internal/performancemeasurement"
	"git.sys-tem.org/caos/db4bigdata/internal/util"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

// This package should run to generate and import new data into the database system
func main() {
	util.SetupLogs()
	util.SetupLogs()
	argsWithoutProg := os.Args[1:]

	agentId, _ := strconv.Atoi(argsWithoutProg[0])
	startDate, _ := time.Parse("2006-01-02", argsWithoutProg[1])
	endDate, _ := time.Parse("2006-01-02", argsWithoutProg[2])
	startTime := time.Now()

	pm := performancemeasurement.New(db.MSQL, "mssql_"+argsWithoutProg[0]+"_"+argsWithoutProg[1]+"_"+argsWithoutProg[2])
	pm.Start("MSSQL calculate performance", 1*time.Second)

	//Connection to server

	dsn := "sqlserver://sa:changeMe1234@localhost:1433?database=master"
	mssqldb, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	util.SetupLogs()
	logrus.Println("hello")
	mssql, err := db.ConnectStorage(db.MSQL)
	if err != nil {
		logrus.Fatalln(err)
	}
	mssql.Migrate(&model.Person{}, &model.Role{}, &model.Provision{}, &model.Invoice{})

	//Calculate the Provisions and write into table
	//write all agents in invoices into string
	var agents []string
	invoice_t := model.Invoice{}
	mssqldb.Model(&invoice_t).Where("Invoice_Date between ? and ?", startDate, endDate).Distinct().Pluck("Agent_ID", &agents)
	prov_delete := &model.ProvisionDistribution{}
	mssqldb.Where("provision_part > 0").Delete(&prov_delete)
	//write all provisions into table
	for _, agent := range agents {
		mssqldb.Exec(`WITH
			temp2 as(
			select agent_id , supervisor_id
			from hierarchies
			where agent_id = ?
			union all
			select a.agent_id, a.supervisor_id
			from hierarchies a inner join temp2 on temp2.supervisor_id = a.agent_id
					)
			insert into provision_distributions
				select i.[Invoice_ID], t.agent_id,
				case
					when (t.agent_id = ? and (select count(*) from temp2) > 1 )
						then i.net_Sum * 0.7 * 0.1
					when (t.agent_id = 1079
						and (select count(*) from temp2) = 1) then i.net_Sum * 0.1
					else i.net_Sum *0.1 * 0.3/((select count(*) from temp2)-1)
				end provision
				from temp2 t, [dbo].[invoices] i
				where i.agent_id = ?
				and invoice_date between ? and ?
				and invoice_id not in (select invoice_id from provision_distributions)
				order by Invoice_ID`, agent, agent, agent, startDate, endDate)
	}

	distr := model.ProvisionDistribution{AgentID: agentId}
	if err := mssql.Find(&distr, &distr); err != nil {
		logrus.Errorln(err)
	}
	jd, err := json.MarshalIndent(&distr, "", "	")
	if err != nil {
		logrus.Errorln(err)
	}
	logrus.Println(string(jd))

	/*
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
	*/

	elapsed := time.Since(startTime)
	pm.MeasureTime("find_mssql", startTime)
	pm.Stop()

	logrus.Info("Finished to calculate provision in ", elapsed)
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
