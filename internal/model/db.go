package model

import (
	"net/url"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type Database interface {
	Save(obj interface{}) error
	Delete(obj interface{}) error
	Find(qry interface{}, obj interface{}) error
	Migrate(inf ...interface{}) error
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
	godotenv.Load(".env")
	// connect
	switch st {
	case MongoDB:
		url, err := url.Parse(os.Getenv("MONGO_URL"))
		if err != nil {
			return nil, err
		}
		conf := MongoConfig{
			URL:      *url,
			UserName: os.Getenv("MONGO_USER"),
			Password: os.Getenv("MONGO_PASSWORD"),
		}
		// if passwd, ok := url.URL.User.Password(); ok {
		// 	conf.AccessKeySecret = passwd
		// }
		return ConnectMongo(&conf)
	case MSQL:
		url, err := url.Parse(os.Getenv("MSSQL_URL"))
		if err != nil {
			return nil, err
		}
		conf := MsSQLConfig{
			URL:      *url,
			UserName: os.Getenv("MSSQL_USER"),
			Password: os.Getenv("MSSQL_PASSWORD"),
		}
		// if passwd, ok := url.URL.User.Password(); ok {
		// 	conf.Password = passwd
		// }
		return ConnectMsSQL(&conf)
	case Neo4J:
		url, err := url.Parse(os.Getenv("NEO4J_URL"))
		if err != nil {
			return nil, err
		}
		conf := Neo4jConfig{
			URL:      *url,
			UserName: os.Getenv("NEO4J_USER"),
			Password: os.Getenv("NEO4J_PASSWORD"),
		}
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
	var tp reflect.Type
	t := reflect.TypeOf(inf)
	if t.Kind() == reflect.Ptr {
		// logrus.Println("converting")
		tp = t.Elem()
	} else {
		tp = t
	}
	return tp
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

// inf has to be kind of SLICE
func getInterfacePointerSliceFromInterface(inf interface{}) []interface{} {
	v := getDirectStructFromInterface(inf)
	var objs []interface{}
	for i := 0; i < v.Len(); i++ {
		objs = append(objs, v.Index(i).Addr().Interface())
	}
	return objs
}

// inf has to be kind of SLICE
func getInterfaceSliceFromInterface(inf interface{}) []interface{} {
	v := getDirectStructFromInterface(inf)
	var objs []interface{}
	for i := 0; i < v.Len(); i++ {
		objs = append(objs, v.Index(i).Interface())
	}
	return objs
}
