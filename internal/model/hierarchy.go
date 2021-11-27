package model

type Hierarchy struct {
	Agent Person
}

type Status int

const (
	inactive Status = 0
	active   Status = 1
)
