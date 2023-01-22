// Challenge App main Package for Server Application
// PROJECT = "CHALLENGE SERVER"
// AUTHOR  = "j.rebhan@gmail.com"
// VERSION = "1.0.0"
// DATE    = "2023-01-22 17:30"
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/browser"
)

var ( //defaults if no values in "server.config.json" file specified
	LOGSQL           = true
	PORT             = 8080
	OPENLOCALBROWSER = false
)

type ServerConf struct {
	LogSql           bool `json:"logsql"`
	Port             int  `json:"port"`
	OpenLocalBrowser bool `json:"openlocalbrowser"`
}

// read server configuration from file "server.config.json"
func readConfig() ServerConf {
	var conf ServerConf

	jsonString, err := os.ReadFile("./server.config.json")
	if err != nil {
		fmt.Print(err)
		return ServerConf{LOGSQL, PORT, OPENLOCALBROWSER}
	}
	err = json.Unmarshal(jsonString, &conf)
	if err != nil {
		fmt.Print(err)
		return ServerConf{LOGSQL, PORT, OPENLOCALBROWSER}
	}
	return conf
}

// about handler
func getAbout(c echo.Context) error {
	params := map[string]any{"App": "ChallengeApp", "date": "2023-21-01", "Author": "j.rebhan@gmail.com", "Version": 1.0}
	return c.JSON(http.StatusOK, params)
}

// dbversion handler
func getDbVersion(c echo.Context) error {
	version := version(DATABASE_NAME)
	params := map[string]any{"Database": "Sqlite", "Version": version}
	return c.JSON(http.StatusOK, params)
}

// getTopZones handler
func getTopZones(c echo.Context) error {
	order := c.Param("order")
	if order != "dropoffs" && order != "pickups" {
		return c.JSON(http.StatusBadRequest, map[string]string{"parameter missing": "order=dropoffs|pickups"})
	}

	sTop_zones := getTopZonesDb(DATABASE_NAME, order, LOGSQL)
	top_zones := map[string]any{"top_zones": sTop_zones}

	return c.JSON(http.StatusOK, top_zones)
}

// getZoneTrips Handler
func getZoneTrips(c echo.Context) error {
	zoneStr := c.Param("zone")
	zone, err := strconv.Atoi(zoneStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"zone parameter error": err.Error()})
	}
	dateStr := c.Param("date")
	_, err2 := time.Parse("2006-01-02", dateStr) //note: "2006-01-02" is a Golang Format specifier, not a hardcoded value
	if err2 != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"date parameter error": err2.Error()})
	}

	zoneTrips := getZoneTripsDb(DATABASE_NAME, zone, dateStr, LOGSQL)
	return c.JSON(http.StatusOK, zoneTrips)
}

// utiliity function "parseListYellowQuerystring" for getListYellow Handler
func parseListYellowQuerystring(paramstr string) string {
	params := strings.Split(paramstr, "/")

	var (
		sqlQuery    string
		sortSlc     []string
		filterSlc   []string
		limitQuery  string
		offsetQuery string
	)

	operator2sql := map[string]string{"eq": " = ", "gt": " > ", "gte": " >= ", "lt": " < ", "lte": " <= "}

	for _, paramval := range params {
		paramsSlc := strings.Split(paramval, "=")
		key, val := paramsSlc[0], paramsSlc[1]
		switch key {
		case "sort":
			sortField := strings.Split(val, ".")
			sortOrder := strings.ToLower(sortField[1])
			if sortOrder == "asc" || sortOrder == "desc" {
				sortSlc = append(sortSlc, strings.Replace(val, ".", " ", -1))
			}

		case "filter":
			filter := strings.Split(val, ":")
			if len(filter) == 3 {
				field, fieldOperator, fieldValue := filter[0], filter[1], filter[2]
				if strings.Contains(field, "date") {
					fieldValue = "date('" + fieldValue + "')"
				}
				operator, isValidOperator := operator2sql[fieldOperator]
				if isValidOperator {
					filterSql := field + operator + fieldValue
					filterSlc = append(filterSlc, filterSql)
				}
			}
		case "limit":
			limitQuery = " LIMIT " + val
		case "offset":
			offsetQuery = " OFFSET " + val
		}
	}

	if len(filterSlc) > 0 {
		sqlQuery = " AND " + strings.Join(filterSlc, " AND ")
	}
	if len(sortSlc) > 0 {
		sqlQuery += " ORDER BY " + strings.Join(sortSlc, ", ")
	}
	if len(limitQuery) > 0 {
		sqlQuery += limitQuery

		if len(offsetQuery) > 0 { //OFFSET only if LIMIT specified
			sqlQuery += offsetQuery
		}
	} else {
		sqlQuery += " LIMIT 1000 " //max limit: 1000 if not set by query

		if len(offsetQuery) > 0 { //OFFSET only if LIMIT specified
			sqlQuery += offsetQuery
		}
	}
	return sqlQuery
}

// getListYellow Handler
func getListYellow(c echo.Context) error {

	params := c.ParamValues()
	sqlFilter := parseListYellowQuerystring(params[0])

	sListYellow := getListYellowDb(DATABASE_NAME, sqlFilter, LOGSQL)
	list_yellow := map[string]any{"list_yellow": sListYellow}

	return c.JSON(http.StatusOK, list_yellow)
}

// main server function
// reads "server.config.json" params
// configures routes with handlers
// starts server at specified port
// eventually opens local browser with index.html for testing
func main() {
	serverConf := readConfig()
	LOGSQL = serverConf.LogSql
	PORT = serverConf.Port
	OPENLOCALBROWSER = serverConf.OpenLocalBrowser

	e := echo.New()
	e.Static("/", "public")

	e.GET("/api/v1/about", getAbout)
	e.GET("/api/v1/dbversion", getDbVersion)
	e.POST("/api/v1/dbversion", getDbVersion)

	e.GET("/api/v1/top-zones/:order", getTopZones)
	e.POST("/api/v1/top-zones/:order", getTopZones)

	e.GET("/api/v1/zone-trips/:zone/:date", getZoneTrips)
	e.POST("/api/v1/zone-trips/:zone/:date", getZoneTrips)

	e.GET("/api/v1/list-yellow/*", getListYellow)
	e.POST("/api/v1/list-yellow/*", getListYellow)

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	portStr := strconv.Itoa(serverConf.Port)
	fmt.Printf("http://%s:"+portStr, hostname)
	if serverConf.OpenLocalBrowser {
		browser.OpenURL(fmt.Sprintf("http://%s:"+portStr, hostname))
	}
	e.Logger.Fatal(e.Start(":" + portStr))
}
