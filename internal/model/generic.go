package model

import "github.com/mindstand/gogm/v2"

type Neo4jBaseNode struct {
	Id *int64 `json:"-" gogm:"pk=default" gorm:"-" bson:"-"`
	// LoadMap represents the state of how a node was loaded for neo4j.
	// This is used to determine if relationships are removed on save
	// field -- relations
	LoadMap map[string]*gogm.RelationConfig `json:"-" gogm:"-" gorm:"-" bson:"-"`
}
