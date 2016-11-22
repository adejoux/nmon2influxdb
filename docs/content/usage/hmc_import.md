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
{{< /highlight >}}

# Parameters

  * **hmc**: HMC to use to fetch PCM data
  * **hmcuser**: HMC user to use for connection
  * **hmcpass**: add per cpu metrics
  * **managed_system**: fetch HMC PCM data only for this managed system

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
{{< /highlight >}}

It's possible to set all CLI parameters. It's also possible to change the InfluxDB database name with **hmc_database** and change the data retention with **hmc_data_retention**.

# Examples

Loading HMC metrics from HMC **myhmc**:

{{< highlight batch >}}
nmon2influxdb hmc import
Getting list of managed systems
MANAGED SYSTEM: p750A
partition powerVC: 2940 points
MANAGED SYSTEM: p720-NIM_RETIRED
Error getting PCM data
{{< /highlight >}}

Note: parameters can also be set in the configuration file **~/.nmon2influxdb.cfg**.

Loading HMC metrics from HMC **myhmc** for system **mysystem** only:

{{< highlight batch >}}
nmon2influxdb hmc import -m p750A
Getting list of managed systems
MANAGED SYSTEM: p750A
partition powerVC: 2940 points
Skipping system: p720-NIM_RETIRED
{{< /highlight >}}
