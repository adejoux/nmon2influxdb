// nmon2influxdb
// import nmon report in InfluxDB
// author: adejoux@djouxtech.net

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/adejoux/grafanaclient"
	"github.com/codegangsta/cli"
)

var nmonFileRegexp = regexp.MustCompile(`\.(nmon|nmon.gz|nmon.bz2)$`)

const panelSize = "300px"
const linux = "linux"
const aix = "aix"

// NmonDashboard entry point for nmon dashboard sub command
func NmonDashboard(c *cli.Context) {

	if len(c.Args()) < 1 {
		fmt.Printf("file name needs to be provided\n")
		os.Exit(1)
	}

	// parsing parameters
	config := ParseParameters(c)

	file := c.Args().First()

	if nmonFileRegexp.MatchString(file) {
		NmonDashboardFile(config, file)
		return
	}

	NmonDashboardTemplate(config, file)

}

//NmonDashboardFile export dashboard to file
func NmonDashboardFile(config *Config, file string) {
	nmonFile := NmonFile{Name: file, FileType: ".nmon"}
	nmon := InitNmon(config, nmonFile)
	if config.DashboardWriteFile {
		nmon.WriteDashboard()
		return
	}

	if nmon.OS != linux && nmon.OS != aix {
		fmt.Printf("Error: unable to find if it's a Linux or AIX nmon file !\n")
		os.Exit(1)
	}

	var dashboard grafanaclient.Dashboard
	if nmon.OS == linux {
		dashboard = nmon.GenerateLinuxDashboard()
	}

	if nmon.OS == aix {
		dashboard = nmon.GenerateAixDashboard()
	}
	err := nmon.UploadDashboard(dashboard)
	check(err)
	return
}

// NmonDashboardTemplate generates dashboard from toml template
func NmonDashboardTemplate(config *Config, file string) {
	nmon := InitNmonTemplate(config)
	dashboard, err := grafanaclient.ConvertTemplate(file)
	if err != nil {
		fmt.Printf("Cannot convert template !\n")
		check(err)
	}
	err = nmon.UploadDashboard(dashboard)
	check(err)
	return

}

