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

	"github.com/adejoux/influxdbclient"
	"github.com/codegangsta/cli"
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
	ImportDataRetention  string
	ImportSSHUser        string `toml:"import_ssh_user"`
	ImportSSHKey         string `toml:"import_ssh_key"`
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
	currUser, _ := user.Current()
	home := currUser.HomeDir
	sshKey := filepath.Join(home, "/.ssh/id_rsa")

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
		ImportLogRetention:   "2d",
		ImportSSHUser:        currUser.Username,
		ImportSSHKey:         sshKey,
		DashboardWriteFile:   false,
		ImportSkipMetrics:    "JFSINODE|TOP|PCPU",
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

// AddDashboardParams initialize default parameters for dashboard
func (config *Config) AddDashboardParams() {
	dfltConfig := InitConfig()
	dfltConfig.LoadCfgFile()

	config.GrafanaAccess = dfltConfig.GrafanaAccess
	config.GrafanaURL = dfltConfig.GrafanaURL
	config.GrafanaDatasource = dfltConfig.GrafanaDatasource
	config.GrafanaUser = dfltConfig.GrafanaUser
	config.GrafanaPassword = dfltConfig.GrafanaPassword
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
	if c.IsSet("cpus") {
		config.ImportAllCpus = c.Bool("cpus")
	}
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

	if config.ImportBuildDashboard {
		config.AddDashboardParams()
	}

	return

}

// connect connect to the specified influxdb database
func (config *Config) connectDB(db string) *influxdbclient.InfluxDB {
	influxdbConfig := influxdbclient.InfluxDBConfig{
		Host:     config.InfluxdbServer,
		Port:     config.InfluxdbPort,
		Database: db,
		User:     config.InfluxdbUser,
		Pass:     config.InfluxdbPassword,
		Debug:    config.Debug,
	}

	influxdb, err := influxdbclient.NewInfluxDB(influxdbConfig)
	check(err)

	return &influxdb
}

// GetDataDB create or get the influxdb database like defined in config
func (config *Config) GetDataDB() *influxdbclient.InfluxDB {

	influxdb := config.connectDB(config.InfluxdbDatabase)

	if exist, _ := influxdb.ExistDB(config.InfluxdbDatabase); exist != true {
		_, createErr := influxdb.CreateDB(config.InfluxdbDatabase)
		check(createErr)
	}
	// Get default retention policy name
	policyName, policyErr := influxdb.GetDefaultRetentionPolicy()
	check(policyErr)
	// update default retention policy if ImportDataRetention is set
	if len(config.ImportDataRetention) > 0 {
		_, err := influxdb.UpdateRetentionPolicy(policyName, config.ImportDataRetention, true)
		check(err)
	}
	return influxdb
}

// GetLogDB create or get the influxdb database like defined in config
func (config *Config) GetLogDB() *influxdbclient.InfluxDB {

	influxdb := config.connectDB(config.ImportLogDatabase)

	if exist, _ := influxdb.ExistDB(config.ImportLogDatabase); exist != true {
		_, err := influxdb.CreateDB(config.ImportLogDatabase)
		check(err)
		_, err = influxdb.SetRetentionPolicy("log_retention", config.ImportLogRetention, true)
		check(err)
	} else {
		_, err := influxdb.UpdateRetentionPolicy("log_retention", config.ImportLogRetention, true)
		check(err)
	}
	return influxdb
}
