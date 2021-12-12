package model

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/mindstand/gogm/v2"
	"github.com/sirupsen/logrus"
)

type Neo4jConfig struct {
	URL      url.URL
	UserName string
	Password string
}

type Neo4j struct {
	db *gogm.Gogm
}

func ConnectNeo4j(conf *Neo4jConfig) (Database, error) {
	iport, err := strconv.Atoi(conf.URL.Port())
	if err != nil {
		return nil, err
	}
	config := gogm.Config{
		Host:     conf.URL.Host,
		Port:     iport,
		Username: conf.UserName,
		Password: conf.Password,
	}
	conn, err := gogm.New(&config, gogm.DefaultPrimaryKeyStrategy, &Person{}, &Role{}, &Invoice{}, &Hierarchy{}, &Provision{})
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	neo4j := &Neo4j{
		db: conn,
	}
	return neo4j, nil
}

// TODO: implement save logic
func (neo4j *Neo4j) Save(obj interface{}) error {
	logrus.Println(obj)
	return nil
}

func (neo4j *Neo4j) SavePersons(persons *[]Person) error {
	return nil
}

// TODO: implement delete logic
func (neo4j *Neo4j) Delete(obj interface{}) error {
	return nil
}

// Returns sql-Result
func (neo4j *Neo4j) Find(qry string, target interface{}) error {
	// mssql.conn.Exec(qry)
	t := reflect.TypeOf(target)
	logrus.Println(t)
	logrus.Println(getAsAbstractStructFieldSetFromInterface(target))
	// logrus.Println(f.Tag.Get("mssql"))
	return nil
}

func (neo4j *Neo4j) Migrate(inf ...interface{}) error {
	return fmt.Errorf("no implementation")
}

// Closes the database connection (should only be used if you close it on purpose)
func (neo4j *Neo4j) Close() error {
	return nil
}
