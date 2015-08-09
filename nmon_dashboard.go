// nmon2influxdb
// import nmon report in InfluxDB
// author: adejoux@djouxtech.net

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/adejoux/grafanaclient"
	"github.com/codegangsta/cli"
	"os"
)

func NmonDashboardFile(c *cli.Context) {

	if len(c.Args()) < 1 {
		fmt.Printf("file name needs to be provided\n")
		os.Exit(1)
	}
	// parsing parameters
	params := ParseParameters(c)

	nmon := InitNmon(params)
	if params.File {
		nmon.WriteDashboard()
		return
	}
	dashboard := nmon.GenerateDashboard()
	err := nmon.UploadDashboard(dashboard)
	check(err)
	return
}

func NmonDashboardTemplate(c *cli.Context) {
	if len(c.Args()) < 1 {
		fmt.Printf("file name needs to be provided\n")
		os.Exit(1)
	}
	// parsing parameters
	params := ParseParameters(c)
	nmon := InitNmonTemplate(params)
	dashboard, err := grafanaclient.ConvertTemplate(params.Filepath)
	if err != nil {
		fmt.Printf("Cannot convert template !\n")
		check(err)
	}
	err = nmon.UploadDashboard(dashboard)
	check(err)
	return

}

func (nmon *Nmon) WriteDashboard() {

	dashboard := nmon.GenerateDashboard()

	// open output file
	filename := nmon.Hostname + "_dashboard"
	file, err := os.Create(filename)
	check(err)
	defer file.Close()

	// make a write buffer
	writer := bufio.NewWriter(file)
	b, _ := json.Marshal(dashboard)
	r := bytes.NewReader(b)
	r.WriteTo(writer)
	writer.Flush()

	fmt.Printf("Writing GRAFANA dashboard: %s\n", filename)

}

func (nmon *Nmon) GenerateDashboard() grafanaclient.Dashboard {

	db := grafanaclient.Dashboard{Editable: true}

	db.Title = fmt.Sprintf("%s nmon report", nmon.Hostname)

	infoRow := grafanaclient.NewRow()
	infoRow.Title = "INFORMATION"
	infoRow.Collapse = true
	panel := grafanaclient.Panel{Type: "text", Editable: true, Mode: "html"}
	panel.Content = nmon.TextContent
	infoRow.Panels = append(infoRow.Panels, panel)
	db.Rows = append(db.Rows, infoRow)

	cpuRow := grafanaclient.NewRow()
	cpuRow.Title = "CPU"
	cpuPanel := BuildGrafanaGraphPanel(nmon.Hostname, "cpu %", "CPU_ALL", "^User%|^Sys%|^Wait%|^Idle%")
	cpuRow.Panels = append(cpuRow.Panels, cpuPanel)
	ecPanel := BuildGrafanaGraphPanel(nmon.Hostname, "EC%", "CPU_ALL", "^EC_User%|^EC_Sys%|^EC_Wait%|^EC_Idle%")
	cpuRow.Panels = append(cpuRow.Panels, ecPanel)
	db.Rows = append(db.Rows, cpuRow)
	db.GTime = grafanaclient.GTime{From: nmon.StartTime(), To: nmon.StopTime()}
	return db
}

func BuildGrafanaGraphPanel(host string, title string, measurement string, filter string) grafanaclient.Panel {
	panel := grafanaclient.NewPanel()
	panel.Title = title
	target := grafanaclient.NewTarget()

	target.Measurement = measurement
	hostTag := grafanaclient.Tag{Key: "host", Value: host}
	target.Tags = append(target.Tags, hostTag)
	if len(filter) > 0 {
		fieldsTag := grafanaclient.Tag{Key: "name", Value: "/" + filter + "/", Condition: "AND"}
		target.Tags = append(target.Tags, fieldsTag)
	}

	target.GroupByTags = []string{"name"}

	panel.Targets = append(panel.Targets, target)

	return panel
}

func (nmon *Nmon) InitGrafanaSession() *grafanaclient.Session {
	//check if datasource for nmon2influxdb exist
	grafana := grafanaclient.NewSession(nmon.Params.Guser, nmon.Params.Gpass, nmon.Params.Gurl)
	err := grafana.DoLogon()
	check(err)

	resDs, err := grafana.GetDataSource(nmon.Params.DS)
	if resDs.Name == "" {
		var ds = grafanaclient.DataSource{Name: nmon.Params.DS,
			Type:     "influxdb_09",
			Access:   "direct",
			URL:      nmon.DbURL(),
			User:     nmon.Params.User,
			Password: nmon.Params.Password,
			Database: nmon.Params.Db,
		}
		err = grafana.CreateDataSource(ds)
		check(err)
		fmt.Printf("Grafana %s DataSource created.\n", nmon.Params.DS)
	}

	return grafana
}

func (nmon *Nmon) UploadDashboard(dashboard grafanaclient.Dashboard) (err error) {
	grafana := nmon.InitGrafanaSession()

	err = grafana.UploadDashboard(dashboard, true)
	if err != nil {
		fmt.Printf("Unable to upload Grafana dashboard ! \n")
	}

	fmt.Printf("Dashboard uploaded to grafana\n")
	return
}
