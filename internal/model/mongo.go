package model

import (
	"context"
	"net/url"
	"reflect"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@127.0.0.1"))
	// This does not work due to expected Database Object as return value
	myMongo := &MyMongo{client, &ctx}
	return myMongo, err
}

// TODO: implement save logic
func (mongo *MyMongo) Save(obj interface{}) error {
	db := mongo.conn.Database("my_go_db")
	coll := db.Collection("persons")

	res, err := coll.InsertOne(*mongo.context, obj.data)
	if err != nil {
		logrus.Fatal("Error when opening file: ", err)
	}
	logrus.Print(res.InsertedID)
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
