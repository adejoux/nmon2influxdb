---
date: 2016-04-15T13:56:31+02:00
title: nmon2influxdb
type: index
---

# nmon2influxdb
## NMON data made dynamic

This application take a nmon file and upload it in a [InfluxDB](influxdata.com/time-series-platform/influxdb/) database.

It generates also a dashboard to allow data visualization in [Grafana](http://grafana.org/).

It's working on linux,Windows and Mac OS. But you will need boot2docker or manual installation to have grafana and influxdb running Windows.

# gallery
{{< gallery image="nmon2influxdb.png" >}}
{{< gallery image="nmon2influxdb2.png" addclass="hidden" >}}
{{< gallery image="nmon2influxdb3.png" addclass="hidden" >}}
{{< gallery image="nmon2influxdb4.png" addclass="hidden" >}}

# features

* Import data from NMON files in InfluxDB database.

* Generate Grafana dashboards based on NMON files data.

* Create InfluxDB datasource in Grafana automatically.

* Upload directly dashboard in Grafana.

* Provide a templating system to generate custom dashboard.

* Use Grafana to customize charts and use data from multiple NMON files.
