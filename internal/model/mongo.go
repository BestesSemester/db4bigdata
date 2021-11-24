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

type Mongo struct {
	conn *mongo.Client
}

func ConnectMongo(conf *MongoConfig) (Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@127.0.0.1"))
	// This does not work due to expected Database Object as return value
	return client, err
}

// TODO: implement save logic
func (mongo *Mongo) Save(obj interface{}) error {
	return nil
}

// TODO: implement delete logic
func (mongo *Mongo) Delete(obj interface{}) error {
	return nil
}

// Returns sql-Result
func (mongo *Mongo) Find(qry string, target interface{}) error {
	// mssql.conn.Exec(qry)
	t := reflect.TypeOf(target)
	logrus.Println(t)
	logrus.Println(getAsAbstractStructFieldSetFromInterface(target))
	// logrus.Println(f.Tag.Get("mssql"))
	return nil
}

// Closes the database connection (should only be used if you close it on purpose)
func (mongo *Mongo) Close() error {
	err := mongo.Close()
	return err
}
