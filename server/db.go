// Challenge App main Package for Server Application
// PROJECT = "CHALLENGE SERVER"
// AUTHOR  = "j.rebhan@gmail.com"
// VERSION = "1.0.0"
// DATE    = "2023-01-22 17:30"
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite" //sqlite driver
)

const (
	SQLITE        = "sqlite"
	VERSION       = "select sqlite_version()"
	DATABASE_NAME = "../data/challenge.db"
)

// get database version function
func version(dbname string) string {

	db, err := sql.Open(SQLITE, dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var ver string
	rows, err := db.Query(VERSION)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&ver)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return ver
}

type Trip struct {
	Zone     string `json:"zone"`
	Pu_total int32  `json:"pu_total"`
	Do_total int32  `json:"do_total"`
}

// getTopZonesDb function, returns slice of Trip structs
func getTopZonesDb(dbname string, order string, showLog bool) []Trip {

	puSql := `WITH dof AS (
		SELECT z.locationid,
			   count(t.do_locationid) do_total
		  FROM trips t,
			   zones z
		 WHERE t.do_locationid = z.locationid
		 GROUP BY z.locationid
	)
	SELECT z.zone,
		   count(t.pu_locationid) pu_total,
		   IFNULL(dof.do_total, 0) do_total
	  FROM trips t,
		   zones z
	LEFT JOIN dof ON t.pu_locationid = dof.locationid
	 WHERE t.pu_locationid = z.locationid
	 GROUP BY z.zone
	 ORDER BY count(t.pu_locationid) DESC
	 LIMIT 5`

	doSql := `WITH puf AS (
		SELECT z.locationid,
			   count(t.pu_locationid) pu_total
		  FROM trips t,
			   zones z
		 WHERE t.pu_locationid = z.locationid
		 GROUP BY z.locationid
	)
	SELECT z.zone,
		   count(t.do_locationid) do_total,
		   IFNULL(puf.pu_total, 0) pu_total
	  FROM trips t,
		   zones z
		   LEFT JOIN
		   puf ON t.do_locationid = puf.locationid
	 WHERE t.do_locationid = z.locationid
	 GROUP BY z.zone
	 ORDER BY count(t.do_locationid) DESC
	 LIMIT 5`

	var stmt string
	if order == "pickups" {
		stmt = puSql
	} else if order == "dropoffs" {
		stmt = doSql
	} else {
		log.Fatal(errors.New("'order' parameter sould be either 'pickups' or 'dropoffs'"))
	}
	if showLog {
		log.Printf("getTopZonesDb sql: %s", stmt)
	}

	db, err := sql.Open(SQLITE, dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	start := time.Now()
	rows, err := db.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	sTrip := make([]Trip, 0)
	for rows.Next() {

		var zone string
		var pu_total int32
		var do_total int32

		err = rows.Scan(&zone, &do_total, &pu_total)
		if err != nil {
			log.Fatal(err)
		}
		trip := Trip{zone, do_total, pu_total}
		sTrip = append(sTrip, trip)

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("done in %v", time.Since(start))
	return sTrip
}

type ZoneTrip struct {
	Zone string `json:"zone"`
	Date string `json:"date"`
	Pu   int32  `json:"pu"`
	Do   int32  `json:"do"`
}

// getZoneTripsDb function, returns ZoneTrip struct
func getZoneTripsDb(dbname string, zoneparam int, dateparam string, showLog bool) ZoneTrip {

	zoneTripSql := `WITH res AS (
		SELECT z.zone,
			   count(trips_pu.pu_locationid) pu
		  FROM trips trips_pu,
			   zones z
		 WHERE trips_pu.pu_locationid 		= z.locationid AND
			   trips_pu.pu_locationid 		= %d AND
			   date(trips_pu.pu_datetime) 	= date('%s')
		)
		SELECT z.zone					 	zone,
		   res.pu 						 	pu_count,
		   count(trips_do.do_locationid) 	do_count
	  	FROM trips trips_do,
		   zones z,
		   res
	 	WHERE trips_do.do_locationid 		= z.locationid AND
		   trips_do.do_locationid 			= %d AND
		   date(trips_do.do_datetime) 		= date('%s')`

	stmt := fmt.Sprintf(zoneTripSql, zoneparam, dateparam, zoneparam, dateparam)
	if showLog {
		log.Printf("getZoneTripsDb sql: %s", stmt)
	}
	db, err := sql.Open(SQLITE, dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	start := time.Now()
	rows, err := db.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var zoneTrip ZoneTrip
	for rows.Next() {

		var zone string
		var pu_count int32
		var do_count int32

		err = rows.Scan(&zone, &pu_count, &do_count)
		if err != nil {
			log.Fatal(err)
		}
		zoneTrip = ZoneTrip{zone, dateparam, pu_count, do_count}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("done in %v", time.Since(start))
	return zoneTrip
}

type YellowTrip struct {
	Pu_datetime   string `json:"pu_datetime"`
	Do_datetime   string `json:"do_datetime"`
	Pu_locationid int32  `json:"pu_locationid"`
	Do_locationid int32  `json:"do_locationid"`
}

// getListYellowDb function, returns slice of YellowTrip structs
func getListYellowDb(dbname string, whereSql string, showLog bool) []YellowTrip {
	listYellowSql := `SELECT 
		pu_datetime,
		do_datetime,
		pu_locationid,
		do_locationid
	FROM trips t
	WHERE COLOR = 'yellow' 
		%s`

	stmt := fmt.Sprintf(listYellowSql, whereSql)
	if showLog {
		log.Printf("getListYellowDb sql: %s", stmt)
	}

	db, err := sql.Open(SQLITE, dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	start := time.Now()
	rows, err := db.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	YellowTripSlc := make([]YellowTrip, 0)
	for rows.Next() {

		var pu_datetime string
		var do_datetime string
		var pu_locationid int32
		var do_locationid int32

		err = rows.Scan(&pu_datetime, &do_datetime, &pu_locationid, &do_locationid)
		if err != nil {
			log.Fatal(err)
		}
		yellowTrip := YellowTrip{pu_datetime, do_datetime, pu_locationid, do_locationid}
		YellowTripSlc = append(YellowTripSlc, yellowTrip)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("done in %v", time.Since(start))
	return YellowTripSlc
}
