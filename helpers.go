// nmon2influxdb
// import nmon data in InfluxDB
// author: adejoux@djouxtech.net

package main

import (
	"github.com/codegangsta/cli"
	"log"
	"strings"
)

//
//helper functions
//
func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func ReplaceComma(s string) string {
	return "<tr><td>" + strings.Replace(s, ",", "</td><td>", 1) + "</td></tr>"
}

type Params struct {
	NoDisks   bool
	CpuAll    bool
	Server    string
	User      string
	Port      string
	File      bool
	Guser     string
	Gpass     string
	Gurl      string
	Db        string
	DS        string
	Password  string
	Template  string
	Metric    string
	StatsHost string
	Sort      string
	Limit     int
	Host      string
	Filter    string
	From      string
	To        string
	TZ        string
	Aggregate string
	Debug     bool
}

func ParseParameters(c *cli.Context) (params *Params) {

	return &Params{Metric: c.String("metric"),
		StatsHost: c.String("statshost"),
		From:      c.String("from"),
		To:        c.String("to"),
		Aggregate: c.String("aggregate"),
		NoDisks:   c.Bool("nodisks"),
		CpuAll:    c.Bool("cpus"),
		File:      c.Bool("file"),
		Filter:    c.String("filter"),
		Host:      c.String("host"),
		Guser:     c.String("guser"),
		Gpass:     c.String("gpassword"),
		Gurl:      c.String("gurl"),
		DS:        c.String("datasource"),
		Debug:     c.GlobalBool("debug"),
		Server:    c.GlobalString("server"),
		User:      c.GlobalString("user"),
		Port:      c.GlobalString("port"),
		Db:        c.GlobalString("db"),
		Password:  c.GlobalString("pass"),
		TZ:        c.GlobalString("tz"),
		Template:  c.String("template"),
	}
}

func ParseStatsParameters(c *cli.Context) (params *Params) {
	return &Params{
		Metric:    c.String("metric"),
		StatsHost: c.String("statshost"),
		From:      c.String("from"),
		To:        c.String("to"),
		Sort:      c.String("sort"),
		Limit:     c.Int("limit"),
		Filter:    c.String("filter"),
		Host:      c.String("host"),
		Aggregate: c.String("aggregate"),
		Debug:     c.GlobalBool("debug"),
		Server:    c.GlobalString("server"),
		User:      c.GlobalString("user"),
		Port:      c.GlobalString("port"),
		Db:        c.GlobalString("db"),
		Password:  c.GlobalString("pass"),
		TZ:        c.GlobalString("tz"),
	}
}
