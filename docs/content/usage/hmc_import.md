---
date: 2016-11-21T11:24:01+01:00
title: HMC import
menu:
  main:
    parent: Usage
    identifier: /usage/hmc_import
    weight: 15
---

{{< highlight batch >}}
NAME:
   nmon2influxdb hmc import - import hmc PCM data

USAGE:
   nmon2influxdb hmc import [command options] [arguments...]

OPTIONS:
   --hmc "myhmc"		hmc server
   --hmcuser "hscroot"		hmc user
   --hmcpass "abc123"		hmc password
   --managed_system, -m 	only import this managed system
   --managed_system-only, --sys-only	skip partition metrics
   --samples "0"			import latest <value> samples
   --timeout 30				set a connection timeout
{{< /highlight >}}

# Parameters

  * **hmc**: HMC to use to fetch PCM data
  * **hmcuser**: HMC user to use for connection
  * **hmcpass**: add per cpu metrics
  * **managed_system**: fetch HMC PCM data only for this managed system
  * **--sys-only**: skip partition metrics
  * **--samples <value>**: fetch the latest <value> samples. Each sample is averaging 30 seconds.
  * **--timeout <value>**: set a connection timeout

# Environment variables

Environment variables can be specified to setup default parameter values.

  * **NMON2INFLUXDB_HMC_SERVER**
  * **NMON2INFLUXDB_HMC_USER**

**Note:** HMC password cannot be set by environment variables.

# Configuration file parameters


{{< highlight toml >}}
hmc_server="mylab"
hmc_user="hscroot"
hmc_password="abc123"
hmc_managed_system="mysystem"
hmc_database="nmon2influxdbHMC"
hmc_data_retention="40d"
hmc_samples=10
hmc_timeout=30
{{< /highlight >}}

It's possible to set all CLI parameters. It's also possible to change the InfluxDB database name with **hmc_database** and change the data retention with **hmc_data_retention**.

# Examples

Loading HMC metrics from HMC **myhmc**:

{{< highlight batch >}}
nmon2influxdb hmc import
Getting list of managed systems

p750A
managed system                     p750A:     2673 points fetched.
Partition           BCK_BCK DR #adxlpar2:     2916 points fetched.
Partition               BCK DR #adxlpar2:     2916 points fetched.
Partition                        powerVC:     2916 points fetched.

POWER8-S824A
managed system              POWER8-S824A:    59532 points fetched.
Partition                       WM-SLES1:    17958 points fetched.
Partition                 LV-PCM-Manager:     8398 points fetched.
Partition                     PowerVC-LE:     7163 points fetched.
Partition                   LVL-cluster2:     7163 points fetched.
Partition                   lvl-cluster1:     7163 points fetched.
Partition                       WM-SLES2:    18031 points fetched.
{{< /highlight >}}

Note: parameters can also be set in the configuration file **~/.nmon2influxdb.cfg**.

Loading HMC metrics from HMC **myhmc** for system **mysystem** only:

{{< highlight batch >}}
nmon2influxdb hmc import -m POWER8-S824A
Fetching latest 2 hours performance metrics. See hmc_samples parameter.
Getting list of managed systems
Skipping system: p750A
POWER8-S824A
managed system              POWER8-S824A:    59532 points fetched.
Partition                       WM-SLES1:    17958 points fetched.
Partition                 LV-PCM-Manager:     8398 points fetched.
Partition                     PowerVC-LE:     7163 points fetched.
Partition                   LVL-cluster2:     7163 points fetched.
Partition                   lvl-cluster1:     7163 points fetched.
Partition                       WM-SLES2:    18031 points fetched.
{{< /highlight >}}