// WriteDashboard to file
func (nmon *Nmon) WriteDashboard() {

	var dashboard grafanaclient.Dashboard

	if nmon.OS == linux {
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

//GenerateAixDashboard custom minimal dashboard for AIX
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
		Title:          "CPU Total",
		Measurement:    "CPU_ALL",
		Filters:        NameFilter("^User%|^Sys%|^Wait%|^Idle%"),
		Group:          []string{"name"},
		Stack:          true,
		TableLegend:    true,
		LeftYAxisLabel: "%"})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Logical Partition",
		Measurement:    "LPAR",
		Filters:        NameFilter("PhysicalC|entitled|virtualC"),
		Group:          []string{"name"},
		Stack:          false,
		TableLegend:    true,
		LeftYAxisLabel: "cores"})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Run queue",
		Measurement:    "PROC",
		Filters:        NameFilter("Runnable"),
		Group:          []string{"name"},
		Stack:          false,
		TableLegend:    true,
		LeftYAxisLabel: "# threads"})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Asynchronous I/O",
		Measurement:    "PROCAIO",
		Group:          []string{"name"},
		Stack:          false,
		TableLegend:    true,
		LeftYAxisLabel: "count"})

	row := BuildGrafanaRow("CPU", panels)
	row.Height = panelSize
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Physical Memory",
		Measurement:    "MEM",
		Filters:        NameFilter("MB"),
		Group:          []string{"name"},
		Stack:          false,
		LeftYAxisLabel: "MB"})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Memory Usage",
		Measurement:    "MEMUSE",
		Filters:        NameFilter("%"),
		Group:          []string{"name"},
		Stack:          false,
		LeftYAxisLabel: "%"})
	row = BuildGrafanaRow("Memory", panels)
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Disk Adapter throughput KB/s",
		Measurement:    "IOADAPT",
		Filters:        NameFilter("KB"),
		Group:          []string{"name"},
		Stack:          true,
		LeftYAxisLabel: "KB/s"})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Paging",
		Measurement:    "PAGE",
		Filters:        NameFilter("pgs"),
		Group:          []string{"name"},
		Stack:          false,
		LeftYAxisLabel: "pg/s"})

	row = BuildGrafanaRow("IO ADAPTER", panels)
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "Disk Adapter transfers",
		Measurement: "IOADAPT",
		Filters:     NameFilter("xfer"),
		Group:       []string{"name"},
		Span:        12,
		Table:       true})

	row = BuildGrafanaRow("Disk Adapter table", panels)
	row.Height = "450px"
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Network I/O",
		Measurement:    "NET",
		Filters:        NameFilter("KB"),
		Group:          []string{"name"},
		Stack:          true,
		TableLegend:    true,
		LeftYAxisLabel: "KB/s",
		NegativeY:      "/read/",
	})
	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "Network Packets",
		Measurement: "NETPACKET",
		Group:       []string{"name"},
		Stack:       true})
	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "Network Errors",
		Measurement: "NETERROR",
		Group:       []string{"name"},
		Stack:       true})
	row = BuildGrafanaRow("Network", panels)
	db.Rows = append(db.Rows, row)

	if len(nmon.DataSeries["SEA"].Columns) > 0 {
		panels = new(NmonPanels)
		panels.AddPanel(&NmonPanel{Host: host,
			Title:          "SEA",
			Measurement:    "SEA",
			Filters:        NameFilter("KB"),
			Group:          []string{"name"},
			Stack:          true,
			LeftYAxisLabel: "KB/s"})
		if len(nmon.DataSeries["SEACHPHY"].Columns) > 0 {
			panels.AddPanel(&NmonPanel{Host: host,
				Title:          "SEA Physical Adapter Traffic Stats",
				Measurement:    "SEACHPHY",
				Filters:        NameFilter("KB"),
				Group:          []string{"name"},
				Stack:          true,
				LeftYAxisLabel: "KB/s"})
		}
		row = BuildGrafanaRow("Shared Ethernet Adapter", panels)
		db.Rows = append(db.Rows, row)
	}

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "TOP",
		Measurement:    "TOP",
		Filters:        NameFilter("%CPU"),
		Group:          []string{"name", "command"},
		Function:       "sum",
		LeftYAxisLabel: "%"})
	row = BuildGrafanaRow("TOP", panels)
	row.Collapse = true
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Filespace %Used",
		Measurement:    "JFSFILE",
		Group:          []string{"name"},
		Span:           12,
		Table:          true,
		LeftYAxisLabel: "%"})
	row = BuildGrafanaRow("Filesystems", panels)
	row.Collapse = true
	row.Height = "450px"
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Disk Read KB/s",
		Measurement:    "DISKREAD",
		Group:          []string{"name"},
		Stack:          false,
		LeftYAxisLabel: "KB/s"})

	if len(nmon.DataSeries["DISKREADSERV"].Columns) > 0 {
		panels.AddPanel(&NmonPanel{Host: host,
			Title:          "Disk Read Service Time msec/xfer",
			Measurement:    "DISKREADSERV",
			Group:          []string{"name"},
			Stack:          false,
			LeftYAxisLabel: "ms"})
	}

	row = BuildGrafanaRow("DISK READ", panels)
	row.Collapse = true
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Disk Write KB/s",
		Measurement:    "DISKWRITE",
		Group:          []string{"name"},
		Stack:          false,
		LeftYAxisLabel: "KB/s"})

	if len(nmon.DataSeries["DISKREADSERV"].Columns) > 0 {
		panels.AddPanel(&NmonPanel{Host: host,
			Title:          "Disk Write Service Time msec/xfer",
			Measurement:    "DISKWRITESERV",
			Group:          []string{"name"},
			Stack:          false,
			LeftYAxisLabel: "ms"})
	}

	row = BuildGrafanaRow("DISK WRITE", panels)
	row.Collapse = true
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Transfers from disk (reads) per second",
		Measurement:    "DISKRXFER",
		Group:          []string{"name"},
		Stack:          false,
		LeftYAxisLabel: "xfers/s"})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Disk transfers per second",
		Measurement:    "DISKXFER",
		Group:          []string{"name"},
		Stack:          false,
		LeftYAxisLabel: "xfers/s"})

	row = BuildGrafanaRow("Disk transfers", panels)
	row.Collapse = true
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Disk IO Reads per second",
		Measurement:    "DISKRIO",
		Group:          []string{"name"},
		Stack:          false,
		LeftYAxisLabel: "IOPs"})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Disk IO Writes per second",
		Measurement:    "DISKWIO",
		Group:          []string{"name"},
		Stack:          false,
		LeftYAxisLabel: "IOPs"})

	row = BuildGrafanaRow("Disk I/O", panels)
	row.Collapse = true
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Disk %Busy",
		Measurement:    "DISKRIO",
		Group:          []string{"name"},
		Stack:          false,
		LeftYAxisLabel: "%"})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Disk Wait Queue Time msec/xfer",
		Measurement:    "DISKWAIT",
		Group:          []string{"name"},
		Stack:          false,
		LeftYAxisLabel: "msec/xfer"})

	row = BuildGrafanaRow("Disk Busy/Wait queue activity", panels)
	row.Collapse = true
	db.Rows = append(db.Rows, row)

	db.GTime = grafanaclient.GTime{From: nmon.StartTime(), To: nmon.StopTime()}
	return db

}

