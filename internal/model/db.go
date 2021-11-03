package model

type Database interface {
	Save(obj interface{})
	Delete(obj interface{})
	Find(qry string)
	Close()
}

type StorageType int

const (
	mongodb StorageType = 0
	mssql   StorageType = 1
	neo4j   StorageType = 3
)

func ConnectStorage(st StorageType) (Database, error) {
	// connect
	switch st {
	case mongodb:
		conf := MongoConfig{}
		// if passwd, ok := url.URL.User.Password(); ok {
		// 	conf.AccessKeySecret = passwd
		// }
		return ConnectMongo(&conf)
	case mssql:
		conf := MsSQLConfig{}
		// if passwd, ok := url.URL.User.Password(); ok {
		// 	conf.Password = passwd
		// }
		return ConnectMsSQL(&conf)
	case neo4j:
		conf := Neo4jConfig{}
		return ConnectNeo4j(&conf)
	default:
		return nil, nil
	}
}
