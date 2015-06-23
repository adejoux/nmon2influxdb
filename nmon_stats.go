// nmon2influxdb
// import nmon report in InfluxDB
// author: adejoux@djouxtech.net
package main

import (
	"fmt"
	"github.com/adejoux/influxdbclient"
	"github.com/codegangsta/cli"
	"time"
)

const querytimeformat = "2006-01-02 15:04:05"

func NmonStat(c *cli.Context) {
	// parsing parameters
	params := ParseStatsParameters(c)

	influxdb := influxdbclient.NewInfluxDB()
	influxdb.SetDebug(params.Debug)

	influxdb.InitSession(params.Host(), params.Db, params.User, params.Password)

	metric := params.StatsHost + "_" + params.Metric

	fromUnix, _ := ConvertTimeStamp(params.From, params.TZ)
	fromTime := time.Unix(fromUnix, 0).UTC().Format(querytimeformat)
	toUnix, _ := ConvertTimeStamp(params.To, params.TZ)
	toTime := time.Unix(toUnix, 0).UTC().Format(querytimeformat)
	result, err := influxdb.ReadPoints("*", metric, fromTime, toTime, "")
	if err != nil {
		check(err)
	}

	//generate stats
	stats := influxdb.BuildStats(result)

	DisplayStats(&stats, params.Sort, params.Limit)
}

func DisplayStats(stats *influxdbclient.DataStats, sort string, limit int) {
	fmt.Printf("%20s|%10s|%10s|%10s|%10s|%10s\n", "field", "Min", "Mean", "Median", "Max", "Points #")
	stats.FieldSort(sort)
	for i, stat := range *stats {
		fmt.Printf("%20s|%10.2f|%10.2f|%10.2f|%10.2f|%8d\n", stat.Name, stat.Min, stat.Mean, stat.Median, stat.Max, stat.Length)
		if limit > 0 {
			if i > limit {
				break
			}
		}
	}
}
