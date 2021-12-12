package model

import (
	"net/url"
	"reflect"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type MsSQLConfig struct {
	Url      url.URL
	UserName string
	Password string
	Database string
}

type MsSQL struct {
	db *gorm.DB
}

func ConnectMsSQL(conf *MsSQLConfig) (Database, error) {
	mssql := &MsSQL{}
	// conn, err := sql.Open("mssql", conf.Database)
	dsn := "sqlserver://sa:1234myFancyPasswort@127.0.0.1?database=master"
	dbconn, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	mssql.db = dbconn
	mssql.db = mssql.db.Session(&gorm.Session{FullSaveAssociations: true})
	mssql.migrate()
	return mssql, nil
}

func (mssql *MsSQL) migrate() {
	mssql.db.AutoMigrate(&Role{})
	mssql.db.AutoMigrate(&Person{})
}

// TODO: implement save logic
func (mssql *MsSQL) Save(obj interface{}) error {
	return nil
}

func (mssql *MsSQL) SavePersons(persons *[]Person) error {
	for i := range *persons {
		p := *persons
		sl := p[i]
		// sl := p[0]
		// logrus.Println(p)
		mssql.db.Create(&sl)
	}
	return nil
}

// TODO: implement delete logic
func (mssql *MsSQL) Delete(obj interface{}) error {
	return nil
}

// Returns sql-Result
func (mssql *MsSQL) Find(qry string, target interface{}) error {
	// mssql.conn.Exec(qry)
	t := reflect.TypeOf(target)
	logrus.Println(t)
	logrus.Println(getAsAbstractStructFieldSetFromInterface(target))
	// logrus.Println(f.Tag.Get("mssql"))
	return nil
}

// Closes the database connection (should only be used if you close it on purpose)
func (mssql *MsSQL) Close() error {
	// err := mssql.conn.Close()
	return nil
}
