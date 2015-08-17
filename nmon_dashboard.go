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
	dashboard, err := grafanaclient.ConvertTemplate(params.Name)
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

	panels := new(NmonPanels)

	panels.AddPanel(nmon.Hostname, "CPU", "CPU_ALL", "^User%|^Sys%|^Wait%|^Idle%", true)
	panels.AddPanel(nmon.Hostname, "LPAR", "LPAR", "PhysicalC|entitled|virtualC", false)

	row := BuildGrafanaRow("CPU", panels)
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(nmon.Hostname, "I/O Adapters", "IOADAPT", "KB", true)
	panels.AddPanel(nmon.Hostname, "PAGE", "PAGE", "pgs", false)
	row = BuildGrafanaRow("IO ADAPTER", panels)
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(nmon.Hostname, "Network", "NET", "KB", true)
	if len(nmon.DataSeries["SEA"].Columns) > 0 {
		panels.AddPanel(nmon.Hostname, "SEA", "SEA", "KB", true)
	}
	row = BuildGrafanaRow("NET", panels)
	db.Rows = append(db.Rows, row)

	db.GTime = grafanaclient.GTime{From: nmon.StartTime(), To: nmon.StopTime()}
	return db

}

type NmonPanel struct {
	Host        string
	Title       string
	Measurement string
	Filter      string
	Stack       bool
}

type NmonPanels []NmonPanel

func (panels *NmonPanels) AddPanel(host string, title string, measurement string, filter string, stack bool) {
	*panels = append(*panels, NmonPanel{Host: host,
		Title:       title,
		Measurement: measurement,
		Filter:      filter,
		Stack:       stack})
}

func BuildGrafanaRow(title string, panels *NmonPanels) grafanaclient.Row {
	row := grafanaclient.NewRow()
	row.Title = title

	for _, panel := range *panels {
		row.Panels = append(row.Panels, BuildGrafanaGraphPanel(panel))
	}

	return row
}

func BuildGrafanaGraphPanel(np NmonPanel) grafanaclient.Panel {
	panel := grafanaclient.NewPanel()
	panel.Title = np.Title
	target := grafanaclient.NewTarget()
	target.Alias = "$tag_name"
	target.Measurement = np.Measurement
	hostTag := grafanaclient.Tag{Key: "host", Value: np.Host}
	target.Tags = append(target.Tags, hostTag)
	if len(np.Filter) > 0 {
		fieldsTag := grafanaclient.Tag{Key: "name", Value: "/" + np.Filter + "/", Condition: "AND"}
		target.Tags = append(target.Tags, fieldsTag)
	}

	if np.Stack {
		panel.Stack = true
		panel.Fill = 1
		panel.Tooltip = grafanaclient.Tooltip{ValueType: "individual"}
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
