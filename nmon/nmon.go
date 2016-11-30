// nmon2influxdb
// import nmon data in InfluxDB
// author: adejoux@djouxtech.net

package nmon

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/adejoux/nmon2influxdb/nmon2influxdblib"
)

// Nmon structure used to manage nmon files
type Nmon struct {
	Hostname    string
	Serial      string
	OS          string
	TimeStamps  map[string]string
	TextContent string
	DataSeries  map[string]DataSerie
	Debug       bool
	Config      *nmon2influxdblib.Config
	starttime   time.Time
	stoptime    time.Time
	Location    *time.Location
	TagParsers  nmon2influxdblib.TagParsers
}

// DataSerie structure contains the columns and points to insert in InfluxDB
type DataSerie struct {
	Columns []string
}

// AppendText add text section to dashboard
func (nmon *Nmon) AppendText(text string) {
	nmon.TextContent += nmon2influxdblib.ReplaceComma(text)
}

// NewNmon initialize a Nmon structure
func NewNmon() *Nmon {
	return &Nmon{DataSeries: make(map[string]DataSerie), TimeStamps: make(map[string]string)}

}

// BuildPoint create a point and convert string value to float when possible
func (nmon *Nmon) BuildPoint(serie string, values []string) map[string]interface{} {
	columns := nmon.DataSeries[serie].Columns
	//TODO check output
	point := make(map[string]interface{})

	for i, rawvalue := range values {
		// try to convert string to integer
		value, err := strconv.ParseFloat(rawvalue, 64)
		if err != nil {
			//if not working, use string
			point[columns[i]] = rawvalue
		} else {
			//send integer if it worked
			point[columns[i]] = value
		}
	}

	return point
}

//GetTimeStamp retrieves the TimeStamp corresponding to the entry
func (nmon *Nmon) GetTimeStamp(label string) (timeStamp string, err error) {
	if t, ok := nmon.TimeStamps[label]; ok {
		timeStamp = t
	} else {
		errorMessage := fmt.Sprintf("TimeStamp %s not found", label)
		err = errors.New(errorMessage)
	}
	return
}

//InitNmonTemplate init nmon structure when creating dashboard
func InitNmonTemplate(config *nmon2influxdblib.Config) (nmon *Nmon) {
	nmon = NewNmon()
	nmon.Config = config
	nmon.SetLocation(config.Timezone)
	return
}

//InitNmon init nmon structure for nmon file import
func InitNmon(config *nmon2influxdblib.Config, nmonFile nmon2influxdblib.File) (nmon *Nmon) {
	nmon = NewNmon()
	nmon.Config = config
	nmon.SetLocation(config.Timezone)
	nmon.Debug = config.Debug

	var lines []string
	if len(nmonFile.Host) > 0 {
		scanner, err := nmonFile.GetRemoteScanner()
		nmon2influxdblib.CheckError(err)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		scanner.Close()
	} else {
		scanner, err := nmonFile.GetScanner()
		nmon2influxdblib.CheckError(err)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		scanner.Close()
	}

	var userSkipRegexp *regexp.Regexp

	if len(config.ImportSkipMetrics) > 0 {
		skipped := strings.Replace(config.ImportSkipMetrics, ",", "|", -1)
		userSkipRegexp = regexp.MustCompile(skipped)
	}

	for _, line := range lines {

		if cpuallRegexp.MatchString(line) && !config.ImportAllCpus {
			continue
		}

		if diskallRegexp.MatchString(line) && config.ImportSkipDisks {
			continue
		}

		if timeRegexp.MatchString(line) {
			matched := timeRegexp.FindStringSubmatch(line)
			nmon.TimeStamps[matched[1]] = matched[2]
			continue
		}

		if hostRegexp.MatchString(line) {
			matched := hostRegexp.FindStringSubmatch(line)
			nmon.Hostname = strings.ToLower(matched[1])
			continue
		}

		if serialRegexp.MatchString(line) {
			matched := serialRegexp.FindStringSubmatch(line)
			nmon.Serial = strings.ToLower(matched[1])
			continue
		}

		if osRegexp.MatchString(line) {
			matched := osRegexp.FindStringSubmatch(line)
			nmon.OS = strings.ToLower(matched[1])
			continue
		}

		if infoRegexp.MatchString(line) {
			matched := infoRegexp.FindStringSubmatch(line)
			nmon.AppendText(matched[1])
			continue
		}

		if !headerRegexp.MatchString(line) {
			if len(line) == 0 {
				continue
			}

			if badRegexp.MatchString(line) {
				continue
			}
			elems := strings.Split(line, ",")

			if len(elems) < 3 {
				if config.Debug == true {
					fmt.Printf("ERROR: parsing the following line : %s\n", line)
				}
				continue
			}
			name := elems[0]
			if len(config.ImportSkipMetrics) > 0 {
				if userSkipRegexp.MatchString(name) {
					continue
				}
			}

			if config.Debug == true {
				fmt.Printf("ADDING serie %s\n", name)
			}

			dataserie := nmon.DataSeries[name]
			dataserie.Columns = elems[2:]
			nmon.DataSeries[name] = dataserie
		}
	}
	return
}

//SetTimeFrame set the current timeframe for the dashboard
func (nmon *Nmon) SetTimeFrame() {
	keys := make([]string, 0, len(nmon.TimeStamps))
	for k := range nmon.TimeStamps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	nmon.starttime, _ = nmon.ConvertTimeStamp(nmon.TimeStamps[keys[0]])
	nmon.stoptime, _ = nmon.ConvertTimeStamp(nmon.TimeStamps[keys[len(keys)-1]])
}

// StartTime returns the starting timestamp for dashboard
func (nmon *Nmon) StartTime() string {
	if nmon.starttime == (time.Time{}) {
		nmon.SetTimeFrame()
	}
	return nmon.starttime.UTC().Format(time.RFC3339)
}

// StopTime returns the ending timestamp for dashboard
func (nmon *Nmon) StopTime() string {
	if nmon.stoptime == (time.Time{}) {
		nmon.SetTimeFrame()
	}
	return nmon.stoptime.UTC().Format(time.RFC3339)
}

const timeformat = "15:04:05,02-Jan-2006"

//SetLocation set the timezone used to input metrics in InfluxDB
func (nmon *Nmon) SetLocation(tz string) (err error) {
	var loc *time.Location
	if len(tz) > 0 {
		loc, err = time.LoadLocation(tz)
		if err != nil {
			loc = time.FixedZone("Europe/Paris", 2*60*60)
		}
	} else {
		timezone, _ := time.Now().In(time.Local).Zone()
		loc, err = time.LoadLocation(timezone)
		if err != nil {
			loc = time.FixedZone("Europe/Paris", 2*60*60)
		}
	}

	nmon.Location = loc
	return
}

//ConvertTimeStamp convert the string timestamp in time.Time structure
func (nmon *Nmon) ConvertTimeStamp(s string) (time.Time, error) {
	var err error
	if s == "now" {
		return time.Now().Truncate(24 * time.Hour), err
	}

	t, err := time.ParseInLocation(timeformat, s, nmon.Location)
	return t, err
}

//DbURL generates InfluxDB server url
func (nmon *Nmon) DbURL() string {
	return "http://" + nmon.Config.InfluxdbServer + ":" + nmon.Config.InfluxdbPort
}
