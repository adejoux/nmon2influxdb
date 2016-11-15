// nmon2influxdb
// import nmon report in InfluxDB
// author: adejoux@djouxtech.net
package nmon

import (
	"fmt"
	"regexp"

	"github.com/adejoux/influxdbclient"
	"github.com/adejoux/nmon2influxdb/nmon2influxdblib"
	"github.com/codegangsta/cli"
	//	"os"
)

//ListMeasurement list all measurements in INFLUXDB database
func ListMeasurement(c *cli.Context) {
	// parsing parameters
	config := nmon2influxdblib.ParseParameters(c)

	influxdb := config.ConnectDB(config.InfluxdbDatabase)
	filters := new(influxdbclient.Filters)

	if len(config.ListHost) > 0 {
		filters.Add("host", config.ListHost, "text")
	}

	measurements, err := influxdb.ListMeasurement(filters)
	nmon2influxdblib.CheckError(err)
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
