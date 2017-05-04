// nmon2influxdb
// import nmon report in InfluxDB
// author: adejoux@djouxtech.net

package nmon

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"math"

	"github.com/adejoux/influxdbclient"
	"github.com/adejoux/nmon2influxdb/nmon2influxdblib"
	"github.com/codegangsta/cli"
)

var hostRegexp = regexp.MustCompile(`^AAA.host.(\S+)`)
var serialRegexp = regexp.MustCompile(`^AAA.SerialNumber.(\S+)`)
var osRegexp = regexp.MustCompile(`^AAA.*(Linux|AIX)`)
var timeRegexp = regexp.MustCompile(`^ZZZZ.(T\d+).(.*)$`)
var intervalRegexp = regexp.MustCompile(`^AAA.interval.(\d+)`)
var headerRegexp = regexp.MustCompile(`^AAA|^BBB|^UARG|\WT\d{4,16}`)
var infoRegexp = regexp.MustCompile(`^AAA.(.*)`)
var cpuallRegexp = regexp.MustCompile(`^CPU\d+|^SCPU\d+|^PCPU\d+`)
var diskallRegexp = regexp.MustCompile(`^DISK`)
var skipRegexp = regexp.MustCompile(`T0+\W|^Z|^TOP.%CPU`)
var statsRegexp = regexp.MustCompile(`\W(T\d{4,16})`)
var topRegexp = regexp.MustCompile(`^TOP.\d+.(T\d+)`)
var nfsRegexp = regexp.MustCompile(`^NFS`)
var nameRegexp = regexp.MustCompile(`(\d+)$`)

