// nmon2influxdb
// import nmon data in InfluxDB
// author: adejoux@djouxtech.net

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Nmon struct {
	Hostname    string
	TimeStamps  map[string]string
	TextContent string
	DataSeries  map[string]DataSerie
	Debug       bool
	Params      *Params
	starttime   time.Time
	stoptime    time.Time
}

//
// DataSerie structure
// contains the columns and points to insert in InfluxDB
//

type DataSerie struct {
	Columns []string
}

func (nmon *Nmon) AppendText(text string) {
	nmon.TextContent += ReplaceComma(text)
}

// initialize a Nmon structure
func NewNmon() *Nmon {
	return &Nmon{DataSeries: make(map[string]DataSerie), TimeStamps: make(map[string]string)}

}

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

func (nmon *Nmon) GetTimeStamp(label string) (t string, err error) {
	if t, ok := nmon.TimeStamps[label]; ok {
		return t, err
	} else {
		error_message := fmt.Sprintf("TimeStamp %s not found", label)
		err = errors.New(error_message)
	}
	return t, err
}

func InitNmonTemplate(params *Params) (nmon *Nmon) {
	nmon = NewNmon()
	nmon.Params = params
	return
}

func InitNmon(params *Params) (nmon *Nmon) {
	nmon = NewNmon()
	nmon.Params = params
	file, err := os.Open(params.Name)
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

		if timeRegexp.MatchString(scanner.Text()) {
			matched := timeRegexp.FindStringSubmatch(scanner.Text())
			nmon.TimeStamps[matched[1]] = matched[2]
			continue
		}

		if hostRegexp.MatchString(scanner.Text()) {
			matched := hostRegexp.FindStringSubmatch(scanner.Text())
			nmon.Hostname = strings.ToLower(matched[1])
			continue
		}

		if infoRegexp.MatchString(scanner.Text()) {
			matched := infoRegexp.FindStringSubmatch(scanner.Text())
			nmon.AppendText(matched[1])
			continue
		}

		if !headerRegexp.MatchString(scanner.Text()) {
			if len(scanner.Text()) == 0 {
				continue
			}

			elems := strings.Split(scanner.Text(), ",")

			if len(elems) < 3 {
				if params.Debug == true {
					fmt.Printf("ERROR: parsing the following line : %s\n", scanner.Text())
				}
				continue
			}
			name := elems[0]
			if params.Debug == true {
				fmt.Printf("ADDING serie %s\n", name)
			}

			dataserie := nmon.DataSeries[name]
			dataserie.Columns = elems[2:]
			nmon.DataSeries[name] = dataserie
		}
	}
	return
}

func (nmon *Nmon) SetTimeFrame() {
	keys := make([]string, 0, len(nmon.TimeStamps))
	for k := range nmon.TimeStamps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	nmon.starttime, _ = ConvertTimeStamp(nmon.TimeStamps[keys[0]], nmon.Params.TZ)
	nmon.stoptime, _ = ConvertTimeStamp(nmon.TimeStamps[keys[len(keys)-1]], nmon.Params.TZ)
}

func (nmon *Nmon) StartTime() string {
	if nmon.starttime == (time.Time{}) {
		nmon.SetTimeFrame()
	}
	return nmon.starttime.UTC().Format(time.RFC3339)
}

func (nmon *Nmon) StopTime() string {
	if nmon.stoptime == (time.Time{}) {
		nmon.SetTimeFrame()
	}
	return nmon.stoptime.UTC().Format(time.RFC3339)
}

const timeformat = "15:04:05,02-Jan-2006"

func ConvertTimeStamp(s string, tz string) (time.Time, error) {
	var err error
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

	t, err := time.ParseInLocation(timeformat, s, loc)
	return t, err
}

func (nmon *Nmon) DataSource() string {
	return nmon.Params.DS
}

func (nmon *Nmon) DbURL() string {
	return "http://" + nmon.Params.Server + ":" + nmon.Params.Port
}