//GenerateLinuxDashboard custom minimal dashboard for Linux
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
		Title:          "CPU Total",
		Measurement:    "CPU_ALL",
		Filters:        NameFilter("%"),
		Group:          []string{"name"},
		Stack:          true,
		TableLegend:    true,
		LeftYAxisLabel: "%"})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "SCAN",
		Measurement: "VM",
		Filters:     NameFilter("scan"),
		Group:       []string{"name"},
		Stack:       false,
		TableLegend: true})

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
	row.Height = panelSize
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:       "Memory MB",
		Measurement: "MEM",
		Filters:     NameFilter("^active|memtotal|cached|inactive"),
		Group:       []string{"name"},
		Stack:       true})

	row = BuildGrafanaRow("MEM", panels)
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Filespace %Used",
		Measurement:    "JFSFILE",
		Filters:        NameFilter(""),
		Group:          []string{"name"},
		Stack:          false,
		LeftYAxisLabel: "%"})

	row = BuildGrafanaRow("FS", panels)

	db.Rows = append(db.Rows, row)
	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "sdX Disk Write KB/s",
		Measurement:    "DISKWRITE",
		Filters:        NameFilter("sd"),
		Group:          []string{"name"},
		Stack:          true,
		LeftYAxisLabel: "KB/s"})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "dm Disk Write KB/s",
		Measurement:    "DISKWRITE",
		Filters:        NameFilter("dm"),
		Group:          []string{"name"},
		Stack:          true,
		LeftYAxisLabel: "KB/s"})

	row = BuildGrafanaRow("DISK WRITE", panels)
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "sdX Disk Read KB/s",
		Measurement:    "DISKREAD",
		Filters:        NameFilter("sd"),
		Group:          []string{"name"},
		Stack:          true,
		LeftYAxisLabel: "KB/s"})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "dm Disk Read KB/s",
		Measurement:    "DISKREAD",
		Filters:        NameFilter("dm"),
		Group:          []string{"name"},
		Stack:          true,
		LeftYAxisLabel: "KB/s"})

	row = BuildGrafanaRow("DISK READ", panels)
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Network",
		Measurement:    "NET",
		Filters:        NameFilter("eth|em|en"),
		Group:          []string{"name"},
		Stack:          true,
		TableLegend:    true,
		LeftYAxisLabel: "KB/s"})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "Docker Network",
		Measurement:    "NET",
		Filters:        NameFilter("docker"),
		Group:          []string{"name"},
		Stack:          true,
		TableLegend:    true,
		LeftYAxisLabel: "KB/s"})

	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "KVM Network",
		Measurement:    "NET",
		Filters:        NameFilter("virbr"),
		Group:          []string{"name"},
		Stack:          true,
		TableLegend:    true,
		LeftYAxisLabel: "KB/s"})

	row = BuildGrafanaRow("NET", panels)
	row.Height = panelSize
	db.Rows = append(db.Rows, row)

	panels = new(NmonPanels)
	panels.AddPanel(&NmonPanel{Host: host,
		Title:          "TOP",
		Measurement:    "TOP",
		Filters:        NameFilter("%CPU"),
		Group:          []string{"name", "command"},
		Function:       "sum",
		LeftYAxisLabel: "KB/s"})
	row = BuildGrafanaRow("TOP", panels)
	db.Rows = append(db.Rows, row)

	db.GTime = grafanaclient.GTime{From: nmon.StartTime(), To: nmon.StopTime()}
	return db

}

