// main package for loader application
// PROJECT = "CHALLENGE LOADER"
// AUTHOR  = "j.rebhan@gmail.com"
// VERSION = "1.0.0"
// DATE    = "2023-01-22 12:20:00"
package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" //pure go sqlite driver for compatiblity reasons, cgo based "https://github.com/mattn/go-sqlite3" is faster
)

const SQLITE = "sqlite"

func createTable(dbname string, createSql string) {

	db, err := sql.Open(SQLITE, dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := db.Exec(createSql); err != nil {
		log.Fatal(err)
	}
}

func deleteTableValues(dbname string, deleteSql string) {

	db, err := sql.Open(SQLITE, dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := db.Exec(deleteSql); err != nil {
		log.Fatal(err)
	}
}

func insertRecords(dbname string, insertSql string, records [][]string, startRow int) {

	records = records[startRow:]

	db, err := sql.Open(SQLITE, dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(insertSql)
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {

		var rec []any
		for i := 0; i < len(record); i++ {
			rec = append(rec, record[i])
		}
		_, err = stmt.Exec(rec...)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := stmt.Close(); err != nil {
		log.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
