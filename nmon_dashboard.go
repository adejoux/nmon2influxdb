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

	nmon := InitNmon(params, c.Args().First())
	if params.File {
		nmon.WriteDashboard()
		return
	}

	if nmon.OS != "linux" && nmon.OS != "aix" {
		fmt.Printf("Error: unable to find if it's a Linux or AIX nmon file !\n")
		os.Exit(1)
	}

	var dashboard grafanaclient.Dashboard
	if nmon.OS == "linux" {
		dashboard = nmon.GenerateLinuxDashboard()
	}

	if nmon.OS == "aix" {
		dashboard = nmon.GenerateAixDashboard()
	}
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
	dashboard, err := grafanaclient.ConvertTemplate(c.Args().First())
	if err != nil {
		fmt.Printf("Cannot convert template !\n")
		check(err)
	}
	err = nmon.UploadDashboard(dashboard)
	check(err)
	return

}

func (nmon *Nmon) WriteDashboard() {

	var dashboard grafanaclient.Dashboard

	if nmon.OS == "linux" {
		dashboard = nmon.GenerateLinuxDashboard()
	}
	if nmon.OS == "aix" {
		dashboard = nmon.GenerateAixDashboard()
	}

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

func (nmon *Nmon) GenerateAixDashboard() grafanaclient.Dashboard {

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

	host := nmon.Hostname

	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "CPU",
		Measurement: "CPU_ALL",
		Filters:     NameFilter("^User%|^Sys%|^Wait%|^Idle%"),
		Group:       []string{"name"},
		Stack:       false})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "LPAR",
		Measurement: "LPAR",
		Filters:     NameFilter("PhysicalC|entitled|virtualC"),
		Group:       []string{"name"},
		Stack:       false})

	row := BuildGrafanaRow("CPU", panels)
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "I/O Adapters",
		Measurement: "IOADAPT",
		Filters:     NameFilter("KB"),
		Group:       []string{"name"},
		Stack:       true})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "PAGE",
		Measurement: "PAGE",
		Filters:     NameFilter("pgs"),
		Group:       []string{"name"},
		Stack:       false})

	row = BuildGrafanaRow("IO ADAPTER", panels)
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "Network",
		Measurement: "NET",
		Filters:     NameFilter("KB"),
		Group:       []string{"name"},
		Stack:       true})

	if len(nmon.DataSeries["SEA"].Columns) > 0 {
		panels.AddPanel(&NmonPanel{Host: host,
			Title:       "SEA",
			Measurement: "SEA",
			Filters:     NameFilter("KB"),
			Group:       []string{"name"},
			Stack:       true})
	}
	row = BuildGrafanaRow("NET", panels)
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "TOP",
		Measurement: "TOP",
		Filters:     NameFilter("%CPU"),
		Group:       []string{"name", "command"},
		Stack:       true})
	row = BuildGrafanaRow("TOP", panels)
	db.Rows = append(db.Rows, row)

	db.GTime = grafanaclient.GTime{From: nmon.StartTime(), To: nmon.StopTime()}
	return db

}

