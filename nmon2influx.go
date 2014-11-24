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
    "regexp"
    "encoding/json"
    "bufio"
    "strings"
    "strconv"
    "os"
    "time"
)

var hostRegexp = regexp.MustCompile(`^AAA,host,(\S+)`)
var timeRegexp = regexp.MustCompile(`^AAA,time,(\S+)`)
var dateRegexp = regexp.MustCompile(`^AAA,date,(\S+)`)
var intervalRegexp = regexp.MustCompile(`^AAA,interval,(\d+)`)
var headerRegexp = regexp.MustCompile(`^AAA|^BBB|^UARG|,T\d`)
var infoRegexp = regexp.MustCompile(`AAA,(.*)`)
var diskRegexp = regexp.MustCompile(`^DISK`)
var statsRegexp = regexp.MustCompile(`[^Z]+,T(\d+)`)

type Config struct {
    Hostname string
    Date string
    Time string
    Interval int64
    startTime int64
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func (c *Config) StartTime() int64 {
    if c.startTime == 0 {
        const timeformat = "02-Jan-2006 15:04:05"
        t, err := time.Parse(timeformat, c.Date + " " + c.Time)
        check(err)
        c.startTime = t.Unix()
        fmt.Println(t.Unix())
    }

    return c.startTime
}

func (c *Config) GetTimestamp(step int64) int64 {
  return  c.StartTime() + step * c.Interval
}

func StringToInt64(s string) int64 {
    intvalue, err := strconv.Atoi(s)
    check(err)

    return int64(intvalue)
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

type Influx struct {
    Client *influxdb.Client
    MaxPoints int
    DataSeries map[string]DataSerie
    Hostname string
    TextContent string
}

type DataSerie struct {
    Columns []string
    PointSeq int
    Points [50][]interface{}
}

func (influx *Influx) GetColumns(serie string) ([]string) {
   return influx.DataSeries[serie].Columns
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

func (influx *Influx) WriteTemplate() {

    tmpl := template.New("influxtempl")
    tmpl.Parse(influxtempl)

    // open output file
    filename := influx.Hostname + "_dashboard"
    fo, err := os.Create(filename)
    check(err)

    // make a write buffer
    w := bufio.NewWriter(fo)
    err2 := tmpl.Execute(w, influx)
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


func (influx *Influx) InitSession(admin string, pass string) {
    database := "nmon_reports"
    client, err := influxdb.NewClient(&influxdb.ClientConfig{})
    check(err)

    admins, err := client.GetClusterAdminList()
    check(err)

    if len(admins) == 1 {
        fmt.Printf("No administrator defined. Creating user %s with password %s\n", admin, pass)
        if err := client.CreateClusterAdmin(admin, pass); err != nil {
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

    users, err := client.GetDatabaseUserList(database)
    check(err)

    dbuser := database + "user"
    dbpass := "pass"

    if len(users) == 0 {
        fmt.Printf("Creating database user : %s\n", dbuser)
        if err := client.CreateDatabaseUser(database, dbuser, dbpass); err != nil {
            panic(err)
        }

        if err := client.AlterDatabasePrivilege(database, dbuser, true); err != nil {
            panic(err)
        }
    }

    client, err = influxdb.NewClient(&influxdb.ClientConfig{
        Username: dbuser,
        Password: dbpass,
        Database: database,

        })
    check(err)

    client.DisableCompression()
    influx.Client = client
}

func NewInflux() *Influx {
    return &Influx{DataSeries: make(map[string]DataSerie), MaxPoints: 50}

}

func (influx *Influx) AppendText(text string) {
    influx.TextContent +=  ReplaceComma(text)
}

func ReplaceComma(s string) (string) {
    return "<tr><td>" + strings.Replace(s, ",", "</td><td>", 1) + "</td></tr>"
}

func main() {
    // parsing parameters
    file := flag.String("file", "nmonfile", "nmon file")
    nodata := flag.Bool("nodata", false, "generate template only")
    notmpl := flag.Bool("notmpl", false, "only upload data")
    nodisk := flag.Bool("nodisk", false, "skip disk metrics")
    admin := flag.String("admin", "admin", "influxdb administor user")
    pass := flag.String("pass", "admin", "influxdb administor password")

    flag.Parse()

    if *file == "nmonfile" {
        fmt.Printf("error: no file provided\n")
        os.Exit(1)
    }

    var config Config
    influx := NewInflux()

    scanner := ParseFile(*file)

    for scanner.Scan() {
        switch {
            case diskRegexp.MatchString(scanner.Text()):
                if *nodisk == true {
                    continue
                }
            case hostRegexp.MatchString(scanner.Text()):
                matched := hostRegexp.FindStringSubmatch(scanner.Text())
                influx.Hostname = matched[1]
            case timeRegexp.MatchString(scanner.Text()):
                matched := timeRegexp.FindStringSubmatch(scanner.Text())
                config.Time = matched[1]
            case dateRegexp.MatchString(scanner.Text()):
                matched := dateRegexp.FindStringSubmatch(scanner.Text())
                config.Date = matched[1]
            case infoRegexp.MatchString(scanner.Text()):
                matched := infoRegexp.FindStringSubmatch(scanner.Text())
                influx.AppendText(matched[1])
            case intervalRegexp.MatchString(scanner.Text()):
                matched := intervalRegexp.FindStringSubmatch(scanner.Text())
                config.Interval = StringToInt64(matched[1])
            case ! headerRegexp.MatchString(scanner.Text()):
                elems := strings.Split(scanner.Text(), ",")
                dataserie := influx.DataSeries[elems[0]]
                dataserie.Columns = elems[2:]
                influx.DataSeries[elems[0]]=dataserie
        }
    }

    if *nodata == false {
        influx.InitSession(*admin, *pass)
        scanner = ParseFile(*file)

        for scanner.Scan() {
            switch {
                case diskRegexp.MatchString(scanner.Text()):
                if *nodisk == true {
                    continue
                }
                case statsRegexp.MatchString(scanner.Text()):
                    matched := statsRegexp.FindStringSubmatch(scanner.Text())
                    step := StringToInt64(matched[1])
                    timestamp := config.GetTimestamp(step)
                    elems := strings.Split(scanner.Text(), ",")
                    influx.AddData(elems[0], timestamp, elems[2:])
            }
        }
        // flushing remaining data
        for serie := range influx.DataSeries {
            influx.WriteData(serie)
        }
    }

    if *notmpl == false {
        influx.WriteTemplate()
    }
}
