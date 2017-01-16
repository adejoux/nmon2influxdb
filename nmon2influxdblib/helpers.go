// nmon2influxdb
// import nmon data in InfluxDB
// author = adejoux@djouxtech.net

package nmon2influxdblib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
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

//SPrintHTTPResponse print raw http response for debugging purpose
func SPrintHTTPResponse(response *http.Response) string {
	responseDump, err := httputil.DumpResponse(response, true)
	if err != nil {
		fmt.Println(err)
	}
	return string(responseDump)
}

//SPrintHTTPRequest print raw http request for debugging purpose
func SPrintHTTPRequest(request *http.Request) string {
	requestDump, err := httputil.DumpRequest(request, true)
	if err != nil {
		log.Println(err)
	}
	return string(requestDump)
}

//SPrintPrettyJSON helper used to display JSON output in a nicer way
func SPrintPrettyJSON(contents []byte) string {
	text := GetPrettyJSON(contents)
	return string(text.Bytes())
}

//GetPrettyJSON returns pretty json string
func GetPrettyJSON(contents []byte) bytes.Buffer {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, contents, "", "\t")
	if error != nil {
		log.Println("JSON parse error: ", error)
	}

	return prettyJSON
}
