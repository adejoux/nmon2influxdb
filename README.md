nmon2influx
===========

This application take a nmon file and upload it in a [InfluxDB](influxdb.com) database.
It generates also a dashboard to allow data visualization in [Grafana](http://grafana.org/).
It's working on linux only for now.

Download
========

I made available a compiled version for linux x86_64 on Dropbox : https://www.dropbox.com/s/zylemeaas1w2o2g/nmon2influx

Else you can download the git repository and build the binary from source(you need to have a working GO environment):

~~~
git clone https://github.com/adejoux/nmon2influx
cd nmon2influx
go build
~~~

nmon2influx will upload all nmon data in influxDB in series.

If you import a nmon file for server **testsrv**, it will create series starting by the server hostname.

For example, the data in CPU_ALL will be available in testsrv_CPU_ALL.

You can list all series stored in InfluxDB by executing "list series" command.

InfluxDB and Grafana
====================

You need a working InfluxDB database to use **nmon2influx**.

You can use also my influx_grafana docker image :

    # docker pull adejoux/influxdb_grafana

To start the container :

    # docker run -d -p 80:80 -p 8083:8083 -p 8086:8086 --name="nmon_reports3" -t adejoux/influxdb_grafana

The git repository of this image is available at : [docker_influxdb_grafana](https://github.com/adejoux/docker_influxdb_grafana)

Grafana will be available at url : [http://localhost/grafana](http://localhost/grafana)

InfluxDB administration interface will be available at : [http://localhost:8083](http://localhost:8083)


Usage
=======

The different parameters are :
  * -file="nmonfile": nmon file
  *  -nodashboard=false: only upload data
  *  -nodata=false: generate dashboard only
  *  -nodisk=false: skip disk metrics
  *  -pass="admin": influxdb administor password
  *  -host="localhost:8086": influxdb server and port
  *  -tmplfile="tmplfile": grafana dashboard template
  *  -user="admin": influxdb administor user

Examples
========

Importing a nmon file :

~~~
# nmon2influx -file=testsrv_141114_0000.nmon
administrator not defined. Creating user admin with password admin
Creating database : nmon_reports
Creating database : grafana
Writing GRAFANA dashboard: testsrv_dashboard
~~~

Importing a nmon file without the disk data :
~~~
# nmon2influx -nodisk -file=testsrv_141114_0000.nmon
Writing GRAFANA dashboard: testsrv_dashboard
~~~

Do not generate the grafana dashboard :

~~~
# nmon2influx -nodashboard -file=testsrv_141114_0000.nmon
~~~

Generating only the grafana dashboard :

~~~
# nmon2influx -nodata -file=testsrv_141114_0000.nmon
Writing GRAFANA dashboard: testsrv_dashboard
~~~

Generating only a grafana dashboard but using a specific template :

~~~
# nmon2influx -nodata -file=testsrv_141114_0000.nmon -tmplfile=grafana.json.tmpl
Writing GRAFANA dashboard: testsrv_dashboard
~~~

template
========

**nmon2influx** use a default template internally. It's the one you find in the file **grafana.json.tmpl**.

It's possible to override it. Grafana dashboard are defined in json.
nmon2influx use test/template to extend the json file.
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
