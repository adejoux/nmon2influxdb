// nmon2influx
// import nmon report in Influxdb
//version: 0.1
// author: adejoux@djouxtech.net

package main

import (
    influxdb "github.com/influxdb/influxdb/client"
    "text/template"
    "flag"
    "fmt"
    "path"
    "sort"
    "regexp"
    "encoding/json"
    "bufio"
    "strings"
    "strconv"
    "os"
    "time"
)
const timeformat = "15:04:05,02-Jan-2006"
var hostRegexp = regexp.MustCompile(`^AAA,host,(\S+)`)
var timeRegexp = regexp.MustCompile(`^ZZZZ,([^,]+),(.*)$`)
var intervalRegexp = regexp.MustCompile(`^AAA,interval,(\d+)`)
var headerRegexp = regexp.MustCompile(`^AAA|^BBB|^UARG|,T\d`)
var infoRegexp = regexp.MustCompile(`AAA,(.*)`)
var diskRegexp = regexp.MustCompile(`^DISK`)
var statsRegexp = regexp.MustCompile(`[^Z]+,(T\d+)`)


//
//helper functions
//
func check(e error) {
    if e != nil {
        panic(e)
    }
}

func ConvertTimeStamp(s string) int64 {
  t, err := time.Parse(timeformat, s)
  check(err)
  return t.Unix()
}

func ParseFile(filepath string) *bufio.Scanner {
    file, err := os.Open(filepath)
    check(err)

    //defer file.Close()
    reader := bufio.NewReader(file)
    scanner := bufio.NewScanner(reader)
    scanner.Split(bufio.ScanLines)
    return scanner
}

func (influx *Influx) AppendText(text string) {
    influx.TextContent +=  ReplaceComma(text)
}

func ReplaceComma(s string) (string) {
    return "<tr><td>" + strings.Replace(s, ",", "</td><td>", 1) + "</td></tr>"
}

//
// DataSerie structure
// contains the columns and points to insert in InfluxDB
//

type DataSerie struct {
    Columns []string
    PointSeq int
    Points [50][]interface{}
}

//
// influx structure
// contains the main structures and methods used to parse nmon files and upload data in Influxdb
//

type Influx struct {
    Client *influxdb.Client
    MaxPoints int
    DataSeries map[string]DataSerie
    TimeStamps map[string]int64
    Hostname string
    TextContent string
    starttime int64
    stoptime int64
}

// initialize a Influx structure
func NewInflux() *Influx {
    return &Influx{DataSeries: make(map[string]DataSerie), TimeStamps: make(map[string]int64),  MaxPoints: 50}

}

func (influx *Influx) GetTimeStamp(label string) int64 {
    if val, ok := influx.TimeStamps[label]; ok {
        return val
    } else {
        fmt.Printf("no time label for %s\n", label)
        os.Exit(1)
    }

    return 0
}

func (influx *Influx) GetColumns(serie string) ([]string) {
   return influx.DataSeries[serie].Columns
}

func (influx *Influx) GetFilteredColumns(serie string, filter string) ([]string) {
    var res []string
    for _, field := range influx.DataSeries[serie].Columns {
        if strings.Contains(field,filter) {
            res = append(res,field)
        }
    }
    return res
}

func (influx *Influx) AddData(serie string, timestamp int64, elems []string) {

    dataSerie := influx.DataSeries[serie]

    if len(dataSerie.Columns) == 0 {
        //fmt.Printf("No defined fields for %s. No datas inserted\n", serie)
        return
    }

    if len(dataSerie.Columns) != len(elems) {
        return
    }

    point := []interface{}{}
    point = append(point, timestamp)
    for i := 0; i < len(elems); i++ {
        // try to convert string to integer
        value, err := strconv.ParseFloat(elems[i],64)
        if err != nil {
            //if not working, use string
            point = append(point, elems[i])
        } else {
            //send integer if it worked
            point = append(point, value)
        }
    }

    if dataSerie.PointSeq == influx.MaxPoints  {
        influx.WriteData(serie)
        dataSerie.PointSeq = 0
    }

    dataSerie.Points[dataSerie.PointSeq] = point
    dataSerie.PointSeq += 1
    influx.DataSeries[serie]=dataSerie
}

func (influx *Influx) WriteTemplate(tmplfile string) {

    var tmplname string
    tmpl := template.New("grafana")

    if _, err := os.Stat(tmplfile); os.IsNotExist(err) {
        fmt.Printf("no such file or directory: %s\n", tmplfile)
        fmt.Printf("ERROR: unable to parse grafana template. Using default template.\n")
        tmpl.Parse(influxtempl)
        tmplname="grafana"
    } else {
        tmpl.ParseFiles(tmplfile)
        tmplname=path.Base(tmplfile)
    }

    // open output file
    filename := influx.Hostname + "_dashboard"
    fo, err := os.Create(filename)
    check(err)

    // make a write buffer
    w := bufio.NewWriter(fo)
    err2 := tmpl.ExecuteTemplate(w, tmplname, influx)
    check(err2)
    w.Flush()
    fo.Close()

    fmt.Printf("Writing GRAFANA dashboard: %s\n",filename)

}

