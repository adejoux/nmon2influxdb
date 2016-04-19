---
date: 2016-04-15T16:27:34+02:00
title: dashboard
menu:
  main:
    parent: Usage
    identifier: /usage/dashboard
    weight: 20
---

{{< highlight batch >}}
NAME:
   nmon2influxdb dashboard - generate a dashboard from a nmon file or template

USAGE:
   nmon2influxdb dashboard [command options] [arguments...]

OPTIONS:
   --template, -t 		optional json template file to use
   --file, -f			generate Grafana dashboard file [$NMON2INFLUXDB_DASHBOARD_TO_FILE]
   --guser "admin"		grafana user
   --gpassword, --gpass "admin"	grafana password
   --gaccess "direct"		grafana datasource access mode : direct or proxy
   --gurl "http://uby:3000"	grafana url
   --datasource "nmon2influxdb"	grafana datasource
{{< /highlight >}}

# Parameters

  * **template**: specify a json grafana template to upload
  * **file**: generate a json file instead of uploading the dashboard to the grafana server
  * **guser**: grafana user. By default, it's **admin**.
  * **gpassword**: grafana password. By default, it's **admin**.
  * **gaccess**: Specify how grafana will connect to the InfluxDB database. See the section **database connection**
  * **gurl**: grafana server url.
  * **datasource**: datasource used by grafana dashboard

# Environment variables

Environment variables can be specified to setup default parameter values.

  * **NMON2INFLUXDB_DASHBOARD_TO_FILE**

# database connection

Grafana has two way to access the InfluxDB database:

* **direct**: the connection is done directly from the user's browser. Its possible when the database is directly accessible.

* **proxy**: the connection is made by the Grafana server itself.

# Examples

Upload a dashboard to Grafana based on a NMON file :
{{< highlight batch >}}
nmon2influxdb dashboard testsrv_141114_0000.nmon
Dashboard uploaded to grafana
{{< /highlight >}}

Generate a dashboard file from the NMON file :
{{< highlight batch >}}
nmon2influxdb dashboard -f testsrv_141114_0000.nmon
Writing GRAFANA dashboard: testsrv_dashboard
{{< /highlight >}}
