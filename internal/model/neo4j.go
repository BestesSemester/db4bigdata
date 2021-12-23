package model

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/mindstand/gogm/v2"
	"github.com/sirupsen/logrus"
)

type Neo4jBaseNode struct {
	Id *int64 `json:"-" gogm:"pk=default" gorm:"-" bson:"-"`

	// LoadMap represents the state of how a node was loaded for neo4j.
	// This is used to determine if relationships are removed on save
	// field -- relations
	LoadMap map[string]*gogm.RelationConfig `json:"-" gogm:"-" gorm:"-" bson:"-"`
}

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

	conn, err := gogm.New(&config, gogm.DefaultPrimaryKeyStrategy, &Person{}, &Role{}, &Invoice{})
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
	case reflect.Array | reflect.Slice:
		logrus.Println("found iterable")
		objs := getInterfacePointerSliceFromInterface(obj)
		for i, o := range objs {
			logrus.Printf("Saving object no. %i", i)
			err := neo4j.session.SaveDepth(context.Background(), o, 10)
			if err != nil {
				logrus.Errorln(err)
				return err
			}
		}
	case reflect.Struct:
		if err := neo4j.session.SaveDepth(context.Background(), obj, 2); err != nil {
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

// Returns Neo4j-Result
func (neo4j *Neo4j) Find(qry interface{}, target interface{}) error {
	// mssql.conn.Exec(qry)
	t := reflect.TypeOf(target)
	logrus.Println(t)
	logrus.Println(getAsAbstractStructFieldSetFromInterface(target))
	// logrus.Println(f.Tag.Get("mssql"))
	return nil
}

//Migrate - does nothing
func (neo4j *Neo4j) Migrate(inf ...interface{}) error {
	return fmt.Errorf("no implementation")
}

// Closes the database connection (should only be used if you close it on purpose)
func (neo4j *Neo4j) Close() error {
	return nil
}
