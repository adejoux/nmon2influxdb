# nmon2influxdb


This application take a nmon file and upload it in a [InfluxDB](influxdb.com) database.
It generates also a dashboard to allow data visualization in [Grafana](http://grafana.org/).
It's working on linux only for now.

# Download

Go to my [github repository Releases section](https://github.com/adejoux/nmon2influxdb/releases)

Else you can download the git repository and build the binary from source(you need to have a working GO environment):

~~~
git clone https://github.com/adejoux/nmon2influxdb
cd nmon2influxdb
go build
~~~

# InfluxDB and Grafana

You need a working InfluxDB database to use **nmon2influxdb**.

Both InfluxDB and Grafana are easy to install on linux. On Windows or Mac OS it's doable but more complicated. Using a Docker container is easier.

You can use my docker-influxdb-grafana docker image :

    # docker pull adejoux/docker-influxdb-grafana

To start the container :

    # docker run -d -p 3000:3000 -p 8083:8083 -p 8086:8086 --name="nmon_reports" -t adejoux/docker-influxdb-grafana

The git repository of this image is available at : [docker_influxdb_grafana](https://github.com/adejoux/docker_influxdb_grafana)

## On Linux


Grafana will be available at url : [http://localhost:3000](http://localhost:3000)

InfluxDB administration interface will be available at : [http://localhost:8083](http://localhost:8083)

## On Windows/Mac OS

You will need to know the IP address of your boot2docker VM.
You will need to use the options in command line. For example, if my VM has ip **192.168.99.100** :
~~~
nmon2influxdb -s 192.168.99.100 import test.nmon
~~~

For dashboard, you will need to specify the grafana url too :
~~~
nmon2influxdb dashboard file --gurl "http://192.168.99.100:3000" test.nmon
~~~



Usage
=======


~~~
#nmon2influxdb
NAME:
   nmon2influxdb - upload NMON stats to InfluxDB database

USAGE:
   nmon2influxdb [global options] command [command options] [arguments...]

VERSION:
   0.7.0

AUTHOR(S):
   Alain Dejoux <adejoux@djouxtech.net>

COMMANDS:
   import import a nmon file
   dashboard  generate a dashboard from a nmon file or template
   stats  generate stats from a InfluxDB metric
   list   list InfluxDB metrics or measurements
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --server, -s "localhost" InfluxDB server and port
   --port, -p "8086"    InfluxDB port
   --db, -d "nmon_reports"  InfluxDB database
   --user, -u "root"    InfluxDB administrator user
   --pass "root"    InfluxDB administrator pass
   --debug      debug mode
   --tz, -t       timezone
   --help, -h     show help
   --version, -v    print the version
~~~

## import command
The **import** command is the one used to import NMON files in InfluxDB.

You can skip disks metrics with **nd** option. Can be interesting when you have a system with a lot of disks and where you are interested only in CPU performance.
~~~
nmon2influxdb import -h
NAME:
   import - import a nmon file

USAGE:
   command import [command options] [arguments...]

OPTIONS:
   --nodisks, --nd  add disk metrics
   --cpus, -c   add per cpu metrics
~~~

## dashboard

The **dashboard** command has two subcommands to create dashboards in Grafana.
~~~
nmon2influxdb dashboard -h
NAME:
   nmon2influxdb dashboard - generate a dashboard from a nmon file or template

USAGE:
   nmon2influxdb dashboard command [command options] [arguments...]

COMMANDS:
   file   generate a dashboard from a nmon file
   template generate a dashboard from a TOML template
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --help, -h show help
~~~

### dashboard file

The **dashboard file** subcommand will create and upload a dashboard from a NMON file.

It's great as a starting point with Grafana.

~~~
NAME:
   file - generate a dashboard from a nmon file

USAGE:
   command file [command options] [arguments...]

OPTIONS:
   --template, -t       optional json template file to use
   --file, -f       generate Grafana dashboard file
   --guser "admin"      grafana user
   --gpassword, --gpass "admin"   grafana password
   --gurl "http://localhost:3000" grafana url
   --datasource "nmon2influxdb"   grafana datasource
~~~

### dashboard template

The **dashboard template** subcommand will create and upload a dashboard from a TOML template file.

~~~
NAME:
   template - generate a dashboard from a TOML template

USAGE:
   command template [command options] [arguments...]

OPTIONS:
   --template, -t       optional json template file to use
   --file, -f       generate Grafana dashboard file
   --guser "admin"      grafana user
   --gpassword, --gpass "admin"   grafana password
   --gurl "http://localhost:3000" grafana url
   --datasource "nmon2influxdb"   grafana datasource
~~~

## stats

Dashboards are not always the best way to analyze data. When you have hundred or thousands disks it's better to have a summary of stats.

~~~
nmon2influxdb stats -h
NAME:
   stats - generate stats from a InfluxDB metric

USAGE:
   command stats [command options] [arguments...]

OPTIONS:
   --metric, -m   mandatory metric for stats
   --statshost, -s  host metrics
   --from, -f     from date
   --to, -t     to date
   --sort "mean"  field to perform sort
   --limit, -l "0"  limit the output
   --filter     specify a filter on fields
~~~



TOML dashboard templates
========================

Grafana is really great to build templates but you need to learn it to build advanced dashboard. It's pretty easy but I wanted to provide ready to use templates. JSON is pretty hard to read for humans so I implemented TOML templates.

For example, this template will display informations coming from two differents VIO servers :

~~~ toml
title = "templated dual vio"
[templates]
  [[templates.template]]
    # name of the template variable
    name = "vios"
    # the default nmon2influxdb database in InfluxDB
    datasource = "nmon2influxdb"
    # query used to retrieve values in InfluxDB
    query =  "SHOW TAG VALUES FROM CPU_ALL WITH KEY = host"
    # Regular expression used to filter the query result
    regex = "vios"
[[row]]
title = "LPAR"
    [[row.panel]]
    title = "shared processor"
        [[row.panel.metric]]
            measurement = "LPAR"
            # use template variable $vios
            hosts = ["$vios"]
            fields = ["PhysicalCPU", "virtualCPUs", "entitled"]
        [[row.panel.override]]
            alias = "$tag_host PhysicalCPU"
            stack = true
            fill = 1
        [row.panel.tooltip]
            value_type = "individual"

    [[row.panel]]
        title = "shared processor pool"
        [[row.panel.metric]]
            measurement = "LPAR"
            # use template variable $vios
            hosts = ["$vios"]
            fields = ["PoolIdle", "poolCPUs"]

[[row]]
title = "SEA"
    [[row.panel]]
    title = "SEA WRITE"
    stack = true
    fill = 1
        [[row.panel.metric]]
            measurement = "SEA"
            # use template variable $vios
            hosts = ["$vios"]
            # regular expressions are used on fields
            fields = ["write-KB"]
    [[row.panel]]
    title = "SEA READ"
    stack = true
    fill = 1
        [[row.panel.metric]]
            measurement = "SEA"
            # use template variable $vios
            hosts = ["$vios"]
            # regular expressions are used on fields
            fields = ["read-KB"]
[time]
  from = "2015-07-06T13:03:56.000Z"
  to = "2015-07-07T14:23:44.000Z"
~~~

Examples
========

Importing a nmon file :

~~~
# nmon2influxdb import testsrv_141114_0000.nmon
##################################################################################
File testsrv_141114_0000.nmon imported !
~~~

Each '#' character in the output means 10 000 points was inserted in InfluxDB.

Upload a dashboard to Grafana based on a NMON file :
~~~
nmon2influxdb dashboard file testsrv_141114_0000.nmon
Dashboard uploaded to grafana
~~~

Generate a dashboard file from the NMON file :
~~~
nmon2influxdb dashboard file -f testsrv_141114_0000.nmon
Writing GRAFANA dashboard: testsrv_dashboard
~~~

Importing a nmon file without the disk data :
~~~
# nmon2influx import -nodisks testsrv_141114_0000.nmon
Writing GRAFANA dashboard: testsrv_dashboard
~~~

Generating stats for DISKREADSERV metric :
~~~
nmon2influxdb stats -m DISKREADSERV -s lpar1
          field|     Min|    Mean|  Median|     Max|Points #
        hdisk10|    0.40|    2.42|    2.10|   12.90|    1200
        hdisk11|    0.60|    2.63|    2.20|   14.00|    1200
        hdisk12|    0.50|    2.74|    2.30|   16.30|    1200
        hdisk13|    0.00|    0.52|    0.00|   16.30|    1200
         hdisk7|    0.00|    0.01|    0.00|    0.80|    1200
         hdisk8|    0.00|    0.04|    0.00|    1.00|    1200
         hdisk9|    0.50|    2.33|    2.10|   11.30|    1200
~~~

Generating stats for **DISKREADSERV** metric on host **lpar1** on a time period :
~~~
nmon2influxdb stats -m DISKREADSERV -s lpar1 -f "08:04:05,29-May-2015" -t "15:04:05,29-May-2015"
          field|     Min|    Mean|  Median|     Max|Points #
        hdisk10|    2.20|    3.80|    3.65|    5.60|     132
        hdisk11|    2.40|    4.10|    3.90|    5.80|     132
        hdisk12|    2.30|    4.39|    4.05|    6.50|     132
        hdisk13|    0.00|    0.73|    0.00|   13.90|     132
         hdisk7|    0.00|    0.02|    0.00|    0.20|     132
         hdisk8|    0.00|    0.10|    0.00|    1.00|     132
         hdisk9|    2.10|    3.68|    3.55|    5.20|     132
~~~

Generating stats for **DISKREADSERV** metric on host **lpar1** and limit output to the 5 most active disk :

~~~
nmon2influxdb stats -m DISKREADSERV -s lpar1 -l 5
          field|     Min|    Mean|  Median|     Max|Points #
        hdisk10|    2.20|    3.80|    3.65|    5.60|     132
        hdisk11|    2.40|    4.10|    3.90|    5.80|     132
        hdisk12|    2.30|    4.39|    4.05|    6.50|     132
        hdisk13|    0.00|    0.73|    0.00|   13.90|     132
         hdisk7|    0.00|    0.02|    0.00|    0.20|     132
~~~

Listing measurements available for a specific host and containing DISK in the name :

~~~
nmon2influxdb list measurement --host feddy --filter DISK
measurements
DISKBSIZE
DISKBUSY
DISKREAD
DISKWRITE
DISKXFER
~~~


Copyright
==========

The code is licensed as GNU AGPLv3. See the LICENSE file for the full license.

Copyright (c) 2014 Alain Dejoux <adejoux@djouxtech.net>
