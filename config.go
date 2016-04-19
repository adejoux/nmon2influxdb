// nmon2influxdb
// import nmon data in InfluxDB

// author: adejoux@djouxtech.net

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/naoina/toml"
)

// Config is the configuration structure used by nmon2influxdb
type Config struct {
	Debug                bool
	Timezone             string
	InfluxdbUser         string
	InfluxdbPassword     string
	InfluxdbServer       string
	InfluxdbPort         string
	InfluxdbDatabase     string
	GrafanaUser          string
	GrafanaPassword      string
	GrafanaURL           string `toml:"grafana_URL"`
	GrafanaAccess        string
	GrafanaDatasource    string
	ImportSkipDisks      bool
	ImportAllCpus        bool
	ImportBuildDashboard bool
	ImportForce          bool
	ImportSkipMetrics    string
	ImportLogDatabase    string
	ImportLogRetention   string
	DashboardWriteFile   bool
	StatsLimit           int
	StatsSort            string
	StatsFilter          string
	StatsFrom            string
	StatsTo              string
	StatsHost            string
	Metric               string `toml:"metric,omitempty"`
	ListFilter           string `toml:",omitempty"`
	ListHost             string `toml:",omitempty"`
}

// InitConfig setup initial configuration with sane values
func InitConfig() Config {
	return Config{Debug: false,
		Timezone:             "Europe/Paris",
		InfluxdbUser:         "root",
		InfluxdbPassword:     "root",
		InfluxdbServer:       "localhost",
		InfluxdbPort:         "8086",
		InfluxdbDatabase:     "nmon_reports",
		GrafanaUser:          "admin",
		GrafanaPassword:      "admin",
		GrafanaURL:           "http://localhost:3000",
		GrafanaAccess:        "direct",
		GrafanaDatasource:    "nmon2influxdb",
		ImportSkipDisks:      false,
		ImportAllCpus:        false,
		ImportBuildDashboard: false,
		ImportForce:          false,
		ImportLogDatabase:    "nmon2influxdb_log",
		ImportLogRetention:   "1d",
		DashboardWriteFile:   false,
		ImportSkipMetrics:    "JFSINODE|TOP",
		StatsLimit:           20,
		StatsSort:            "mean",
		StatsFilter:          "",
		StatsFrom:            "",
		StatsTo:              "",
		StatsHost:            "",
	}
}

//GetCfgFile returns the current configuration file path
func GetCfgFile() (cfgfile string) {
	currUser, _ := user.Current()
	home := currUser.HomeDir
	cfgfile = filepath.Join(home, ".nmon2influxdb.cfg")
	return
}

//IsNotFile returns true if the file doesn't exist
func IsNotFile(file string) bool {
	stat, err := os.Stat(file)
	if err != nil {
		return true
	}
	if stat.Mode().IsRegular() {
		return false
	}

	return true
}

//BuildCfgFile creates a default configuration file
func (config *Config) BuildCfgFile(cfgfile string) {
	file, err := os.Create(cfgfile)
	check(err)
	defer file.Close()
	writer := bufio.NewWriter(file)
	b, err := toml.Marshal(*config)
	check(err)
	r := bytes.NewReader(b)
	r.WriteTo(writer)
	writer.Flush()
	fmt.Printf("Generating default configuration file : %s\n", cfgfile)
}

// LoadCfgFile loads current configuration file settings
func (config *Config) LoadCfgFile() {

	cfgfile := GetCfgFile()
	if IsNotFile(cfgfile) {
		config.BuildCfgFile(cfgfile)
	}

	file, err := os.Open(cfgfile)
	if err != nil {
		fmt.Printf("Error opening configuration file %s\n", cfgfile)
		return
	}

	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		check(err)
	}

	if err := toml.Unmarshal(buf, &config); err != nil {
		check(err)
	}

}
