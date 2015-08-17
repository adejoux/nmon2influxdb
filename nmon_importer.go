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
	"sort"
	"strconv"
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
var nfsRegexp = regexp.MustCompile(`^NFS`)
var nameRegexp = regexp.MustCompile(`(\d+)$`)

func NmonImport(c *cli.Context) {

	if len(c.Args()) < 1 {
		fmt.Printf("file name needs to be provided\n")
		os.Exit(1)
	}
	// parsing parameters
	params := ParseParameters(c)

	nmon := InitNmon(params)

	influxdb := influxdbclient.NewInfluxDB(params.Server, params.Port, params.Db, params.User, params.Password)
	nmon.Debug = params.Debug
	influxdb.SetDebug(params.Debug)

	influxdb.Connect()

	if exist, _ := influxdb.ExistDB(params.Db); exist != true {
		_, err := influxdb.CreateDB(params.Db)
		check(err)
	}

	file, err := os.Open(params.Name)
	check(err)

	defer file.Close()
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	sort.Strings(lines)

	for _, line := range lines {

		if cpuallRegexp.MatchString(line) && !params.CpuAll {
			continue
		}

		if diskallRegexp.MatchString(line) && params.NoDisks {
			continue
		}

		if skipRegexp.MatchString(line) {
			continue
		}

		if statsRegexp.MatchString(line) {
			matched := statsRegexp.FindStringSubmatch(line)
			elems := strings.Split(line, ",")
			timeStr, err := nmon.GetTimeStamp(matched[1])
			check(err)
			name := elems[0]
			timestamp, err := ConvertTimeStamp(timeStr, nmon.Params.TZ)

			for i, value := range elems[2:] {
				tags := map[string]string{"host": nmon.Hostname, "name": nmon.DataSeries[name].Columns[i]}

				// try to convert string to integer
				converted, err := strconv.ParseFloat(value, 64)
				if err != nil {
					//if not working, skip to next value. We don't want text values in InfluxDB.
					continue
				}

				//send integer if it worked
				field := map[string]interface{}{"value": converted}

				measurement := ""
				if nfsRegexp.MatchString(name) {
					measurement = name
				} else {
					measurement = nameRegexp.ReplaceAllString(name, "")
				}

				influxdb.AddPoint(measurement, timestamp, field, tags)

				if influxdb.PointsCount() == 10000 {
					err = influxdb.WritePoints()
					check(err)
					influxdb.ClearPoints()
					fmt.Printf("#")
				}
			}
		}
	}
	// flushing remaining data
	influxdb.WritePoints()
	fmt.Printf("#\n")

	fmt.Printf("File %s imported !\n", params.Name)
}
