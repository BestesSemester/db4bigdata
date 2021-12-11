package model

import (
	"context"
	"net/url"
	"reflect"
	"time"

	// "github.com/kamva/mgm/v3"
	// "git.sys-tem.org/caos/db4bigdata/internal/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName string = "myDB"

type MongoConfig struct {
	Url      url.URL
	UserName string
	Password string
}

type MyMongo struct {
	conn    *mongo.Client
	context *context.Context
}

// This function is not tested
func Initialize(mongo *MyMongo) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := mongo.conn.Database(dbName).Drop(ctx)

	if err == nil {
		return true
	}
	return false
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
	coll := mongo.conn.Database(dbName).Collection(t.Elem().Name())
	switch t.Kind() {
	case reflect.Slice:
		objs := getInterfaceSliceFromInterface(obj)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := coll.InsertMany(ctx, objs)
		if err != nil {
			logrus.Fatal("Error when inserting objects: ", err)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coll := mongo.conn.Database(dbName).Collection(("Person"))

	cursor, err := coll.Find(ctx, bson.D{{"name", qry}})
	if err = cursor.All(ctx, target); err != nil {
		logrus.Fatal(err)
	}
	defer cursor.Close(ctx)

	return nil
}

// Closes the database connection (should only be used if you close it on purpose)
func (mongo *MyMongo) Close() error {
	err := mongo.conn.Disconnect(*mongo.context)
	return err
}
