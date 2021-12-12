//the importer tasks are the following:
// 1. open file
// 2. read bytes
// 3. unmarshal json to object

package importer

import (
	"encoding/json"
	"io/ioutil"

	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"github.com/sirupsen/logrus"
)

func ImportPersonsFromJSON(jsonfile string) {
	// Let's first read the `config.json` file
	content, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		logrus.Fatal("Error when opening file: ", err)
	}
	var persons []model.Person
	err = json.Unmarshal(content, &persons)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Println(len(persons))
	ImportObjectsMongo(persons)
}

func ImportHierarchyFromJSON(jsonfile string) {
	// Let's first read the `config.json` file
	content, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		logrus.Fatal("Error when opening file: ", err)
	}
	var hierarchy []model.Hierarchy
	err = json.Unmarshal(content, &hierarchy)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Println(len(hierarchy))
	ImportObjectsMongo(hierarchy)
}

func ImportInvoiceFromJSON(jsonfile string) {
	// Let's first read the `config.json` file
	content, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		logrus.Fatal("Error when opening file: ", err)
	}
	var invoice []model.Invoice
	err = json.Unmarshal(content, &invoice)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Println(len(invoice))
	ImportObjectsMongo(invoice)
}

// func ImportProvisiondistributionFromJSON(jsonfile string) {
// 	// Let's first read the `config.json` file
// 	content, err := ioutil.ReadFile(jsonfile)
// 	if err != nil {
// 		logrus.Fatal("Error when opening file: ", err)
// 	}
// 	var provdib []model.Provdib
// 	err = json.Unmarshal(content, &provdib)
// 	if err != nil {
// 		logrus.Fatal(err)
// 	}
// 	logrus.Println(len(provdib))
// 	db, err := model.ConnectStorage(model.MSQL)
// 	if err != nil {
// 		logrus.Fatal("Error in saving", err)
// 	}
// 	err = db.Save(provdib)
// 	if err != nil {
// 		logrus.Fatal(err)
// 	}

func ImportObjectsMongo(obj interface{}) {
	mongo, err := model.ConnectStorage(model.MongoDB)
	if err != nil {
		logrus.Fatal("Import to MongoDB failed: ", err)
	}
	err = mongo.Save(obj)
	if err != nil {
		logrus.Fatal(err)
	}
}

func ImportObjectsMsql(obj interface{}) {
	msql, err := model.ConnectStorage(model.MSQL)
	if err != nil {
		logrus.Fatal("Import to MSQL failed: ", err)
	}
	err = msql.Save(obj)
	if err != nil {
		logrus.Fatal(err)
	}
}

func ImportObjectsNeo4j(obj interface{}) {
	neo4j, err := model.ConnectStorage(model.Neo4J)
	if err != nil {
		logrus.Fatal("Import to Neo4j failed: ", err)
	}
	err = neo4j.Save(obj)
	if err != nil {
		logrus.Fatal(err)
	}
}

// }
