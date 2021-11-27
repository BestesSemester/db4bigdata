package model

type Person struct {
	Name      string `bson:"Name"`
	FirstName string `bson:"FirstName"`
}
