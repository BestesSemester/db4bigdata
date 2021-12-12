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

func ImportPersonsFromJSON(jsonfile string, persons *[]model.Person) {
	// Let's first read the `config.json` file
	content, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		logrus.Fatal("Error when opening file: ", err)
	}
	err = json.Unmarshal(content, &persons)
	if err != nil {
		logrus.Fatal(err)
	}
}

func ImportHierarchyFromJSON(jsonfile string, hierarchy *[]model.Hierarchy) {
	// Let's first read the `config.json` file
	content, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		logrus.Fatal("Error when opening file: ", err)
	}
	err = json.Unmarshal(content, &hierarchy)
	if err != nil {
		logrus.Fatal(err)
	}
}

func ImportInvoiceFromJSON(jsonfile string, invoices *[]model.Invoice) {
	// Let's first read the `config.json` file
	content, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		logrus.Fatal("Error when opening file: ", err)
	}
	err = json.Unmarshal(content, invoices)
	if err != nil {
		logrus.Fatal(err)
	}
}
