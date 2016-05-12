---
date: 2016-04-15T16:06:33+02:00
title: import
menu:
  main:
    parent: Usage
    identifier: /usage/import
    weight: 10
---


{{< highlight batch >}}
NAME:
   nmon2influxdb import - import nmon files

USAGE:
   nmon2influxdb import [command options] [arguments...]

OPTIONS:
   --skip_metrics 			skip metrics [$NMON2INFLUXDB_SKIP_METRICS]
   --nodisks, --nd			skip disk metrics [$NMON2INFLUXDB_SKIP_DISKS]
   --cpus, -c				add per cpu metrics [$NMON2INFLUXDB_ADD_ALL_CPU]
   --build, -b				build dashboard [$NMON2INFLUXDB_BUILD_DASHBOARD]
   --force, -f				force import [$NMON2INFLUXDB_FORCE]
   --log_database "nmon2influxdb_log"	influxdb database used to log imports
   --log_retention "1d"			import log retention
{{< /highlight >}}

# Parameters

  * **skip_metrics**: specify a list of nmon metrics to not import
  * **nodisks**: skip disk metrics import
  * **cpus**: add per cpu metrics
  * **build**: automatically build the corresponding grafana dashboard
  * **force**: force import instead of skipping if already imported
  * **log_database**: the database used to log nmon files import
  * **log_retention**: will delete import file log information after 1 day by default

# Environment variables

Environment variables can be specified to setup default parameter values.

  * **NMON2INFLUXDB_SKIP_METRICS**
  * **NMON2INFLUXDB_SKIP_DISKS**
  * **NMON2INFLUXDB_ADD_ALL_CPU**
  * **NMON2INFLUXDB_BUILD_DASHBOARD**
  * **NMON2INFLUXDB_FORCE**


# Examples

Importing nmon files:

{{< highlight batch >}}
# nmon2influxdb import testsrv_141114_0000.nmon testsrv_141115_0000.nmon
##################################################################################
File testsrv_141114_0000.nmon imported !
##################################################################################
File testsrv_141115_0000.nmon imported !
{{< /highlight >}}

It's possible to specify a directory:
{{< highlight batch >}}
# nmon2influxdb import /data/nmon/
{{< /highlight >}}

Or use shell completion:
{{< highlight batch >}}
# nmon2influxdb import /data/nmon/*nmon
{{< /highlight >}}

**Note**: by default, only files with the extension **.nmon** are imported.

Importing a nmon file without the disk data :
{{< highlight batch >}}
# nmon2influx import --nodisks testsrv_141114_0000.nmon
{{< /highlight >}}

Importing a remote nmon file:
{{< highlight batch >}}
nmon2influxdb import adxlpar1:/log/nmon/lpar02_151104_1204.nmon
###############################
File /log/nmon/lpar02_151104_1204.nmon imported : 316320 points !
{{< /highlight >}}

Importing nmon files from a remote directory using the remote user **batch**:
{{< highlight batch >}}
nmon2influxdb import adxlpar1:/log/nmon
###############################
File /log/nmon/lpar02_151104_1204.nmon imported : 316320 points !
file not changed since last import: /log/nmon/lpar01_151104_1116.nmon
file not changed since last import: /log/nmon/lpar01_151104_1200.nmon
file not changed since last import: /log/nmon/lpar02_110415.nmon
{{< /highlight >}}