// NmonPanel custom Panel fro Grafana
type NmonPanel struct {
	Host            string
	Title           string
	Measurement     string
	Function        string
	Filters         []grafanaclient.Tag
	Group           []string
	Stack           bool
	Table           bool
	TableLegend     bool
	LeftYAxisLabel  string
	RightYAxisLabel string
	NegativeY       string
	Span            int
}

//NmonPanels array of NmonPanel
type NmonPanels []NmonPanel

//AddPanel append NmonPanel
func (panels *NmonPanels) AddPanel(npanel *NmonPanel) {
	*panels = append(*panels, *npanel)
}

//BuildGrafanaRow generate a row composed of panels
func BuildGrafanaRow(title string, panels *NmonPanels) grafanaclient.Row {
	row := grafanaclient.NewRow()
	row.Title = title

	for _, panel := range *panels {
		if panel.Table {
			row.Panels = append(row.Panels, BuildGrafanaTablePanel(panel))
		} else {
			row.Panels = append(row.Panels, BuildGrafanaGraphPanel(panel))
		}
	}

	return row
}

//NameFilter add a Grafana filter on name tag
func NameFilter(filter string) (tags []grafanaclient.Tag) {
	tags = append(tags, grafanaclient.Tag{Key: "name", Value: "/" + filter + "/", Condition: "AND"})
	return
}

//TagsFilter add a standard grafana filter
func TagsFilter(filters map[string]string) (tags []grafanaclient.Tag) {
	for _, key := range filters {
		tags = append(tags, grafanaclient.Tag{Key: key, Value: "/" + filters[key] + "/", Condition: "AND"})
	}
	return
}

//BuildGrafanaGraphPanel generates a grafana graph panel
func BuildGrafanaGraphPanel(np NmonPanel) grafanaclient.Panel {
	panel := grafanaclient.NewPanel()
	panel.Title = np.Title
	if np.Span > 0 {
		panel.Span = np.Span
	}
	target := grafanaclient.NewTarget()
	if len(np.Function) > 0 {
		var selects grafanaclient.Selects
		fieldFunction := grafanaclient.Select{Type: "field", Params: []string{"value"}}
		selects = append(selects, fieldFunction)
		function := grafanaclient.Select{Type: np.Function, Params: []string{}}
		selects = append(selects, function)
		target.Select = []grafanaclient.Selects{selects}
	}
	target.Measurement = np.Measurement
	hostTag := grafanaclient.Tag{Key: "host", Value: np.Host}
	target.Tags = append(target.Tags, hostTag)

	for _, filt := range np.Filters {
		target.Tags = append(target.Tags, filt)
	}

	panel.LeftYAxisLabel = np.LeftYAxisLabel
	panel.RightYAxisLabel = np.RightYAxisLabel

	if np.TableLegend {
		legend := grafanaclient.NewLegend()
		legend.Values = true
		legend.Min = true
		legend.Avg = true
		legend.Max = true
		legend.AlignAsTable = true
		panel.Legend = legend
	}

	if len(np.NegativeY) > 0 {
		seriesOverride := grafanaclient.NewSeriesOverride(np.NegativeY)
		seriesOverride.Transform = "negative-Y"
		panel.SeriesOverrides = append(panel.SeriesOverrides, seriesOverride)
	}

	if np.Stack {
		panel.Stack = true
		panel.Fill = 1
		panel.Tooltip = grafanaclient.Tooltip{ValueType: "individual"}
	}

	if len(np.Group) > 0 {
		target.GroupBy = grafanaclient.NewGroupBy()
		target.Alias = ""
	}

	for _, field := range np.Group {
		target.Alias += "$tag_" + field + " "
		target.GroupBy = append(target.GroupBy, grafanaclient.GroupBy{Type: "tag", Params: []string{field}})
	}

	panel.Targets = append(panel.Targets, target)

	return panel
}

