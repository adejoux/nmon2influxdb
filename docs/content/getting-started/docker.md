---
date: 2016-04-15T14:18:33+02:00
title: docker
menu:
  main:
    parent: Getting started
    identifier: /getting-started/docker
    weight: 20
---

It's possible to use a Docker container with **InfluxDB** and **Grafana** inside.
It allows to have a dedicated **InfluxDB** database for each performance analysis by spinning a new instance.

* Install Docker on Ubuntu:
{{< highlight batch >}}
sudo apt-get install -y docker.io
{{< /highlight >}}

* **Or** install docker on redhat/fedora:
{{< highlight batch >}}
sudo yum install docker
sudo service docker start
{{< /highlight >}}

* Download InfluxDB + Grafana Container :
{{< highlight batch >}}
sudo docker pull adejoux/docker-influxdb-grafana
{{< /highlight >}}

* Start the docker container :
{{< highlight batch >}}
sudo docker run -d -p 3000:3000 -p 8083:8083 -p 8086:8086 --name="test" -t adejoux/docker-influxdb-grafana
{{< /highlight >}}

* Download nmon2influxdb :
{{< highlight batch >}}
wget https://github.com/adejoux/nmon2influxdb/releases/download/v0.9.0/nmon2influxdb-linux-amd64.gz
gunzip nmon2influxdb-linux-amd64.gz
mv nmon2influxdb-linux-amd64.gz nmon2influxdb
chmod u+x nmon2influxdb
{{< /highlight >}}

# Examples

Grafana will be available at url : **http://[your vm ip]:3000**

InfluxDB administration interface will be available at : **http://[your vm ip]:8083**

  * You can configure hostnames directly in the [configuration file](/configuration/file/):
{{< highlight toml >}}
influxdb_server="yourvm"
influxdb_port="8086"
grafana_URL="http://yourvm:3000"
{{< /highlight >}}

  * Loading a nmon file:
{{< highlight batch >}}
nmon2influxdb import server.nmon
{{< /highlight >}}

  * Creating a dashboard:
{{< highlight batch >}}
nmon2influxdb dashboard server.nmon
{{< /highlight >}}


Sample nmon reports are  [available](https://github.com/adejoux/nmon2influxdb/releases/download/v0.6.0/nmon_samples.tar.gz).
