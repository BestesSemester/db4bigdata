package model

import (
	"context"
	"net/url"
	"reflect"
	"time"

	// "github.com/kamva/mgm/v3"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Url      url.URL
	UserName string
	Password string
}

type MyMongo struct {
	conn    *mongo.Client
	context *context.Context
}

func ConnectMongo(conf *MongoConfig) (Database, error) {
	// Setup the mgm default config
	// err := mgm.SetDefaultConfig(nil, "mgm_lab", options.Client().ApplyURI("mongodb://root:example@127.0.0.1"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@127.0.0.1"))

	myMongo := &MyMongo{client, &ctx}
	return myMongo, err
}

// I would call this method import
func (mongo *MyMongo) Save(obj interface{}) error {
	t := reflect.TypeOf(obj)
	coll := mongo.conn.Database("my_go_db").Collection(t.Elem().Name())
	switch t.Kind() {
	case reflect.Slice:
		objs := getInterfaceSliceFromInterface(obj)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := coll.InsertMany(ctx, objs)
		if err != nil {
			logrus.Fatal("Error when inserting objects: ", err)
			return err
		} else {
			logrus.Println("Inserted ", len(res.InsertedIDs), " documents for Collection \"", t.Elem().Name(), "\"")
		}
	case reflect.Array:

	default:
	}
	return nil
}

// TODO: implement delete logic
func (mongo *MyMongo) Delete(obj interface{}) error {
	return nil
}

// Returns sql-Result
func (mongo *MyMongo) Find(qry string, target interface{}) error {
	// mssql.conn.Exec(qry)
	t := reflect.TypeOf(target)
	logrus.Println(t)
	logrus.Println(getAsAbstractStructFieldSetFromInterface(target))
	// logrus.Println(f.Tag.Get("mssql"))
	return nil
}

// Closes the database connection (should only be used if you close it on purpose)
func (mongo *MyMongo) Close() error {
	err := mongo.conn.Disconnect(*mongo.context)
	return err
}