//BuildGrafanaTablePanel generates a grafana graph panel
func BuildGrafanaTablePanel(np NmonPanel) grafanaclient.Panel {
	panel := grafanaclient.NewPanel()
	panel.Type = "table"
	panel.Title = np.Title
	if np.Span > 0 {
		panel.Span = np.Span
	}
	target := grafanaclient.NewTarget()

	target.Measurement = np.Measurement
	hostTag := grafanaclient.Tag{Key: "host", Value: np.Host}
	target.Tags = append(target.Tags, hostTag)
	panel.PageSize = 20
	target.Transform = "timeseries_to_columns"

	for _, filt := range np.Filters {
		target.Tags = append(target.Tags, filt)
	}

	if len(np.Group) > 0 {
		target.Alias = ""
	}

	for _, field := range np.Group {
		target.Alias += "$tag_" + field + " "
		target.GroupBy = append(target.GroupBy, grafanaclient.GroupBy{Type: "time", Params: []string{"15m"}})
		target.GroupBy = append(target.GroupBy, grafanaclient.GroupBy{Type: "tag", Params: []string{field}})
	}

	panel.Targets = append(panel.Targets, target)

	return panel
}

//InitGrafanaSession connects to grafana instance and setup influxdb datasource
func (nmon *Nmon) InitGrafanaSession() *grafanaclient.Session {
	//check if datasource for nmon2influxdb exist
	grafana := grafanaclient.NewSession(nmon.Config.GrafanaUser, nmon.Config.GrafanaPassword, nmon.Config.GrafanaURL)
	err := grafana.DoLogon()
	check(err)

	resDs, err := grafana.GetDataSource(nmon.Config.GrafanaDatasource)
	check(err)
	if resDs.Name == "" {
		plugins, err := grafana.GetDataSourcePlugins()

		//grafana 3.0 new plugin architecture
		if err.Error() == "HTTP 404: Data source not found" {
			plugins, pluginErr := grafana.GetPlugins("datasource")
			check(pluginErr)

			status := ""
			for _, plugin := range plugins {
				if plugin.ID == "influxdb" {
					status = "ok"
				}
			}

			if status != "ok" {
				fmt.Printf("No plugin for influxDB in Grafana !\n")
				os.Exit(1)
			}
		} else {
			check(err)
			if _, present := plugins["influxdb"]; !present {
				fmt.Printf("No plugin for influxDB in Grafana !\n")
				os.Exit(1)
			}
		}

		var ds = grafanaclient.DataSource{Name: nmon.Config.GrafanaDatasource,
			Type:      "influxdb",
			Access:    nmon.Config.GrafanaAccess,
			URL:       nmon.DbURL(),
			User:      nmon.Config.GrafanaUser,
			Password:  nmon.Config.GrafanaPassword,
			Database:  nmon.Config.InfluxdbDatabase,
			IsDefault: true,
		}
		err = grafana.CreateDataSource(ds)
		check(err)
		fmt.Printf("Grafana %s DataSource created.\n", nmon.Config.GrafanaDatasource)
	}

	return grafana
}

//UploadDashboard upload dashboard to current grafana instance
func (nmon *Nmon) UploadDashboard(dashboard grafanaclient.Dashboard) (err error) {
	grafana := nmon.InitGrafanaSession()

	err = grafana.UploadDashboard(dashboard, true)
	if err != nil {
		fmt.Printf("Unable to upload Grafana dashboard: %s ! \n", err.Error())
		return
	}

	fmt.Printf("Dashboard uploaded to grafana\n")
	return
}
