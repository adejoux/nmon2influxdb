// nmon2influxdb
// import nmon data in InfluxDB
// author: adejoux@djouxtech.net

package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {

	app := cli.NewApp()
	app.Name = "nmon2influxdb"
	app.Usage = "upload NMON stats to InfluxDB database"
	app.Version = "0.4.0"
	app.Commands = []cli.Command{
		{
			Name:  "import",
			Usage: "import a nmon file",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "nodisks,nd",
					Usage: "add disk metrics",
				},
				cli.BoolFlag{
					Name:  "cpus,c",
					Usage: "add per cpu metrics",
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
					Usage: "optional template file to use",
				},
				cli.StringFlag{
					Name:  "guser",
					Usage: "grafana user",
					Value: "admin",
				},
				cli.StringFlag{
					Name:  "gpassword,gpass",
					Usage: "grafana password",
					Value: "admin",
				},
				cli.StringFlag{
					Name:  "gurl",
					Usage: "grafana url",
					Value: "http://localhost:3000",
				},
				cli.StringFlag{
					Name:  "datasource",
					Usage: "grafana datasource",
					Value: "nmon2influxdb",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:  "file",
					Usage: "generate a dashboard from a nmon file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "template,t",
							Usage: "optional json template file to use",
						},
						cli.BoolFlag{
							Name:  "file,f",
							Usage: "generate Grafana dashboard file",
						},
						cli.StringFlag{
							Name:  "guser",
							Usage: "grafana user",
							Value: "admin",
						},
						cli.StringFlag{
							Name:  "gpassword,gpass",
							Usage: "grafana password",
							Value: "admin",
						},
						cli.StringFlag{
							Name:  "gurl",
							Usage: "grafana url",
							Value: "http://localhost:3000",
						},
						cli.StringFlag{
							Name:  "datasource",
							Usage: "grafana datasource",
							Value: "nmon2influxdb",
						},
					},
					Action: NmonDashboardFile,
				},
				{
					Name:  "template",
					Usage: "generate a dashboard from a TOML template",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "template,t",
							Usage: "optional json template file to use",
						},
						cli.BoolFlag{
							Name:  "file,f",
							Usage: "generate Grafana dashboard file",
						},
						cli.StringFlag{
							Name:  "guser",
							Usage: "grafana user",
							Value: "admin",
						},
						cli.StringFlag{
							Name:  "gpassword,gpass",
							Usage: "grafana password",
							Value: "admin",
						},
						cli.StringFlag{
							Name:  "gurl",
							Usage: "grafana url",
							Value: "http://localhost:3000",
						},
						cli.StringFlag{
							Name:  "datasource",
							Usage: "grafana datasource",
							Value: "nmon2influxdb",
						},
					},
					Action: NmonDashboardTemplate,
				},
			},
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
				},
				cli.StringFlag{
					Name:  "from,f",
					Usage: "from date",
				},
				cli.StringFlag{
					Name:  "to,t",
					Usage: "from date",
				},
				cli.StringFlag{
					Name:  "sort",
					Usage: "field to perform sort",
					Value: "mean",
				},
				cli.IntFlag{
					Name:  "limit,l",
					Usage: "limit the output",
				},
			},
			Action: NmonStat,
		},
		{
			Name:   "list",
			Usage:  "list InfluxDB metrics",
			Action: NmonList,
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "server,s",
			Value: "localhost",
			Usage: "InfluxDB server and port",
		},
		cli.StringFlag{
			Name:  "port",
			Value: "8086",
			Usage: "InfluxDB port",
		},
		cli.StringFlag{
			Name:  "db,d",
			Value: "nmon_reports",
			Usage: "InfluxDB database",
		},
		cli.StringFlag{
			Name:  "user,u",
			Value: "root",
			Usage: "InfluxDB administrator user",
		},
		cli.StringFlag{
			Name:  "pass,p",
			Value: "root",
			Usage: "InfluxDB administrator pass",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "debug mode",
		},
		cli.StringFlag{
			Name:  "tz,t",
			Usage: "timezone",
		},
	}
	app.Author = "Alain Dejoux"
	app.Email = "adejoux@djouxtech.net"
	app.Run(os.Args)

}
