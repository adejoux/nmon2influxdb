// nmon2influxdb
// import nmon data in InfluxDB

// author: adejoux@djouxtech.net

package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {

	config := InitConfig()

	config.LoadCfgFile()

	// cannot set values directly for boolean flags
	if config.DashboardWriteFile {
		os.Setenv("NMON2INFLUXDB_DASHBOARD_TO_FILE", "true")
	}

	if config.ImportSkipDisks {
		os.Setenv("NMON2INFLUXDB_SKIP_DISKS", "true")
	}

	if config.ImportAllCpus {
		os.Setenv("NMON2INFLUXDB_ADD_ALL_CPUS", "true")
	}

	if config.ImportBuildDashboard {
		os.Setenv("NMON2INFLUXDB_BUILD_DASHBOARD", "true")
	}

	if len(config.ImportSkipMetrics) > 0 {
		os.Setenv("NMON2INFLUXDB_SKIP_METRICS", config.ImportSkipMetrics)
	}
	app := cli.NewApp()
	app.Name = "nmon2influxdb"
	app.Usage = "upload NMON stats to InfluxDB database"
	app.Version = "0.8.2"
	app.Commands = []cli.Command{
		{
			Name:  "import",
			Usage: "import nmon files",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "skip_metrics",
					Usage:  "skip metrics",
					EnvVar: "NMON2INFLUXDB_SKIP_METRICS",
				},
				cli.BoolFlag{
					Name:   "nodisks,nd",
					Usage:  "skip disk metrics",
					EnvVar: "NMON2INFLUXDB_SKIP_DISKS",
				},
				cli.BoolFlag{
					Name:   "cpus,c",
					Usage:  "add per cpu metrics",
					EnvVar: "NMON2INFLUXDB_ADD_ALL_CPU",
				},
				cli.BoolFlag{
					Name:   "build,b",
					Usage:  "build dashboard",
					EnvVar: "NMON2INFLUXDB_BUILD_DASHBOARD",
				},
			},
			Action: NmonImport,
		},
		{
			Name:  "dashboard",
			Usage: "generate a dashboard from a nmon file or template",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "template,t",
					Usage: "optional json template file to use",
				},
				cli.BoolFlag{
					Name:   "file,f",
					Usage:  "generate Grafana dashboard file",
					EnvVar: "NMON2INFLUXDB_DASHBOARD_TO_FILE",
				},
				cli.StringFlag{
					Name:  "guser",
					Usage: "grafana user",
					Value: config.GrafanaUser,
				},
				cli.StringFlag{
					Name:  "gpassword,gpass",
					Usage: "grafana password",
					Value: config.GrafanaPassword,
				},
				cli.StringFlag{
					Name:  "gaccess",
					Usage: "grafana datasource access mode : direct or proxy",
					Value: config.GrafanaAccess,
				},
				cli.StringFlag{
					Name:  "gurl",
					Usage: "grafana url",
					Value: config.GrafanaUrl,
				},
				cli.StringFlag{
					Name:  "datasource",
					Usage: "grafana datasource",
					Value: config.GrafanaDatasource,
				},
			},
			Action: NmonDashboard,
		},
		{
			Name:  "stats",
			Usage: "generate stats from a InfluxDB metric",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "metric,m",
					Usage: "mandatory metric for stats",
				},
				cli.StringFlag{
					Name:  "statshost,s",
					Usage: "host metrics",
					Value: config.StatsHost,
				},
				cli.StringFlag{
					Name:  "from,f",
					Usage: "from date",
					Value: config.StatsFrom,
				},
				cli.StringFlag{
					Name:  "to,t",
					Usage: "to date",
					Value: config.StatsTo,
				},
				cli.StringFlag{
					Name:  "sort",
					Usage: "field to perform sort",
					Value: config.StatsSort,
				},
				cli.IntFlag{
					Name:  "limit,l",
					Usage: "limit the output",
					Value: config.StatsLimit,
				},
				cli.StringFlag{
					Name:  "filter",
					Usage: "specify a filter on fields",
					Value: config.StatsFilter,
				},
			},
			Action: NmonStat,
		},
		{
			Name:  "list",
			Usage: "list InfluxDB metrics or measurements",
			Subcommands: []cli.Command{
				{
					Name:  "measurement",
					Usage: "list InfluxDB measurements",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "host",
							Usage: "only for specified host",
						},
						cli.StringFlag{
							Name:  "filter,f",
							Usage: "filter measurement",
						},
					},
					Action: NmonListMeasurement,
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "server,s",
			Usage: "InfluxDB server and port",
			Value: config.InfluxdbServer,
		},
		cli.StringFlag{
			Name:  "port,p",
			Usage: "InfluxDB port",
			Value: config.InfluxdbPort,
		},
		cli.StringFlag{
			Name:  "db,d",
			Usage: "InfluxDB database",
			Value: config.InfluxdbDatabase,
		},
		cli.StringFlag{
			Name:  "user,u",
			Usage: "InfluxDB administrator user",
			Value: config.InfluxdbUser,
		},
		cli.StringFlag{
			Name:  "pass",
			Usage: "InfluxDB administrator pass",
			Value: config.InfluxdbPassword,
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug mode",
			EnvVar: "NMON2INFLUXDB_DEBUG",
		},
		cli.StringFlag{
			Name:  "tz,t",
			Usage: "timezone",
			Value: config.Timezone,
		},
	}
	app.Author = "Alain Dejoux"
	app.Email = "adejoux@djouxtech.net"
	app.Run(os.Args)

}
