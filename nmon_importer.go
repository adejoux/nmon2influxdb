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
var osRegexp = regexp.MustCompile(`^AAA,.*(Linux|AIX)`)
var timeRegexp = regexp.MustCompile(`^ZZZZ,([^,]+),(.*)$`)
var intervalRegexp = regexp.MustCompile(`^AAA,interval,(\d+)`)
var headerRegexp = regexp.MustCompile(`^AAA|^BBB|^UARG|,T\d`)
var infoRegexp = regexp.MustCompile(`AAA,(.*)`)
var badRegexp = regexp.MustCompile(`,,`)
var cpuallRegexp = regexp.MustCompile(`^CPU\d+|^SCPU\d+|^PCPU\d+`)
var diskallRegexp = regexp.MustCompile(`^DISK`)
var skipRegexp = regexp.MustCompile(`T0000|^Z|^TOP,%CPU`)
var statsRegexp = regexp.MustCompile(`^[^,]+?,(T\d+)`)
var topRegexp = regexp.MustCompile(`^TOP,\d+,(T\d+)`)
var nfsRegexp = regexp.MustCompile(`^NFS`)
var nameRegexp = regexp.MustCompile(`(\d+)$`)

func NmonImport(c *cli.Context) {

	if len(c.Args()) < 1 {
		fmt.Printf("file name needs to be provided\n")
		os.Exit(1)
	}

	// parsing parameters
	params := ParseParameters(c)

	influxdb := influxdbclient.NewInfluxDB(params.Server, params.Port, params.Db, params.User, params.Password)
	influxdb.SetDebug(params.Debug)
	influxdb.Connect()
	if exist, _ := influxdb.ExistDB(params.Db); exist != true {
		_, err := influxdb.CreateDB(params.Db)
		check(err)
	}

	for _, nmon_file := range c.Args() {

		nmon := InitNmon(params, nmon_file)
		file, err := os.Open(nmon_file)
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
				name := elems[0]

				timeStr, err := nmon.GetTimeStamp(matched[1])
				check(err)

				timestamp, err := ConvertTimeStamp(timeStr, nmon.Params.TZ)

				for i, value := range elems[2:] {
					if len(nmon.DataSeries[name].Columns) < i+1 {
						if nmon.Debug {
							fmt.Printf(line)
							fmt.Printf("Entry added position %d in serie %s since nmon start: skipped\n", i+1, name)
						}
						continue
					}
					column := nmon.DataSeries[name].Columns[i]
					tags := map[string]string{"host": nmon.Hostname, "name": column}

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

			if topRegexp.MatchString(line) {
				matched := topRegexp.FindStringSubmatch(line)
				elems := strings.Split(line, ",")
				timeStr, err := nmon.GetTimeStamp(matched[1])
				check(err)
				timestamp, err := ConvertTimeStamp(timeStr, nmon.Params.TZ)

				if len(elems) < 14 {
					fmt.Printf("error TOP import:")
					fmt.Println(elems)
					continue
				}

				for i, value := range elems[3:12] {
					column := nmon.DataSeries["TOP"].Columns[i]

					var wlmclass string
					if len(elems) < 15 {
						wlmclass = "none"
					} else {
						wlmclass = elems[14]
					}

					tags := map[string]string{"host": nmon.Hostname, "name": column, "pid": elems[1], "command": elems[13], "wlm": wlmclass}

					// try to convert string to integer
					converted, err := strconv.ParseFloat(value, 64)
					if err != nil {
						//if not working, skip to next value. We don't want text values in InfluxDB.
						continue
					}

					//send integer if it worked
					field := map[string]interface{}{"value": converted}

					influxdb.AddPoint("TOP", timestamp, field, tags)

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

		fmt.Printf("File %s imported !\n", nmon_file)
	}
}
