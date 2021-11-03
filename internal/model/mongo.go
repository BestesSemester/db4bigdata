package model

import "net/url"

type MongoConfig struct {
	Url      url.URL
	UserName string
	Password string
}

func ConnectMongo(conf *MongoConfig) (Database, error) {
	return nil, nil
}
