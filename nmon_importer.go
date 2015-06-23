// nmon2influxdb
// import nmon report in InfluxDB
// author: adejoux@djouxtech.net

package main

import (
	"bufio"
	"fmt"
	"github.com/adejoux/influxdbclient"
	"github.com/codegangsta/cli"
	"os"
	"regexp"
	"strings"
)

var hostRegexp = regexp.MustCompile(`^AAA,host,(\S+)`)
var timeRegexp = regexp.MustCompile(`^ZZZZ,([^,]+),(.*)$`)
var intervalRegexp = regexp.MustCompile(`^AAA,interval,(\d+)`)
var headerRegexp = regexp.MustCompile(`^AAA|^BBB|^UARG|^TOP|,T\d`)
var infoRegexp = regexp.MustCompile(`AAA,(.*)`)
var cpuallRegexp = regexp.MustCompile(`^CPU\d+|^SCPU\d+|^PCPU\d+`)
var diskallRegexp = regexp.MustCompile(`^DISK`)
var skipRegexp = regexp.MustCompile(`T0000|^TOP`)
var statsRegexp = regexp.MustCompile(`[^Z]+,(T\d+)`)

func NmonImport(c *cli.Context) {

	if len(c.Args()) < 1 {
		fmt.Printf("file name needs to be provided\n")
		os.Exit(1)
	}
	// parsing parameters
	params := ParseParameters(c)

	nmon := InitNmon(params)

	influxdb := influxdbclient.NewInfluxDB()
	nmon.Debug = params.Debug
	influxdb.SetDebug(params.Debug)

	influxdb.InitSession(nmon.Params.Host(), params.Db, params.User, params.Password)
	influxdb.Label = nmon.Hostname

	// Hack for influxdb 0.8 API
	// set columns for each serie.
	for serie := range nmon.DataSeries {
		ds := influxdb.DataSeries[serie]
		ds.Columns = nmon.DataSeries[serie].Columns
		influxdb.DataSeries[serie] = ds
	}

	if !influxdb.ExistDB(params.Db) {
		err := influxdb.CreateDB(params.Db)
		check(err)
	}

	file, err := os.Open(params.Filepath)
	check(err)

	defer file.Close()
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if cpuallRegexp.MatchString(scanner.Text()) && !params.CpuAll {
			continue
		}

		if diskallRegexp.MatchString(scanner.Text()) && params.NoDisks {
			continue
		}

		if skipRegexp.MatchString(scanner.Text()) {
			continue
		}

		if statsRegexp.MatchString(scanner.Text()) {
			matched := statsRegexp.FindStringSubmatch(scanner.Text())
			elems := strings.Split(scanner.Text(), ",")
			timeStr, err := nmon.GetTimeStamp(matched[1])
			check(err)
			name := elems[0]
			timestamp, err := ConvertTimeStamp(timeStr, nmon.Params.TZ)
			influxdb.AddPoint(name, timestamp, elems[2:])

			if influxdb.MaxPointsCount(name) {
				err = influxdb.WritePoints(name)
				check(err)
				influxdb.ClearPoints(name)
			}
		}
	}
	// flushing remaining data
	for serie := range influxdb.DataSeries {
		influxdb.WritePoints(serie)
	}

	fmt.Printf("File %s imported !\n", params.Filepath)
}
