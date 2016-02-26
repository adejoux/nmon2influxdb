// nmon2influxdb
// import nmon data in InfluxDB
// author = adejoux@djouxtech.net

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

func check_info(e error) {
	if e != nil {
		log.Printf("info: %s", e)
	}
}

func ReplaceComma(s string) string {
	return "<tr><td>" + strings.Replace(s, ",", "</td><td>", 1) + "</td></tr>"
}

func ParseParameters(c *cli.Context) (config *Config) {
	config = new(Config)
	*config = InitConfig()
	config.LoadCfgFile()

	config.Metric = c.String("metric")
	config.StatsHost = c.String("statshost")
	config.StatsFrom = c.String("from")
	config.StatsTo = c.String("to")
	config.ImportSkipDisks = c.Bool("nodisks")
	config.ImportAllCpus = c.Bool("cpus")
	config.ImportBuildDashboard = c.Bool("build")
	config.ImportSkipMetrics = c.String("skip_metrics")
	config.DashboardWriteFile = c.Bool("file")
	config.ListFilter = c.String("filter")
	config.ImportForce = c.Bool("force")
	config.ListHost = c.String("host")
	config.GrafanaUser = c.String("guser")
	config.GrafanaPassword = c.String("gpassword")
	config.GrafanaAccess = c.String("gaccess")
	config.GrafanaUrl = c.String("gurl")
	config.GrafanaDatasource = c.String("datasource")
	config.Debug = c.GlobalBool("debug")
	config.InfluxdbServer = c.GlobalString("server")
	config.InfluxdbUser = c.GlobalString("user")
	config.InfluxdbPort = c.GlobalString("port")
	config.InfluxdbDatabase = c.GlobalString("db")
	config.InfluxdbPassword = c.GlobalString("pass")
	config.Timezone = c.GlobalString("tz")
	config.DashboardTemplate = c.String("template")
	return

}