func (influx *Influx) WriteData(serie string) {

    dataSerie := influx.DataSeries[serie]
    series := &influxdb.Series{}

    series.Name = influx.Hostname + "_" + serie

    series.Columns = append([]string{"time"}, dataSerie.Columns...)

    for i := 0; i < len(dataSerie.Points); i++ {
        if dataSerie.Points[i] == nil {
            break
        }
        series.Points = append(series.Points, dataSerie.Points[i])
    }

    client := influx.Client
    if err := client.WriteSeriesWithTimePrecision([]*influxdb.Series{series}, "s"); err != nil {
        data, err2 := json.Marshal(series)
        if err2 != nil {
            panic(err2)
        }
        fmt.Printf("%s\n", data)
        panic(err)
    }
}


func (influx *Influx) InitSession(host string, user string, pass string) {
    database := "nmon_reports"
    client, err := influxdb.NewClient(&influxdb.ClientConfig{Host: host})
    check(err)

    admins, err := client.GetClusterAdminList()
    check(err)

    if len(admins) == 1 {
        fmt.Printf("No administrator defined. Creating user %s with password %s\n", user, pass)
        if err := client.CreateClusterAdmin(user, pass); err != nil {
            panic(err)
        }
    }

    dbs, err := client.GetDatabaseList()
    check(err)

    dbexists := false

    //checking if database exists
    for _, v := range dbs {
        if v["name"] == database {
          dbexists = true
        }
    }

    if !dbexists {
        fmt.Printf("Creating database : %s\n", database)
        if err := client.CreateDatabase(database); err != nil {
            panic(err)
        }
    }

    dbexists = false
    //checking if grafana database exists
    for _, v := range dbs {
        if v["name"] == "grafana" {
          dbexists = true
        }
    }

    if !dbexists {
        fmt.Printf("Creating database : grafana\n")
        if err := client.CreateDatabase("grafana"); err != nil {
            panic(err)
        }
    }

    client, err = influxdb.NewClient(&influxdb.ClientConfig{
        Host: host,
        Username: user,
        Password: pass,
        Database: database,
        })
    check(err)

    client.DisableCompression()
    influx.Client = client
}

func (influx *Influx) SetTimeFrame() {
    keys := make([]string, 0, len(influx.TimeStamps))
    for k := range influx.TimeStamps {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    influx.starttime=influx.TimeStamps[keys[0]]
    influx.stoptime=influx.TimeStamps[keys[len(keys)-1]]
}

func (influx *Influx) StartTime() string {
    if influx.starttime == 0 {
        influx.SetTimeFrame()
    }
    return time.Unix(influx.starttime,0).Format(time.RFC3339)
}

func (influx *Influx) StopTime() string {
    if influx.stoptime == 0 {
        influx.SetTimeFrame()
    }
    return time.Unix(influx.stoptime,0).Format(time.RFC3339)
}

func main() {
    // parsing parameters
    file := flag.String("file", "nmonfile", "nmon file")
    tmplfile := flag.String("tmplfile", "tmplfile", "grafana dashboard template")
    nodata := flag.Bool("nodata", false, "generate dashboard only")
    nodashboard := flag.Bool("nodashboard", false, "only upload data")
    nodisk := flag.Bool("nodisk", false, "skip disk metrics")
    host := flag.String("host", "localhost:8086", "influxdb server and port")
    user := flag.String("user", "admin", "influxdb administor user")
    pass := flag.String("pass", "admin", "influxdb administor password")

    flag.Parse()

    if *file == "nmonfile" {
        fmt.Printf("error: no file provided\n")
        os.Exit(1)
    }

    influx := NewInflux()
    scanner := ParseFile(*file)

    for scanner.Scan() {
        switch {
            case diskRegexp.MatchString(scanner.Text()):
                if *nodisk == true {
                    continue
                }
            case timeRegexp.MatchString(scanner.Text()):
                matched := timeRegexp.FindStringSubmatch(scanner.Text())
                influx.TimeStamps[matched[1]]=ConvertTimeStamp(matched[2])
            case hostRegexp.MatchString(scanner.Text()):
                matched := hostRegexp.FindStringSubmatch(scanner.Text())
                influx.Hostname = matched[1]
            case infoRegexp.MatchString(scanner.Text()):
                matched := infoRegexp.FindStringSubmatch(scanner.Text())
                influx.AppendText(matched[1])
            case ! headerRegexp.MatchString(scanner.Text()):
                elems := strings.Split(scanner.Text(), ",")
                dataserie := influx.DataSeries[elems[0]]
                dataserie.Columns = elems[2:]
                influx.DataSeries[elems[0]]=dataserie
        }
    }

    if *nodata == false {
        influx.InitSession(*host, *user, *pass)
        scanner = ParseFile(*file)

        for scanner.Scan() {
            switch {
                case diskRegexp.MatchString(scanner.Text()):
                if *nodisk == true {
                    continue
                }
                case statsRegexp.MatchString(scanner.Text()):
                    matched := statsRegexp.FindStringSubmatch(scanner.Text())
                    elems := strings.Split(scanner.Text(), ",")
                    timestamp := influx.GetTimeStamp(matched[1])
                    influx.AddData(elems[0], timestamp, elems[2:])
            }
        }
        // flushing remaining data
        for serie := range influx.DataSeries {
            influx.WriteData(serie)
        }
    }

    if *nodashboard == false {
        influx.WriteTemplate(*tmplfile)
    }
}
