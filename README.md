nmon2influxdb
===========

This application take a nmon file and upload it in a [InfluxDB](influxdb.com) database.
It generates also a dashboard to allow data visualization in [Grafana](http://grafana.org/).
It's working on linux only for now.

Download
========

I made available a compiled version for linux x86_64 on Dropbox : https://www.dropbox.com/s/3keh5d6umnvcwv8/nmon2influxdb

Else you can download the git repository and build the binary from source(you need to have a working GO environment):

~~~
git clone https://github.com/adejoux/nmon2influxdb
cd nmon2influxdb
go build
~~~

nmon2influxdb will upload all nmon data in InfluxDB in series.

If you import a nmon file for server **testsrv**, it will create series starting by the server hostname.

For example, the data in CPU_ALL will be available in testsrv_CPU_ALL.

You can list all series stored in InfluxDB by executing "list series" command.

InfluxDB and Grafana
====================

You need a working InfluxDB database to use **nmon2influxdb**.

You can use also my influxdb_grafana docker image :

    # docker pull adejoux/docker-influxdb-grafana

To start the container :

    # docker run -d -p 3000:3000 -p 8083:8083 -p 8086:8086 --name="nmon_reports" -t adejoux/docker-influxdb-grafana

The git repository of this image is available at : [docker_influxdb_grafana](https://github.com/adejoux/docker_influxdb_grafana)

Grafana will be available at url : [http://localhost:3000](http://localhost:3000)

InfluxDB administration interface will be available at : [http://localhost:8083](http://localhost:8083)


Usage
=======

~~~
nmon2influxdb
NAME:
   nmon2influxdb - upload NMON stats to InfluxDB database

USAGE:
   nmon2influxdb [global options] command [command options] [arguments...]

VERSION:
   0.4.0

AUTHOR(S):
   Alain Dejoux <adejoux@djouxtech.net>

COMMANDS:
   import import a nmon file
   dashboard  generate a dashboard from a nmon file
   stats generate stats from a InfluxDB metric
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --server, -s "localhost" InfluxDB server and port
   --port "8086"    InfluxDB port
   --db, -d "nmon_reports"  InfluxDB database
   --user, -u "root"    InfluxDB administrator user
   --pass, -p "root"    InfluxDB administrator pass
   --debug      debug mode
   --help, -h     show help
   --version, -v    print the version
~~~

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

~~~
nmon2influxdb dashboard -h
NAME:
   dashboard - generate a dashboard from a nmon file

USAGE:
   command dashboard [command options] [arguments...]

OPTIONS:
   --template, -t       optional template file to use
   --file, -f       generate Grafana dashboard file
   --guser "admin"      grafana user
   --gpassword, --gpass "admin"   grafana password
   --gurl "http://localhost:3000" grafana url
   --datasource "nmon2influxdb"   grafana datasource
~~~

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
   --to, -t     from date
   --aggregate, -a  aggregate function
~~~

Examples
========

Importing a nmon file :

~~~
# nmon2influxdb import testsrv_141114_0000.nmon
File testsrv_141114_0000.nmon imported !
~~~

Upload a dashboard to Grafana :
~~~
nmon2influxdb dashboard -f testsrv_141114_0000.nmon
Writing GRAFANA dashboard: testsrv_dashboard
~~~

Generate a dashboard file from the NMON file :
~~~
nmon2influxdb dashboard -f testsrv_141114_0000.nmon
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


template
========

**nmon2influxdb** use a default template internally. It's the one you find in the file **grafana.json.tmpl**.

It's possible to override it. Grafana dashboard are defined in json.
nmon2influxdb use test/template to extend the json file.
For now, it's still not evolved and is mainly standard GO.

* **.GetColumns "NMONCATEGORY"** : will return every columns registered in influxdb for [hostname]_[NMON CATEGORY]

A example :

~~~
{{ range $index, $adapter := .GetColumns "NET" }}{{ if $index}},{{end}}
    {
        "function": "mean",
        "column": "{{$adapter}}",
        "series": "{{$.Hostname}}_NET",
        "query": "select mean(\"{{$adapter}}\") from \"{{$.Hostname}}_NET\" where $timeFilter group by time($interval) order asc",
        "rawQuery": true,
        "alias": "{{$adapter}}"
    }{{end}}
~~~

This will list all available columns in the NET serie of the server.

A example output for **testsrv**:

~~~
{
    "function": "mean",
    "column": "en0-read-KB/s",
    "series": "testsrv_NET",
    "query": "select mean(\"en0-read-KB/s\") from \"testsrv_NET\" where $timeFilter group by time($interval) order asc",
    "rawQuery": true,
    "alias": "en0-read-KB/s"
},
{
    "function": "mean",
    "column": "en1-read-KB/s",
    "series": "testsrv_NET",
    "query": "select mean(\"en1-read-KB/s\") from \"testsrv_NET\" where $timeFilter group by time($interval) order asc",
    "rawQuery": true,
    "alias": "en1-read-KB/s"
},
~~~



* **.GetFilteredColumns "NMONCATEGORY" "FILTER"** : will return columns matching FILTER in [hostname]_[NMON CATEGORY].

A example :

~~~
{{ range $index, $adapter := .GetFilteredColumns "NPIV" "e-KB"}}{{ if $index}},{{end}}
    {
      "function": "mean",
      "column": "{{$adapter}}",
      "series": "{{$.Hostname}}_NPIV",
      "query": "select mean(\"{{$adapter}}\") from \"{{$.Hostname}}_NPIV\" where $timeFilter group by time($interval) order asc",
      "rawQuery": true,
      "alias": "{{$adapter}}"
    }{{end}}
~~~

And the output where only columns matching "e-KB" are listed:

~~~
{
  "function": "mean",
  "column": "vfchost0_write-KB/s",
  "series": "testsrv_NPIV",
  "query": "select mean(\"vfchost0_write-KB/s\") from \"testsrv_NPIV\" where $timeFilter group by time($interval) order asc",
  "rawQuery": true,
  "alias": "vfchost0_write-KB/s"
},
{
  "function": "mean",
  "column": "vfchost2_write-KB/s",
  "series": "testsrv_NPIV",
  "query": "select mean(\"vfchost2_write-KB/s\") from \"testsrv_NPIV\" where $timeFilter group by time($interval) order asc",
  "rawQuery": true,
  "alias": "vfchost2_write-KB/s"
},
~~~


**Important**: you need to always include $adapter between quote else InfluxDB query will not work if some special characters in the columns name(like above).



Copyright
==========

The code is licensed as GNU AGPLv3. See the LICENSE file for the full license.

Copyright (c) 2014 Alain Dejoux <adejoux@djouxtech.net>
