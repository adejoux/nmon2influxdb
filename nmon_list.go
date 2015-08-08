// nmon2influxdb
// import nmon report in InfluxDB
// author: adejoux@djouxtech.net
package main

import (
	"github.com/adejoux/influxdbclient"
	"github.com/codegangsta/cli"
)

func NmonList(c *cli.Context) {
	// parsing parameters
	params := ParseParameters(c)

	influxdb := influxdbclient.NewInfluxDB(params.Server, params.Port, params.Db, params.User, params.Password)
	influxdb.SetDebug(params.Debug)
	influxdb.Connect()
}
