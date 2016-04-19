# nmon2influxdb


This application take a nmon file and upload it in a [InfluxDB](influxdb.com) database.
It generates also a dashboard to allow data visualization in [Grafana](http://grafana.org/).
It's working on linux only for now.

# Demo

A live demo is available at : [demo.nmon2influxdb.org](http://demo.nmon2influxdb.org)

user/password: demo/demo

It's a read only editor user. You can change anything but cannot save it.

# Download

Go to my [github repository Releases section](https://github.com/adejoux/nmon2influxdb/releases)

You can build the binary from source. You need to have a working GO environment, see [golang.org installation instructions](https://golang.org/doc/install). Check GOPATH environment variable to be sure it's defined.

~~~
go get -u github.com/adejoux/nmon2influxdb
cd $GOPATH/src/github.com/adejoux/nmon2influxdb
go build
~~~

**[FULL Documentation available here](https://nmon2influxdb.org)**


Copyright
==========

The code is licensed as GNU AGPLv3. See the LICENSE file for the full license.

Copyright (c) 2014 Alain Dejoux <adejoux@djouxtech.net>
