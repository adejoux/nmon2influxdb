---
date: 2016-04-15T16:38:51+02:00
title: stats
menu:
  main:
    parent: Usage
    identifier: /usage/stats
    weight: 30
---

{{< highlight batch >}}
NAME:
   nmon2influxdb stats - generate stats from a InfluxDB metric

USAGE:
   nmon2influxdb stats [command options] [arguments...]

OPTIONS:
   --metric, -m 	mandatory metric for stats
   --statshost, -s 	host metrics
   --from, -f 		from date
   --to, -t 		to date
   --sort "mean"	field to perform sort
   --limit, -l "20"	limit the output
   --filter 		specify a filter on fields
{{< /highlight >}}

# Parameters

* **metric**: specify the metric from where the stats will be generated.
* **statshost**: stats are generated for this host
* **from**: starting time in NMON format: "HH:mm:SS,DD-MMM-YYYY"
* **to**: end time in NMON format: "HH:mm:SS,DD-MMM-YYYY"
* **sort**: By default the results are sorted based on the **mean** values. Can also be **min**,**max**,**median**
* **limit**: limit the number of lines in result.
* **filter**: filter measurements.

# Examples

Generating stats for DISKREADSERV metric :
{{< highlight batch >}}
nmon2influxdb stats -m DISKREADSERV -s lpar1
          field|     Min|    Mean|  Median|     Max|Points #
        hdisk10|    0.40|    2.42|    2.10|   12.90|    1200
        hdisk11|    0.60|    2.63|    2.20|   14.00|    1200
        hdisk12|    0.50|    2.74|    2.30|   16.30|    1200
        hdisk13|    0.00|    0.52|    0.00|   16.30|    1200
         hdisk7|    0.00|    0.01|    0.00|    0.80|    1200
         hdisk8|    0.00|    0.04|    0.00|    1.00|    1200
         hdisk9|    0.50|    2.33|    2.10|   11.30|    1200
{{< /highlight >}}

Generating stats for **DISKREADSERV** metric on host **lpar1** on a time period :
{{< highlight batch >}}
nmon2influxdb stats -m DISKREADSERV -s lpar1 -f "08:04:05,29-May-2015" -t "15:04:05,29-May-2015"
          field|     Min|    Mean|  Median|     Max|Points #
        hdisk10|    2.20|    3.80|    3.65|    5.60|     132
        hdisk11|    2.40|    4.10|    3.90|    5.80|     132
        hdisk12|    2.30|    4.39|    4.05|    6.50|     132
        hdisk13|    0.00|    0.73|    0.00|   13.90|     132
         hdisk7|    0.00|    0.02|    0.00|    0.20|     132
         hdisk8|    0.00|    0.10|    0.00|    1.00|     132
         hdisk9|    2.10|    3.68|    3.55|    5.20|     132
{{< /highlight >}}

Generating stats for **DISKREADSERV** metric on host **lpar1** and limit output to the 5 most active disk :

{{< highlight batch >}}
nmon2influxdb stats -m DISKREADSERV -s lpar1 -l 5
          field|     Min|    Mean|  Median|     Max|Points #
        hdisk10|    2.20|    3.80|    3.65|    5.60|     132
        hdisk11|    2.40|    4.10|    3.90|    5.80|     132
        hdisk12|    2.30|    4.39|    4.05|    6.50|     132
        hdisk13|    0.00|    0.73|    0.00|   13.90|     132
         hdisk7|    0.00|    0.02|    0.00|    0.20|     132
{{< /highlight >}}
