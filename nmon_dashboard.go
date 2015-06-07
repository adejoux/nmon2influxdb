// nmon2influxdb
// import nmon report in InfluxDB
// author: adejoux@djouxtech.net

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/adejoux/grafanaclient"
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
	if params.File {
		nmon.WriteDashboard(params.Template)
		return
	}
	nmon.UploadDashboard(params.Template)
	return

}

func (nmon *Nmon) WriteDashboard(tmplfile string) {

	dashboard, err := nmon.GenerateDashboard(tmplfile)

	// open output file
	filename := nmon.Hostname + "_dashboard"
	file, err := os.Create(filename)
	check(err)
	defer file.Close()

	// make a write buffer
	writer := bufio.NewWriter(file)
	dashboard.WriteTo(writer)
	writer.Flush()

	fmt.Printf("Writing GRAFANA dashboard: %s\n", filename)

}

func (nmon *Nmon) GenerateDashboard(tmplfile string) (dashboard bytes.Buffer, err error) {

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
	err = tmpl.ExecuteTemplate(&dashboard, tmplname, nmon)

	return
}

func (nmon *Nmon) UploadDashboard(tmplfile string) (err error) {
	dashboard, err := nmon.GenerateDashboard(tmplfile)

	//check if datasource for nmon2influxdb exist
	grafana := grafanaclient.NewSession(nmon.Params.Guser, nmon.Params.Gpass, nmon.Params.Gurl)
	err = grafana.DoLogon()

	if err != nil {
		return
	}

	resDs, err := grafana.GetDataSource(nmon.Params.DS)
	if resDs.Name == "" {
		var ds = grafanaclient.DataSource{Name: nmon.Params.DS,
			Type:     "influxdb_08",
			Access:   "direct",
			URL:      nmon.DbURL(),
			User:     nmon.Params.User,
			Password: nmon.Params.Password,
			Database: nmon.Params.Db,
		}
		err = grafana.CreateDataSource(ds)
		if err != nil {
			return
		}
		fmt.Printf("Grafana %s DataSource created.\n", nmon.Params.DS)
	}

	err = grafana.UploadDashboardString(dashboard.String(), true)
	if err != nil {
		fmt.Printf("Unable to upload Grafana dashboard ! \n")
	}

	fmt.Printf("Dashboard uploaded to grafana\n")
	return
}
