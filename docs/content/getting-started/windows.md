---
date: 2016-04-15T14:27:46+02:00
title: windows
menu:
  main:
    parent: Getting started
    identifier: /getting-started/windows
    weight: 30
---

# docker-toolbox installation

**docker-toolbox** is packaging Docker to make it available easily on Windows.

A simple installer is available on  [docker-toolbox](https://www.docker.com/toolbox)

{{< gallery image="docker_toolbox_link.png" >}}

After downloading it, run DockerToolbox-<version>.exe, and you will see this installer windows :

{{< gallery image="docker_toolbox_installer.png" >}}

It will handle almost everything. You can be asked for Administrator privileges during the installation.

{{< gallery image="docker_toolbox_installer2.png" >}}

After some time, **kitematic** will launch and start to install VirtualBox :

{{< gallery image="kitematic_install.png" >}}

Admin permissions can be requested to install.

After some time, **kitematic** should have installed and started the **boot2docker** VM and you will see this screen :
{{< gallery image="kitematic.png" >}}



# Grafana and InfluxDB container installation



A docker container with InfluxDB and Grafana using automated build is available on docker hub :
{{< gallery image="dockerhub_image.png" >}}

The Dockerfile is pretty simple:
{{< gallery image="dockerhub_dockerfile.png" >}}

Install it directly from **kitematic**.

And you will have a container running both InfluxDB and Grafana :
{{< gallery image="container_started.png" >}}

By clicking on web preview, you will access Grafana login page. Default account is **admin**, password is **admin**.

**kitematic** will map the container's ports randomly.

{{< gallery image="ip_vm.png" >}}

It's only performed at the container's creation and it will not change during the container's life.
Else in the boot2docker VM in VirtualBox named **default**, use docker command line to delete the container and start a new container with the new port mappings :

{{< highlight batch >}}
docker ps
docker rm -f docker-influxdb-grafana
docker run -d -p 3000:3000 -p 8083:8083 -p 8086:8086 --name="nmon_reports" -t adejoux/docker-influxdb-grafana
{{< /highlight >}}

{{< gallery image="docker_cli.png" >}}

kitematic will be able to manage the new container without problem.

# Using nmon2influxdb

**nmon2influxdb binary** file can be used directly from windows command line tool.

The main difference with running it on linux is you will need to specify the url to access Grafana and the server and port to access the InfluxDB database.

For example, importing a nmon file will be :
{{< highlight batch >}}
nmon2influxdb-windows-amd64.exe -s 192.168.99.100 import vios1_150706_1552.nmon
{{< /highlight >}}

For dashboard, you will need both InfluxDB and Grafana parameters the first time. It's because the tool will create the DataSource if it's not existing :
{{< highlight batch >}}
nmon2influxdb-windows-amd64.exe dashboard file --gurl http://192.168.99.100:3000 vios1_150706_1552.nmon
{{< /highlight >}}

{{< gallery image="nmon2influxdb_win.png" >}}

It's also possible to setup the [configuration file](/configuration)

Sample nmon reports are available [here](https://github.com/adejoux/nmon2influxdb/releases/download/v0.6.0/nmon_samples.tar.gz).
