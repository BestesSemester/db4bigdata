package model

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type MsSQLConfig struct {
	URL      url.URL
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
	dsn := fmt.Sprintf("%s://%s:%s@%s?%s", conf.URL.Scheme, conf.UserName, conf.Password, conf.URL.Host, conf.URL.RawQuery)
	dbconn, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	mssql.db = dbconn
	mssql.db = mssql.db.Session(&gorm.Session{FullSaveAssociations: true})
	return mssql, nil
}

// Migrate - migrates the tables for each struct given
func (mssql *MsSQL) Migrate(inf ...interface{}) error {
	for _, inf := range inf {
		if err := mssql.db.AutoMigrate(inf); err != nil {
			return err
		}
	}
	return nil
}

// Save - used to save - TODO: implement for directs
func (mssql *MsSQL) Save(obj interface{}) error {
	// check for indirects (pointer)
	t := getDirectTypeFromInterface(obj)
	switch t.Kind() {
	// Check if the interface is iterable
	case reflect.Slice | reflect.Array:
		logrus.Println("detected iterables")
		if err := mssql.saveIterable(obj); err != nil {
			return err
		}
	case reflect.Struct:
		mssql.db.Save(obj)
	default:
		return fmt.Errorf("unsupported data type: %s", t.Kind())
	}
	return nil
}

func (mssql *MsSQL) saveIterable(obj interface{}) error {
	// iterate over the slice (has to be abstracted, because we are working type-agnostic)
	objs := getInterfacePointerSliceFromInterface(obj)
	for _, o := range objs {
		// save
		mssql.db.Save(o)
	}
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
func (mssql *MsSQL) Find(qry interface{}, target interface{}) error {
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
