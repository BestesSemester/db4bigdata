package model

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

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
		mssql.db.Create(o)
	}
	return nil
}

// TODO: implement delete logic
func (mssql *MsSQL) Delete(obj interface{}) error {
	return nil
}

// Returns sql-Result
func (mssql *MsSQL) Find(qry interface{}, target interface{}) error {
	t := reflect.TypeOf(target)
	logrus.Printf("%d", t.Kind())
	switch t.Kind() {
	case reflect.Ptr:
		fs := getAsAbstractStructFieldSetFromInterface(target)
		joinableFields := []string{}
		preloadableFields := []string{}
		for _, field := range fs.fields {
			if field.tp.Type.Kind() == reflect.Ptr && field.tp.Tag.Get("gorm") != "-" {
				joinableFields = append(joinableFields, field.key)
				_, preloads := mssql.resolveStructFields(field)
				preloadableFields = append(preloadableFields, preloads...)
			}
		}
		logrus.Printf("Joining: %s", joinableFields)
		logrus.Printf("Preloading: %s", preloadableFields)
		tx := mssql.db.Set("gorm:auto_preload", true)
		for _, joinableField := range joinableFields {
			tx = tx.Joins(joinableField)
		}
		for _, preloadField := range preloadableFields {
			tx = tx.Preload(preloadField)
		}
		tx.Debug().Where(qry).Find(&target)
	default:
		logrus.Errorln("no such implementation")
	}
	return nil
}

func (mssql *MsSQL) resolveStructFields(structure abstractStructField) ([]string, []string) {
	logrus.Println(structure)
	joinlist := []string{}
	preloadlist := []string{}
	parent := structure.tp
	for i := 0; i < parent.Type.Elem().NumField(); i++ {
		child := parent.Type.Elem().Field(i)
		if child.Type.Kind() == reflect.Ptr && child.Tag.Get("gorm") != "-" {
			logrus.Println(child)
			joinTableName := parent.Name + "_" + child.Name
			joinlist = append(joinlist, "JOIN "+strings.ToLower(child.Name)+"s "+joinTableName+" on "+strings.ToLower(parent.Name)+"."+strings.ToLower(child.Name)+"_id="+joinTableName+"."+strings.ToLower(child.Name)+"_id")
			preloadlist = append(preloadlist, parent.Name+"."+child.Name)
		}
	}
	return joinlist, preloadlist
}

// Closes the database connection (should only be used if you close it on purpose)
func (mssql *MsSQL) Close() error {
	// err := mssql.conn.Close()
	return nil
}
