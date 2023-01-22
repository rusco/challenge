// main package for loader application
// PROJECT = "CHALLENGE LOADER"
// AUTHOR  = "j.rebhan@gmail.com"
// VERSION = "1.0.0"
// DATE    = "2023-01-22 12:20:00"
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

const (
	DATABASE_NAME = "../data/challenge.db"
	START_ROW_0   = 0
	START_ROW_1   = 1
)

// utility function for "readCsvFile" function to load only certain columns
func in_array(val int, array []int) bool {
	for _, v := range array {
		if val == v {
			return true
		}
	}
	return false
}

// readCsvFile function to convert csv file into slice of string slices where only some columnIndexes are include
func readCsvFile(filePath string, includeIdx ...int) [][]string {

	loadAllColumns := len(includeIdx) == 0

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	var records [][]string
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		var rec []string
		for idx, val := range record {
			if loadAllColumns || in_array(idx, includeIdx) {
				rec = append(rec, val)
			}
		}
		records = append(records, rec)
	}
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
		return nil
	}
	return records
}

// insertTaxiZones function, loads zone data into db
func insertTaxiZones(fileName string) int {
	const (
		TAXI_DELETE_SQL = "DELETE FROM zones"
		TAXI_CREATE_SQL = "CREATE TABLE IF NOT EXISTS zones(locationid INTEGER PRIMARY KEY UNIQUE NOT NULL, borough TEXT, zone TEXT, service_zone TEXT)"
		TAXI_INSERT_SQL = "INSERT INTO zones(locationid, borough, zone, service_zone) values (?, ?, ?, ?)"
	)
	records := readCsvFile(fileName)
	createTable(DATABASE_NAME, TAXI_CREATE_SQL)
	deleteTableValues(DATABASE_NAME, TAXI_DELETE_SQL)
	insertRecords(DATABASE_NAME, TAXI_INSERT_SQL, records, START_ROW_1)
	return len(records)
}

// insertGreenTrips function, loads green trips data into db
func insertGreenTrips(fileName string) int {
	const (
		GREEN_DELETE_SQL = "DELETE FROM trips WHERE color = 'green'"
		GREEN_CREATE_SQL = "CREATE TABLE IF NOT EXISTS trips(pu_datetime TEXT,do_datetime TEXT, pu_locationid INTEGER, do_locationid INTEGER, color TEXT)"
		GREEN_INSERT_SQL = "INSERT INTO trips(pu_datetime, do_datetime, pu_locationid, do_locationid, color) values (?, ?, ?, ?, 'green')"
	)
	records := readCsvFile(fileName, 1, 2, 5, 6)
	createTable(DATABASE_NAME, GREEN_CREATE_SQL)
	deleteTableValues(DATABASE_NAME, GREEN_DELETE_SQL)
	insertRecords(DATABASE_NAME, GREEN_INSERT_SQL, records, START_ROW_0)
	return len(records)
}

// insertYellowTrips function, loads yellow trips data into db
func insertYellowTrips(fileName string) int {
	const (
		YELLOW_DELETE_SQL = "DELETE FROM trips WHERE color = 'yellow'"
		YELLOW_CREATE_SQL = "CREATE TABLE IF NOT EXISTS trips(pu_datetime TEXT,do_datetime TEXT, pu_locationid INTEGER, do_locationid INTEGER, color TEXT)"
		YELLOW_INSERT_SQL = "INSERT INTO trips(pu_datetime, do_datetime, pu_locationid, do_locationid, color) values (?, ?, ?, ?, 'yellow')"
	)
	records := readCsvFile(fileName, 1, 2, 7, 8)
	createTable(DATABASE_NAME, YELLOW_CREATE_SQL)
	deleteTableValues(DATABASE_NAME, YELLOW_DELETE_SQL)
	insertRecords(DATABASE_NAME, YELLOW_INSERT_SQL, records, START_ROW_0)
	return len(records)
}

// main function of loader to be run at comand line
// expects 2 arugments at comand line: 1. filetype: zone|green|yellow, 2. csvfile with complete path
// logs number of records and loading time
// repeted loads deletes previously loaded zones or green|yellow trip data
func main() {
	var (
		start        = time.Now()
		args         = os.Args
		countRecords = 0
	)

	if len(args) != 3 {
		fmt.Println(`2 arguments expected: argument1 = zone|green|yellow, argument 2 = filename.csv`)
		fmt.Println(`Examples:`)
		fmt.Println(` loader[.exe] zone   ../data/taxi_zone_lookup.csv`)
		fmt.Println(` loader[.exe] green  ../data/green_tripdata_2018-01_01-15.csv`)
		fmt.Println(` loader[.exe] yellow ../data/yellow_tripdata_2018-01_01-15.csv`)
		return
	}
	fileType, fileName := args[1], args[2]
	if fileType != "zone" && fileType != "green" && fileType != "yellow" {
		fmt.Println(`Error: 1st argument should be one of  "zone" or "green" or "yellow"`)
		return
	}
	if !strings.HasSuffix(fileName, ".csv") {
		fmt.Println(`Error: 2st argument should be a .csv file with the ending ".csv"`)
		return
	}
	switch fileType {
	case "zone":
		countRecords = insertTaxiZones(fileName)
	case "green":
		countRecords = insertGreenTrips(fileName)
	case "yellow":
		countRecords = insertYellowTrips(fileName)
	}
	log.Printf("done: %d Records of type '%s' loaded in %v.", countRecords, fileType, time.Since(start))
}
