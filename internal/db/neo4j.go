package db

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"git.sys-tem.org/caos/db4bigdata/internal/model"
	"github.com/mindstand/gogm/v2"
	"github.com/sirupsen/logrus"
)

type Neo4jConfig struct {
	URL      url.URL
	UserName string
	Password string
}

type Neo4j struct {
	db      *gogm.Gogm
	session gogm.SessionV2
}

func ConnectNeo4j(conf *Neo4jConfig) (Database, error) {
	iport, err := strconv.Atoi(conf.URL.Port())
	if err != nil {
		return nil, err
	}
	config := gogm.Config{
		Host:             conf.URL.Hostname(),
		Port:             iport,
		Username:         conf.UserName,
		Password:         conf.Password,
		PoolSize:         4,
		Protocol:         "bolt",
		EnableDriverLogs: true,
		TargetDbs:        []string{"my_test_db"},
		LogLevel:         "DEBUG",
		IndexStrategy:    gogm.IGNORE_INDEX,
	}

	conn, err := gogm.New(&config, gogm.DefaultPrimaryKeyStrategy, &model.Person{}, &model.Role{}, &model.Invoice{})
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	session, err := conn.NewSessionV2(gogm.SessionConfig{AccessMode: gogm.AccessModeWrite})
	if err != nil {
		return nil, err
	}

	neo4j := &Neo4j{
		db:      conn,
		session: session,
	}
	return neo4j, nil
}

// TODO: implement save logic
func (neo4j *Neo4j) Save(obj interface{}) error {
	logrus.Println(reflect.TypeOf(obj).Kind())
	t := getDirectTypeFromInterface(obj)
	switch t.Kind() {
	case reflect.Slice | reflect.Array:
		logrus.Println("found iterable")
		objs := getInterfacePointerSliceFromInterface(obj)
		for i, o := range objs {
			logrus.Printf("Saving object no. %d", i)
			logrus.Println(o)
			err := neo4j.session.SaveDepth(context.Background(), o, 1)
			if err != nil {
				logrus.Errorln(err)
				return err
			}
		}
	case reflect.Struct:
		if err := neo4j.session.SaveDepth(context.Background(), obj, 1); err != nil {
			logrus.Errorln(err)
			return err
		}
	default:
		return fmt.Errorf("no compatible type given")
	}
	return nil
}

// TODO: implement delete logic
func (neo4j *Neo4j) Delete(obj interface{}) error {
	return nil
}

// Find - Does nothing
func (neo4j *Neo4j) Find(qry interface{}, target interface{}) error {
	return nil
}

//Migrate - does nothing
func (neo4j *Neo4j) Migrate(inf ...interface{}) error {
	return fmt.Errorf("no implementation")
}

//Exec - Ecexutes Cypher query
func (neo4j *Neo4j) Exec(qry string, inf interface{}) error {
	var qryInterface map[string]interface{}
	inrec, err := json.Marshal(inf)
	if err != nil {
		logrus.Errorln(err)

		return err
	}
	err = json.Unmarshal(inrec, &qryInterface)
	if err != nil {
		logrus.Errorln(err)

		return err
	}

	res, _, err := neo4j.session.QueryRaw(context.Background(), qry, nil)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	if len(res) < 1 || len(res[0]) < 1 {
		return nil
	}
	resj, _ := json.MarshalIndent(&res[0][0], "", "	")
	err = json.Unmarshal(resj, &inf)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	return nil
}

// Closes the database connection (should only be used if you close it on purpose)
func (neo4j *Neo4j) Close() error {
	return nil
}
