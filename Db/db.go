package Db

import (
	a "github.com/varunturlapati/simpleWebSvc/Article"
	"time"
)

type Entry struct {
	Id    string
	Value a.Article
	Ts    time.Time
}

type Db interface {
	PingPong() (string, error)
	AddEntry(*Entry) error
	RemoveEntry(string) error
	ChangeEntry(string, *Entry) error
}
