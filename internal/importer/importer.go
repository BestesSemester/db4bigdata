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
	db, err := model.ConnectStorage(model.MSQL)
	if err != nil {
		logrus.Fatal("Error in saving", err)
	}
	err = db.Save(persons)
	if err != nil {
		logrus.Fatal(err)
	}

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
	db, err := model.ConnectStorage(model.MSQL)
	if err != nil {
		logrus.Fatal("Error in saving", err)
	}
	err = db.Save(hierarchy)
	if err != nil {
		logrus.Fatal(err)
	}

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
	db, err := model.ConnectStorage(model.MSQL)
	if err != nil {
		logrus.Fatal("Error in saving", err)
	}
	err = db.Save(invoice)
	if err != nil {
		logrus.Fatal(err)
	}

}

func ImportProvisiondistributionFromJSON(jsonfile string) {
	// Let's first read the `config.json` file
	content, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		logrus.Fatal("Error when opening file: ", err)
	}
	var provdib []model.Provdib
	err = json.Unmarshal(content, &provdib)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Println(len(provdib))
	db, err := model.ConnectStorage(model.MSQL)
	if err != nil {
		logrus.Fatal("Error in saving", err)
	}
	err = db.Save(provdib)
	if err != nil {
		logrus.Fatal(err)
	}

}
