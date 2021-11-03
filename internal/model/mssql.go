package model

import "net/url"

type MsSQLConfig struct {
	Url      url.URL
	UserName string
	Password string
}

func ConnectMsSQL(conf *MsSQLConfig) (Database, error) {
	return nil, nil
}
