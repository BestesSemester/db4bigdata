package model

import (
	"reflect"
)

type Database interface {
	Save(obj interface{}) error
	Delete(obj interface{}) error
	Find(qry string, obj interface{}) error
	Close() error
}

type abstractStructFieldSet struct {
	fields []abstractStructField
}
type abstractStructField struct {
	key   string
	value interface{}
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
	strct := getDirectTypeFromInterface(inf)
	fields := []reflect.StructField{}
	for i := 0; i < strct.NumField(); i++ {
		fields = append(fields, strct.Field(i))
	}
	return fields
}

func getAsAbstractStructFieldSetFromInterface(inf interface{}) abstractStructFieldSet {
	fields := resolveStructFields(inf)
	// result, err := nil
	afs := abstractStructFieldSet{}
	for k, field := range fields {
		f := abstractStructField{
			key:   field.Name,
			value: getDirectStructFromInterface(inf).Field(k).String(),
		}
		afs.fields = append(afs.fields, f)
	}
	return afs
}

func getDirectTypeFromInterface(inf interface{}) reflect.Type {
	var strct reflect.Type
	t := reflect.TypeOf(inf)
	if t.Kind() == reflect.Ptr {
		strct = t.Elem()
	} else {
		strct = t
	}
	return strct
}

func getDirectStructFromInterface(inf interface{}) reflect.Value {
	var strct reflect.Value
	t := reflect.TypeOf(inf)
	if t.Kind() == reflect.Ptr {
		strct = reflect.ValueOf(inf).Elem()
	} else {
		strct = reflect.ValueOf(inf)
	}
	return strct
}
