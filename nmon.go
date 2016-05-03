// nmon2influxdb
// import nmon data in InfluxDB
// author: adejoux@djouxtech.net

package main

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const gzipfile = ".gz"

// Nmon structure used to manage nmon files
type Nmon struct {
	Hostname    string
	OS          string
	TimeStamps  map[string]string
	TextContent string
	DataSeries  map[string]DataSerie
	Debug       bool
	Config      *Config
	starttime   time.Time
	stoptime    time.Time
	Location    *time.Location
}

// DataSerie structure contains the columns and points to insert in InfluxDB
type DataSerie struct {
	Columns []string
}

// AppendText add text section to dashboard
func (nmon *Nmon) AppendText(text string) {
	nmon.TextContent += ReplaceComma(text)
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
func InitNmonTemplate(config *Config) (nmon *Nmon) {
	nmon = NewNmon()
	nmon.Config = config
	nmon.SetLocation(config.Timezone)
	return
}

//InitNmon init nmon structure for nmon file import
func InitNmon(config *Config, nmonFile NmonFile) (nmon *Nmon) {
	nmon = NewNmon()
	nmon.Config = config
	nmon.SetLocation(config.Timezone)
	nmon.Debug = config.Debug

	scanner, err := nmonFile.GetScanner()
	check(err)
	defer scanner.Close()

	scanner.Split(bufio.ScanLines)

	var userSkipRegexp *regexp.Regexp

	if len(config.ImportSkipMetrics) > 0 {
		skipped := strings.Replace(config.ImportSkipMetrics, ",", "|", -1)
		userSkipRegexp = regexp.MustCompile(skipped)
	}

	for scanner.Scan() {

		if cpuallRegexp.MatchString(scanner.Text()) && !config.ImportAllCpus {
			continue
		}

		if diskallRegexp.MatchString(scanner.Text()) && config.ImportSkipDisks {
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

		if osRegexp.MatchString(scanner.Text()) {
			matched := osRegexp.FindStringSubmatch(scanner.Text())
			nmon.OS = strings.ToLower(matched[1])
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

			if badRegexp.MatchString(scanner.Text()) {
				continue
			}
			elems := strings.Split(scanner.Text(), ",")

			if len(elems) < 3 {
				if config.Debug == true {
					fmt.Printf("ERROR: parsing the following line : %s\n", scanner.Text())
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

// NmonFile structure used to select nmon files to import
type NmonFile struct {
	Name     string
	FileType string
}

// NmonFiles array of NmonFile
type NmonFiles []NmonFile

//Add a file in the NmonFIles structure
func (nmonFiles *NmonFiles) Add(file string, fileType string) {
	*nmonFiles = append(*nmonFiles, NmonFile{Name: file, FileType: fileType})
}

//Valid returns only valid fiels for nmon import
func (nmonFiles *NmonFiles) Valid() (validFiles NmonFiles) {
	for _, v := range *nmonFiles {
		if v.FileType == ".nmon" || v.FileType == gzipfile {
			validFiles = append(validFiles, v)
		}
	}
	return validFiles
}

// fileScanner struct to manage
type fileScanner struct {
	*os.File
	*bufio.Scanner
}

// GetScanner open an nmon file based on file extension and provides a bufio Scanner
func (nmonFile *NmonFile) GetScanner() (*fileScanner, error) {
	file, err := os.Open(nmonFile.Name)
	if err != nil {
		return nil, err
	}

	if nmonFile.FileType == gzipfile {
		gr, err := gzip.NewReader(file)
		if err != nil {
			return nil, err
		}
		reader := bufio.NewReader(gr)
		return &fileScanner{file, bufio.NewScanner(reader)}, nil
	}

	reader := bufio.NewReader(file)
	return &fileScanner{file, bufio.NewScanner(reader)}, nil
}

// NewNmonFiles return a nmonfiles struct
func NewNmonFiles(args []string) *NmonFiles {
	nmonFiles := new(NmonFiles)
	for _, param := range args {
		paraminfo, err := os.Stat(param)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("%s doesn't exist ! skipped.\n", param)
			}
			continue
		}

		if paraminfo.IsDir() {
			entries, err := ioutil.ReadDir(param)
			check(err)
			for _, entry := range entries {
				if !entry.IsDir() {
					file := path.Join(param, entry.Name())
					nmonFiles.Add(file, path.Ext(file))
				}
			}
			continue
		}
		nmonFiles.Add(param, path.Ext(param))
	}
	return nmonFiles
}
