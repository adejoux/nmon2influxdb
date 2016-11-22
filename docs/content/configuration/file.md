---
date: 2016-04-18T18:09:35+02:00
title: file
menu:
  main:
    parent: Configuration
    identifier: /configuration/file
    weight: 10
---

nmon2influxdb will generate a configuration file named **$HOME/.nmon2influxdb.cfg**.

It will allow to change default configuration value in command line. Command line parameters will always have precedence over the configuration file parameters.

{{< highlight toml >}}
# general
debug = false
timezone = "Europe/Paris"

# influxdb
influxdb_user = "root"
influxdb_password = "root"
influxdb_server = "uby"
influxdb_port = "8086"
influxdb_database = "nmon_reports"

# grafana
grafana_user = "admin"
grafana_password = "admin"
grafana_access = "direct"
grafana_url = "http://uby:3000"
grafana_datasource = "nmon2influxdb"

# import
import_skip_disks = false
import_all_cpus = false
import_build_dashboard=false
import_force=false
import_skip_metrics="JFSINODE|TOP"
import_ssh_user = "batchuser"
import_ssh_key = "/home/user/.ssh/id_rsa"

# import log database
import_log_database="nmon2influxdb_log"
import_log_retention="1d"

# dashboard
dashboard_write_file = false

# HMC parameters
hmc_server="mylab"
hmc_user="hscroot"
hmc_password="abc123"
hmc_managed_system="mysystem"
hmc_database="nmon2influxdbHMC"
hmc_data_retention="40d"
{{< /highlight >}}

# Additional parameters

Some parameters are not set by default because they are changing default behavior in big ways or are not useful by default.

##stats parameters
{{< highlight toml >}}
stats_limit=20
stats_sort="mean"
stats_filter=""
stats_from=""
stats_to=""
stats_host=""
{{< /highlight >}}

If you are always querying the same host or applying the same timeframe to your queries you can setup here this values.

##data retention

By default, data are kept indefinitely in InfluxDB. It's possible to change it to have data expiration.

{{< highlight toml >}}
import_data_retention = "30d"
{{< /highlight >}}

This value is updated each time a import is done.

All data older than what are specified in the retention policy are not kept.

**Note:** it's the timestamp associated with data which matters. If you load data from one year ago and you have a retention policy of 30 days, you will not see the data.
