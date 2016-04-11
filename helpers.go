// nmon2influxdb
// import nmon data in InfluxDB
// author = adejoux@djouxtech.net

package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"
)

//
//helper functions
//
func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func checkInfo(e error) {
	if e != nil {
		log.Printf("info: %s", e)
	}
}

//ReplaceComma replaces comma by html tabs tag
func ReplaceComma(s string) string {
	return "<tr><td>" + strings.Replace(s, ",", "</td><td>", 1) + "</td></tr>"
}

// ParseParameters parse parameter from command line in Config struct
func ParseParameters(c *cli.Context) (config *Config) {
	config = new(Config)
	*config = InitConfig()
	config.LoadCfgFile()

	config.Metric = c.String("metric")
	config.StatsHost = c.String("statshost")
	config.StatsFrom = c.String("from")
	config.StatsTo = c.String("to")
	config.StatsLimit = c.Int("limit")
	config.StatsFilter = c.String("filter")
	config.ImportSkipDisks = c.Bool("nodisks")
	config.ImportAllCpus = c.Bool("cpus")
	config.ImportBuildDashboard = c.Bool("build")
	config.ImportSkipMetrics = c.String("skip_metrics")
	config.ImportLogDatabase = c.String("log_database")
	config.ImportLogRetention = c.String("log_retention")
	config.DashboardWriteFile = c.Bool("file")
	config.ListFilter = c.String("filter")
	config.ImportForce = c.Bool("force")
	config.ListHost = c.String("host")
	config.GrafanaUser = c.String("guser")
	config.GrafanaPassword = c.String("gpassword")
	config.GrafanaAccess = c.String("gaccess")
	config.GrafanaURL = c.String("gurl")
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

//Checksum generates SHA1 file checksum
func Checksum(filePath string) (myhash string, err error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	hash := sha1.New()
	if _, err = io.Copy(hash, file); err != nil {
		return
	}
	myhash = hex.EncodeToString(hash.Sum(result))

	return
}
