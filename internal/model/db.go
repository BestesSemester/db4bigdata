package model

import "reflect"

type Database interface {
	Save(obj interface{}) error
	Delete(obj interface{}) error
	Find(qry string, obj interface{}) error
	Close() error
}

type StorageType int

const (
	MongoDB StorageType = 0
	MSQL    StorageType = 1
	Neo4J   StorageType = 2
)

func ConnectStorage(st StorageType) (Database, error) {
	// connect
	switch st {
	case MongoDB:
		conf := MongoConfig{}
		// if passwd, ok := url.URL.User.Password(); ok {
		// 	conf.AccessKeySecret = passwd
		// }
		return ConnectMongo(&conf)
	case MSQL:
		conf := MsSQLConfig{}
		// if passwd, ok := url.URL.User.Password(); ok {
		// 	conf.Password = passwd
		// }
		return ConnectMsSQL(&conf)
	case Neo4J:
		conf := Neo4jConfig{}
		return ConnectNeo4j(&conf)
	default:
		return nil, nil
	}
}

func resolveStructFields(inf interface{}) []reflect.StructField {
	t := reflect.TypeOf(inf)
	var strct reflect.Type
	fields := []reflect.StructField{}
	if t.Kind() == reflect.Ptr {
		strct = t.Elem()
	} else {
		strct = t
	}
	for i := 0; i < strct.NumField(); i++ {
		fields = append(fields, strct.Field(i))
	}
	return fields
}
