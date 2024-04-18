package sqlite

import (
	"database/sql"
	"fmt"
	"os"
)

type Supabase struct {}

var conn Supabase

func (repo *Supabase) Connect() (*sql.DB, error) {
	
	connStr := "user="+ os.Getenv("POSTGRES_USERNAME") +" password="+ os.Getenv("POSTGRES_PASSWORD") +" host="+ os.Getenv("POSTGRES_HOST") +" port="+ os.Getenv("POSTGRES_PORT") +" dbname="+ os.Getenv("POSTGRES_DB") +""
	
	db, err := sql.Open("postgres", connStr)
	
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}