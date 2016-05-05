// nmon2influxdb
// import nmon data in InfluxDB
// author = adejoux@djouxtech.net

package main

import (
	"log"
	"strings"
)

//
//helper functions
//
func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func checkInfo(e error) {
	if e != nil {
		log.Printf("info: %s", e)
	}
}

//ReplaceComma replaces comma by html tabs tag
func ReplaceComma(s string) string {
	return "<tr><td>" + strings.Replace(s, ",", "</td><td>", 1) + "</td></tr>"
}
