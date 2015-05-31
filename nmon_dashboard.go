// nmon2influxdb
// import nmon report in InfluxDB
//version: 0.1
// author: adejoux@djouxtech.net

package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"path"
	"text/template"
)

func NmonDashboard(c *cli.Context) {

	if len(c.Args()) < 1 {
		fmt.Printf("file name needs to be provided\n")
		os.Exit(1)
	}
	// parsing parameters
	params := ParseParameters(c)

	nmon := InitNmon(params)
	nmon.WriteTemplate(params.Template)

}

func (nmon *Nmon) WriteTemplate(tmplfile string) {

	var tmplname string
	tmpl := template.New("grafana")

	if _, err := os.Stat(tmplfile); os.IsNotExist(err) {
		if nmon.Debug {
			fmt.Printf("no such file or directory: %s\n", tmplfile)
			fmt.Printf("Warning: unable to parse grafana template. Using default template.\n")
		}
		tmpl.Parse(influxtempl)
		tmplname = "grafana"
	} else {
		tmpl.ParseFiles(tmplfile)
		tmplname = path.Base(tmplfile)
	}

	// open output file
	filename := nmon.Hostname + "_dashboard"
	fo, err := os.Create(filename)
	check(err)

	// make a write buffer
	w := bufio.NewWriter(fo)
	err2 := tmpl.ExecuteTemplate(w, tmplname, nmon)
	check(err2)
	w.Flush()
	fo.Close()

	fmt.Printf("Writing GRAFANA dashboard: %s\n", filename)

}
