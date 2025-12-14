package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

type Filter struct {
	Key   string
	Value string
}

type Paginate struct {
	Limit   int
	Offset  int
	Filters []Filter
}

func Connect(url string) {
	var err error
	db, err = sqlx.Connect("postgres", url)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func DB() *sqlx.DB {
	return db
}
