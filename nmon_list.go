// nmon2influxdb
// import nmon report in InfluxDB
// author: adejoux@djouxtech.net
package main

import (
	"fmt"
	"regexp"

	"github.com/adejoux/influxdbclient"
	"github.com/codegangsta/cli"
	//	"os"
)

//NmonListMeasurement list all measurements in INFLUXDB database
func NmonListMeasurement(c *cli.Context) {
	// parsing parameters
	config := ParseParameters(c)

	influxdb := config.connectDB(config.InfluxdbDatabase)
	filters := new(influxdbclient.Filters)

	if len(config.ListHost) > 0 {
		filters.Add("host", config.ListHost, "text")
	}

	measurements, err := influxdb.ListMeasurement(filters)
	check(err)
	if measurements != nil {
		fmt.Printf("%s\n", measurements.Name)
		for _, value := range measurements.Datas {
			if len(config.ListFilter) == 0 {
				fmt.Printf("%s\n", value)
				continue
			}
			matched, _ := regexp.MatchString(config.ListFilter, value)
			if matched {
				fmt.Printf("%s\n", value)
			}
		}
	}
}
