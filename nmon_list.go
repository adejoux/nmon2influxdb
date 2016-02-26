// nmon2influxdb
// import nmon report in InfluxDB
// author: adejoux@djouxtech.net
package main

import (
	"fmt"
	"github.com/adejoux/influxdbclient"
	"github.com/codegangsta/cli"
	"regexp"
	//	"os"
)

func NmonListMeasurement(c *cli.Context) {
	// parsing parameters
	config := ParseParameters(c)

	influxdb := influxdbclient.NewInfluxDB(config.InfluxdbServer, config.InfluxdbPort, config.InfluxdbDatabase, config.InfluxdbUser, config.InfluxdbPassword)
	influxdb.SetDebug(config.Debug)
	err := influxdb.Connect()
	check(err)

	filters := new(influxdbclient.Filters)

	if len(config.ListHost) > 0 {
		filters.Add("host", config.ListHost, "text")
	}

	measurements, _ := influxdb.ListMeasurement(filters)
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
