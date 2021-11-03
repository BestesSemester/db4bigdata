package model

import "net/url"

type Neo4jConfig struct {
	Url      url.URL
	UserName string
	Password string
}

func ConnectNeo4j(conf *Neo4jConfig) (Database, error) {
	return nil, nil
}