func (nmon *Nmon) GenerateLinuxDashboard() grafanaclient.Dashboard {

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

	host := nmon.Hostname
	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "CPU",
		Measurement: "CPU_ALL",
		Filters:     NameFilter("%"),
		Group:       []string{"name"},
		Stack:       true})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "SCAN",
		Measurement: "VM",
		Filters:     NameFilter("scan"),
		Group:       []string{"name"},
		Stack:       false})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "STEAL",
		Measurement: "VM",
		Filters:     NameFilter("steal"),
		Group:       []string{"name"},
		Stack:       false})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "COUNTERS",
		Measurement: "VM",
		Filters:     NameFilter("nr"),
		Group:       []string{"name"},
		Stack:       false})

	row := BuildGrafanaRow("CPU", panels)
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "MEM",
		Measurement: "MEM",
		Filters:     NameFilter("^active|memtotal|cached|inactive"),
		Group:       []string{"name"},
		Stack:       true})

	row = BuildGrafanaRow("MEM", panels)
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "FS USAGE",
		Measurement: "JFSFILE",
		Filters:     NameFilter(""),
		Group:       []string{"name"},
		Stack:       false})

	row = BuildGrafanaRow("FS", panels)

	db.Rows = append(db.Rows, row)
	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "DISK WRITE",
		Measurement: "DISKWRITE",
		Filters:     NameFilter("sd"),
		Group:       []string{"name"},
		Stack:       true})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "DM DISK WRITE",
		Measurement: "DISKWRITE",
		Filters:     NameFilter("dm"),
		Group:       []string{"name"},
		Stack:       true})

	row = BuildGrafanaRow("DISK WRITE", panels)
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "DISK READ",
		Measurement: "DISKREAD",
		Filters:     NameFilter("sd"),
		Group:       []string{"name"},
		Stack:       true})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "DM DISK READ",
		Measurement: "DISKREAD",
		Filters:     NameFilter("dm"),
		Group:       []string{"name"},
		Stack:       true})

	row = BuildGrafanaRow("DISK READ", panels)
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "Network",
		Measurement: "NET",
		Filters:     NameFilter("eth|em|en"),
		Group:       []string{"name"},
		Stack:       true})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "Docker Network",
		Measurement: "NET",
		Filters:     NameFilter("docker"),
		Group:       []string{"name"},
		Stack:       true})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "KVM Network",
		Measurement: "NET",
		Filters:     NameFilter("virbr"),
		Group:       []string{"name"},
		Stack:       true})

	row = BuildGrafanaRow("NET", panels)
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "TOP",
		Measurement: "TOP",
		Filters:     NameFilter("%CPU"),
		Group:       []string{"name", "command"},
		Stack:       true})
	row = BuildGrafanaRow("TOP", panels)
	db.Rows = append(db.Rows, row)

	db.GTime = grafanaclient.GTime{From: nmon.StartTime(), To: nmon.StopTime()}
	return db

}

type NmonPanel struct {
	Host        string
	Title       string
	Measurement string
	Filters     []grafanaclient.Tag
	Group       []string
	Stack       bool
}

type NmonPanels []NmonPanel

func (panels *NmonPanels) AddPanel(npanel *NmonPanel) {
	*panels = append(*panels, *npanel)
}

func BuildGrafanaRow(title string, panels *NmonPanels) grafanaclient.Row {
	row := grafanaclient.NewRow()
	row.Title = title

	for _, panel := range *panels {
		row.Panels = append(row.Panels, BuildGrafanaGraphPanel(panel))
	}

	return row
}

func NameFilter(filter string) (tags []grafanaclient.Tag) {
	tags = append(tags, grafanaclient.Tag{Key: "name", Value: "/" + filter + "/", Condition: "AND"})
	return
}

func TagsFilter(filters map[string]string) (tags []grafanaclient.Tag) {
	for _, key := range filters {
		tags = append(tags, grafanaclient.Tag{Key: key, Value: "/" + filters[key] + "/", Condition: "AND"})
	}
	return
}

func BuildGrafanaGraphPanel(np NmonPanel) grafanaclient.Panel {
	panel := grafanaclient.NewPanel()
	panel.Title = np.Title
	target := grafanaclient.NewTarget()
	target.Measurement = np.Measurement
	hostTag := grafanaclient.Tag{Key: "host", Value: np.Host}
	target.Tags = append(target.Tags, hostTag)

	for _, filt := range np.Filters {
		target.Tags = append(target.Tags, filt)
	}

	if np.Stack {
		panel.Stack = true
		panel.Fill = 1
		panel.Tooltip = grafanaclient.Tooltip{ValueType: "individual"}
	}

	if len(np.Group) > 0 {
		target.GroupByTags = np.Group
		target.Alias = ""
	}

	for _, field := range np.Group {
		target.Alias += "$tag_" + field + " "
	}

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
		plugins, err := grafana.GetDataSourcePlugins()
		check(err)
		if _, present := plugins["influxdb"]; !present {
			fmt.Printf("No plugin for influxDB in Grafana !\n")
			os.Exit(1)
		}

		var ds = grafanaclient.DataSource{Name: nmon.Params.DS,
			Type:      plugins["influxdb"].Type,
			Access:    "proxy",
			URL:       nmon.DbURL(),
			User:      nmon.Params.User,
			Password:  nmon.Params.Password,
			Database:  nmon.Params.Db,
			IsDefault: true,
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
