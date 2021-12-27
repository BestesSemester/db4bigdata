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
		mssql.db.Create(o)
	}
	return nil
}

// TODO: implement delete logic
func (mssql *MsSQL) Delete_Statement(obj interface{}) error {
	mssql.db.Delete(&obj)
	return nil
}

// Returns sql-Result
func (mssql *MsSQL) Find(qry interface{}, target interface{}) error {
	var agents []string
	i_t := Invoice{}
	//select the agents in invoices
	mssql.db.Model(&i_t).Distinct().Pluck("Agent_ID", &agents)

	//write all provisions into table
	for _, agent := range agents {
		//logrus.Println(agent)

		mssql.db.Exec(`WITH  temp2 as( select agent_id , supervisor_id
				from hierarchies
				where agent_id = ?
				union all
				select a.agent_id, a.supervisor_id
				from hierarchies a inner join temp2 on temp2.supervisor_id = a.agent_id
						)
		insert into provision_distributions
			select i.[Invoice_ID], t.agent_id, 
			case 
				when (t.agent_id = ? and (select count(*) from temp2) > 1 )
					then i.net_Sum * 0.7 * 0.1	
				when (t.agent_id = 1080
					and (select count(*) from temp2) = 1) then i.net_Sum * 0.1			
				else i.net_Sum *0.1 * 0.3/((select count(*) from temp2)-1) 
			end provision
			from temp2 t, [dbo].[invoices] i
			where i.agent_id = ? 	
			and invoice_id not in (select invoice_id from provision_distributions)	
			order by Invoice_ID`, agent, agent, agent)
	}

	t := reflect.TypeOf(target)
	logrus.Printf("%d", t.Kind())
	switch t.Kind() {
	case reflect.Ptr:
		logrus.Println("Here comes the statement:")

		mssql.db.Where(qry).Find(&target)
		/*fs := getAsAbstractStructFieldSetFromInterface(target)
		joinableFields := []string{}
		for _, field := range fs.fields {
			if field.tp.Type.Kind() == reflect.Ptr && field.tp.Tag.Get("gorm") != "-" {
				joinableFields = append(joinableFields, field.key)
			}
		}
		logrus.Println(joinableFields)
		var tx *gorm.DB
		for _, preloadField := range joinableFields {
			tx = mssql.db.Preload(preloadField).Joins(preloadField)
		}
		tx.Where(qry).First(&target)*/

	default:
		logrus.Errorln("no such implementation")
	}
	// logrus.Println(f.Tag.Get("mssql"))
	return nil
}

// Closes the database connection (should only be used if you close it on purpose)
func (mssql *MsSQL) Close() error {
	// err := mssql.conn.Close()
	return nil
}
