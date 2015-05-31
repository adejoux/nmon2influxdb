// nmon2influx
// import nmon report in Influxdb
//version: 0.1
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
	app.Version = "0.2.0"
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
			Usage: "generate a dashboard from a nmon file",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "template,t",
					Usage: "optional template file to use",
				},
			},
			Action: NmonDashboard,
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
	}
	app.Author = "Alain Dejoux"
	app.Email = "adejoux@djouxtech.net"
	app.Run(os.Args)

}
