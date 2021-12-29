package model

import (
	"fmt"
	"net/url"
	"reflect"

	pb "github.com/cheggaaa/pb/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

var depth = 0

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
		mssql.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(obj)
	default:
		return fmt.Errorf("unsupported data type: %s", t.Kind())
	}
	return nil
}

func (mssql *MsSQL) saveIterable(obj interface{}) error {
	// iterate over the slice (has to be abstracted, because we are working type-agnostic)
	objs := getInterfacePointerSliceFromInterface(obj)
	bar := pb.StartNew(len(objs))
	for _, o := range objs {
		// save
		mssql.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(o)
		bar.Increment()
	}
	bar.Finish()
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
			if field.tp.Type.Kind() == reflect.Ptr && (field.tp.Type.Elem().Kind() == reflect.Struct || field.tp.Type.Elem().Kind() == reflect.Slice || field.tp.Type.Elem().Kind() == reflect.Array) && field.tp.Tag.Get("gorm") != "-" {
				joinableFields = append(joinableFields, field.key)
				preloads := mssql.resolveStructFields(field, field.key, 6)
				if preloads != nil {
					preloadableFields = append(preloadableFields, preloads...)
				}
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

func (mssql *MsSQL) resolveStructFields(structure abstractStructField, parentname string, maxdepth int) []string {
	// logrus.Println(structure)
	if depth >= maxdepth {
		return nil
	}
	depth++
	preloadlist := []string{}
	parent := structure.tp
	parentType := parent.Type
	if parentType.Kind() == reflect.Ptr {
		parentType = parentType.Elem()
	}
	if parentType.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < parentType.NumField(); i++ {
		child := parentType.Field(i)
		if child.Type.Kind() == reflect.Ptr && (child.Type.Elem().Kind() == reflect.Struct || child.Type.Elem().Kind() == reflect.Slice || child.Type.Elem().Kind() == reflect.Array) && child.Tag.Get("gorm") != "-" {
			logrus.Println(child)
			preloadlist = append(preloadlist, parentname+"."+child.Name)
			field := abstractStructField{
				tp: child,
			}
			fieldnames := mssql.resolveStructFields(field, parentname+"."+child.Name, 10)
			// logrus.Println(fieldnames)
			if fieldnames != nil {
				preloadlist = append(preloadlist, fieldnames...)
			}
		}
	}
	return preloadlist
}

// Closes the database connection (should only be used if you close it on purpose)
func (mssql *MsSQL) Close() error {
	// err := mssql.conn.Close()
	return nil
}
