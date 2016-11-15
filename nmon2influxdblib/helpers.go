// nmon2influxdb
// import nmon data in InfluxDB
// author = adejoux@djouxtech.net

package nmon2influxdblib

import (
	"log"
	"strings"
)

//
//helper functions

//CheckError check error message and display it
func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// CheckInfo wrap info message
func CheckInfo(e error) {
	if e != nil {
		log.Printf("info: %s", e)
	}
}

//ReplaceComma replaces comma by html tabs tag
func ReplaceComma(s string) string {
	return "<tr><td>" + strings.Replace(s, ",", "</td><td>", 1) + "</td></tr>"
}
