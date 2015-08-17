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
	params := ParseParameters(c)

	influxdb := influxdbclient.NewInfluxDB(params.Server, params.Port, params.Db, params.User, params.Password)
	influxdb.SetDebug(params.Debug)
	influxdb.Connect()

	filters := new(influxdbclient.Filters)

	if len(params.Host) > 0 {
		filters.Add("host", params.Host, "text")
	}

	measurements, _ := influxdb.ListMeasurement(filters)
	if measurements != nil {
		fmt.Printf("%s\n", measurements.Name)
		for _, value := range measurements.Datas {
			if len(params.Filter) == 0 {
				fmt.Printf("%s\n", value)
				continue
			}
			matched, _ := regexp.MatchString(params.Filter, value)
			if matched {
				fmt.Printf("%s\n", value)
			}
		}
	}
}
