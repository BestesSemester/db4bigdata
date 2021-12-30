package main

import (
	"os"
	"strconv"
	"time"

	"git.sys-tem.org/caos/db4bigdata/internal/db"
	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"git.sys-tem.org/caos/db4bigdata/internal/performancemeasurement"
	"git.sys-tem.org/caos/db4bigdata/internal/util"

	"github.com/sirupsen/logrus"
)

// This package should run to generate and import new data into the database system
func main() {
	util.SetupLogs()

	// Parse command line arguments
	argsWithoutProg := os.Args[1:]
	agentId, _ := strconv.Atoi(argsWithoutProg[0])
	startdate := argsWithoutProg[1]
	enddate := argsWithoutProg[2]

	logrus.Info("Start to calculate provision for agent ", agentId)

	// Start performance measurement
	pm := performancemeasurement.New(db.Neo4J, "neo4j_"+argsWithoutProg[0]+"_"+argsWithoutProg[1]+"_"+argsWithoutProg[2])
	pm.Start("", 1*time.Second)
	startTime := time.Now()

	neo4j, err := db.ConnectStorage(db.Neo4J)
	if err != nil {
		logrus.Fatalln(err)
	}

	//startdate := "2000-01-01"
	//enddate := "2021-12-30"

	p := model.Person{PersonID: agentId}
	//p := model.Person{Name: "Meier"}

	res := struct {
		Provision float32
		PersonID  int
	}{}

	neoqry := `MATCH (i:Invoice)<-[trigger:sold*0..4]-(agent:Person)-[:supervised_by*0..4]->(supervisor:Person)-[:hasRole]->(role:Role)
				WHERE supervisor.person_id = ` + strconv.Itoa(p.PersonID) + ` and i.pay_date > datetime("` + startdate + `") AND i.pay_date < datetime("` + enddate + `")
				CALL apoc.case([
				supervisor.person_id=agent.person_id AND role.role_id=1,
				'RETURN i.provision as provision',
				supervisor.person_id=agent.person_id AND role.role_id>1,
				'RETURN i.netto_sum*0.1*0.7 as provision'],
				'RETURN i.netto_sum*0.1*0.3/(role.role_id-1) as provision',
				{i:i,agent:agent,supervisor:supervisor, role:role})
				YIELD value
				RETURN {PersonID: supervisor.person_id, Provision: SUM(value.provision)}`

	neo4j.Exec(neoqry, &res)

	// Stop performance measurement
	elapsed := time.Since(startTime)
	pm.MeasureTime("find_neo4j", startTime)
	pm.Stop()

	logrus.Info("Finished to calculate provision in ", elapsed)
	logrus.Println(res)

}
