---
date: 2016-04-15T14:07:45+02:00
title: linux
menu:
  main:
    parent: Getting started
    identifier: /getting-started/linux
    weight: 10
---

# Redhat and centos

InfluxDB installation:

{{< highlight batch >}}
wget  https://s3.amazonaws.com/influxdb/influxdb-0.12.1-1.x86_64.rpm
yum install influxdb-0.12.1-1.x86_64.rpm
{{< /highlight >}}

Grafana installation:

{{< highlight batch >}}
wget https://grafanarel.s3.amazonaws.com/builds/grafana-3.0.0-beta41460581169.x86_64.rpm
yum install grafana-3.0.0-beta41460581169.x86_64.rpm
{{< /highlight >}}

nmon2influxdb installation:

{{< highlight batch >}}
wget https://github.com/adejoux/nmon2influxdb/releases/download/v0.9.0/nmon2influxdb-linux-amd64.gz
gunzip nmon2influxdb-linux-amd64.gz
mv nmon2influxdb-linux-amd64 nmon2influxdb
{{< /highlight >}}

# Ubuntu and Debian

InfluxDB installation:

{{< highlight batch >}}
wget  https://s3.amazonaws.com/influxdb/influxdb_0.12.1-1_amd64.deb
dpkg -i influxdb-0.12.1-1.x86_64.rpm
{{< /highlight >}}

Grafana installation:

{{< highlight batch >}}
wget https://grafanarel.s3.amazonaws.com/builds/grafana_3.0.0-beta41460581169_amd64.deb
dpkg -i grafana-3.0.0-beta41460581169.x86_64.rpm
{{< /highlight >}}

nmon2influxdb installation:

{{< highlight batch >}}
wget https://github.com/adejoux/nmon2influxdb/releases/download/v0.9.0/nmon2influxdb-linux-amd64.gz
gunzip nmon2influxdb-linux-amd64.gz
mv nmon2influxdb-linux-amd64 nmon2influxdb
{{< /highlight >}}

# Examples

Grafana will be available at url : [http://localhost:3000](http://localhost:3000)

InfluxDB administration interface will be available at : [http://localhost:8083](http://localhost:8083)

Loading a nmon file:

{{< highlight batch >}}
nmon2influxdb import server.nmon
{{< /highlight >}}

Creating a dashboard:
{{< highlight batch >}}
nmon2influxdb dashboard server.nmon
{{< /highlight >}}


Sample nmon reports are available [here](https://github.com/adejoux/nmon2influxdb/releases/download/v0.6.0/nmon_samples.tar.gz).
