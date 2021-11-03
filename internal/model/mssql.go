package model

import (
	"database/sql"
	"log"
	"net/url"
	"reflect"

	"github.com/sirupsen/logrus"
)

type MsSQLConfig struct {
	Url      url.URL
	UserName string
	Password string
	Database string
}

type MsSQL struct {
	conn *sql.DB
}

func ConnectMsSQL(conf *MsSQLConfig) (Database, error) {
	mssql := &MsSQL{}
	// conn, err := sql.Open("mssql", conf.Database)
	mssql.conn = nil
	return mssql, nil
}

// TODO: implement save logic
func (mssql *MsSQL) Save(obj interface{}) error {
	return nil
}

// TODO: implement delete logic
func (mssql *MsSQL) Delete(obj interface{}) error {
	return nil
}

// Returns sql-Result
func (mssql *MsSQL) Find(qry string, target interface{}) error {
	mssql.conn.Exec(qry)
	t := reflect.TypeOf(target)
	logrus.Println(t)
	fields := resolveStructFields(target)
	// result, err := nil
	for _, field := range fields {
		log.Println(field)
	}
	// logrus.Println(f.Tag.Get("mssql"))
	return nil
}

// Closes the database connection (should only be used if you close it on purpose)
func (mssql *MsSQL) Close() error {
	err := mssql.conn.Close()
	return err
}