//Import is the entry point for subcommand nmon import
func Import(c *cli.Context) {

	if len(c.Args()) < 1 {
		fmt.Printf("file name or directory needs to be provided\n")
		os.Exit(1)
	}

	// parsing parameters
	config := nmon2influxdblib.ParseParameters(c)

	//getting databases connections
	influxdb := config.GetDB("nmon")
	influxdbLog := config.GetLogDB()

	nmonFiles := new(nmon2influxdblib.Files)
	nmonFiles.Parse(c.Args(), config.ImportSSHUser, config.ImportSSHKey)

	tagParsers := nmon2influxdblib.ParseInputs(config.Inputs)

	var userSkipRegexp *regexp.Regexp
	if len(config.ImportSkipMetrics) > 0 {
		skipped := strings.Replace(config.ImportSkipMetrics, ",", "|", -1)
		userSkipRegexp = regexp.MustCompile(skipped)
	}

	for _, nmonFile := range nmonFiles.Valid() {


		// store the list of metrics which was logged as skipped
		LoggedSkippedMetrics := make(map[string]bool)
		var count int64
		count = 0
		nmon := InitNmon(config, nmonFile)

		if len(config.Inputs) > 0 {
			//Build tag parsing
			nmon.TagParsers = tagParsers
		}

		lines := nmonFile.Content()
		log.Printf("NMON file separator: %s\n", nmonFile.Delimiter)
		var last string
		filters := new(influxdbclient.Filters)
		filters.Add("file", path.Base(nmonFile.Name), "text")

		result, err := influxdbLog.ReadLastPoint("value", filters, "timestamp")
		nmon2influxdblib.CheckError(err)

		var lastTime time.Time
		if !nmon.Config.ImportForce && len(result) > 0 {
			lastTime, err = nmon.ConvertTimeStamp(result[1].(string))
		} else {
			lastTime, err = nmon.ConvertTimeStamp("00:00:00,01-JAN-1900")
		}
		nmon2influxdblib.CheckError(err)

		origChecksum, err := influxdbLog.ReadLastPoint("value", filters, "checksum")
		nmon2influxdblib.CheckError(err)

		ckfield := map[string]interface{}{"value": nmonFile.Checksum()}
		if !nmon.Config.ImportForce && len(origChecksum) > 0 {

			if origChecksum[1].(string) == nmonFile.Checksum() {
				fmt.Printf("file not changed since last import: %s\n", nmonFile.Name)
				continue
			}
		}

		for _, line := range lines {

			if cpuallRegexp.MatchString(line) && !config.ImportAllCpus {
				continue
			}

			if diskallRegexp.MatchString(line) && config.ImportSkipDisks {
				continue
			}

			if skipRegexp.MatchString(line) {
				continue
			}

			if statsRegexp.MatchString(line) {
				matched := statsRegexp.FindStringSubmatch(line)
				elems := strings.Split(line, nmonFile.Delimiter)
				name := elems[0]

				if len(config.ImportSkipMetrics) > 0 {
					if userSkipRegexp.MatchString(name) {
						if nmon.Debug {
							if !LoggedSkippedMetrics[name] {
								log.Printf("metric skipped : %s\n", name)
								LoggedSkippedMetrics[name] = true
							}
						}
						continue
					}
				}

				timeStr, getErr := nmon.GetTimeStamp(matched[1])
				if getErr != nil {
					continue
				}
				last = timeStr
				timestamp, convErr := nmon.ConvertTimeStamp(timeStr)
				nmon2influxdblib.CheckError(convErr)
				if timestamp.Before(lastTime) && !nmon.Config.ImportForce {
					continue
				}

				for i, value := range elems[2:] {
					if len(nmon.DataSeries[name].Columns) < i+1 {
						if nmon.Debug {
							log.Printf(line)
							log.Printf("Entry added position %d in serie %s since nmon start: skipped\n", i+1, name)
						}
						continue
					}
					column := nmon.DataSeries[name].Columns[i]
					tags := map[string]string{"host": nmon.Hostname, "name": column}

					// try to convert string to integer
					converted, parseErr := strconv.ParseFloat(value, 64)
					if (parseErr != nil || math.IsNaN(converted)) {
						//if not working, skip to next value. We don't want text values in InfluxDB.
						continue
					}

					//send integer if it worked
					field := map[string]interface{}{"value": converted}

					measurement := ""
					if nfsRegexp.MatchString(name) || cpuallRegexp.MatchString(name) {
						measurement = name
					} else {
						measurement = nameRegexp.ReplaceAllString(name, "")
					}

					// Checking additional tagging
					for key, value := range tags {
						if _, ok := nmon.TagParsers[measurement][key]; ok {
							for _, tagParser := range nmon.TagParsers[measurement][key] {
								if tagParser.Regexp.MatchString(value) {
									tags[tagParser.Name] = tagParser.Value
								}
							}
						}
					}
					influxdb.AddPoint(measurement, timestamp, field, tags)

					if influxdb.PointsCount() >= 5000 {
						err = influxdb.WritePoints()
						nmon2influxdblib.CheckError(err)
						count += influxdb.PointsCount()
						influxdb.ClearPoints()
						fmt.Printf("#")
					}
				}
			}

			if topRegexp.MatchString(line) {
				matched := topRegexp.FindStringSubmatch(line)

				elems := strings.Split(line, nmonFile.Delimiter)
				name := elems[0]
				if len(config.ImportSkipMetrics) > 0 {
					if userSkipRegexp.MatchString(name) {
						if nmon.Debug {
							if !LoggedSkippedMetrics[name] {
								log.Printf("metric skipped : %s\n", name)
								LoggedSkippedMetrics[name] = true
							}
						}
						continue
					}
				}

				timeStr, getErr := nmon.GetTimeStamp(matched[1])
				if getErr != nil {
					continue
				}

				timestamp, convErr := nmon.ConvertTimeStamp(timeStr)
				nmon2influxdblib.CheckError(convErr)

				if len(elems) < 14 {
					log.Printf("error TOP import:")
					log.Println(elems)
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

					if len(nmon.Serial) > 0 {
						tags["serial"] = nmon.Serial
					}

					// try to convert string to integer
					converted, parseErr := strconv.ParseFloat(value, 64)
					if parseErr != nil {
						//if not working, skip to next value. We don't want text values in InfluxDB.
						continue
					}

					//send integer if it worked
					field := map[string]interface{}{"value": converted}

					influxdb.AddPoint("TOP", timestamp, field, tags)

					if influxdb.PointsCount() == 10000 {
						err = influxdb.WritePoints()
						nmon2influxdblib.CheckError(err)
						count += influxdb.PointsCount()
						influxdb.ClearPoints()
						fmt.Printf("#")
					}
				}
			}
		}
		// flushing remaining data
		influxdb.WritePoints()
		count += influxdb.PointsCount()
		fmt.Printf("\nFile %s imported : %d points !\n", nmonFile.Name, count)
		if config.ImportBuildDashboard {
			DashboardFile(config, nmonFile.Name)
		}

		if len(last) > 0 {
			field := map[string]interface{}{"value": last}
			tag := map[string]string{"file": path.Base(nmonFile.Name)}
			lasttime, _ := nmon.ConvertTimeStamp("now")
			influxdbLog.AddPoint("timestamp", lasttime, field, tag)
			influxdbLog.AddPoint("checksum", lasttime, ckfield, tag)
			err = influxdbLog.WritePoints()
			nmon2influxdblib.CheckError(err)
		}
	}
}
