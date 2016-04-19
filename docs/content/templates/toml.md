---
date: 2016-04-18T18:05:46+02:00
title: toml
menu:
  main:
    parent: Templates
    identifier: /templates/toml
    weight: 10
---

TOML templates are available to customize quickly grafana dashboard for deployment in new instances. It's easier to read than JSON and allow quick modifications.

# simple template

{{< highlight toml >}}
title = "dual vio"

[[row]]
title = "LPAR"
    [[row.panel]]
    title = "shared processor"
        [[row.panel.metric]]
            measurement = "LPAR"
            hosts = ["vios1", "vios2"]
            fields = ["PhysicalCPU", "virtualCPUs", "entitled"]
        [[row.panel.override]]
            alias = "$tag_host PhysicalCPU"
            stack = true
            fill = 1
        [row.panel.tooltip]
            value_type = "individual"

    [[row.panel]]
        title = "shared processor pool"
        [[row.panel.metric]]
            measurement = "LPAR"
            hosts = ["vios1"]
            fields = ["PoolIdle", "poolCPUs"]

[[row]]
title = "SEA"
    [[row.panel]]
    title = "SEA WRITE"
    stack = true
    fill = 1
        [[row.panel.metric]]
            measurement = "SEA"
            # regular expression can be used for host. Here it will be "vios".
            hosts = ["vios1"]
            # regular expressions are used on fields
            fields = ["write-KB"]
    [[row.panel]]
    title = "SEA READ"
    stack = true
    fill = 1
        [[row.panel.metric]]
            measurement = "SEA"
            hosts = ["vios1", "vios2"]
            fields = ["read-KB"]
[time]
  from = "2015-07-06T13:03:56.000Z"
  to = "2015-07-07T14:23:44.000Z"
{{< /highlight >}}

Upload the template:
{{< highlight toml >}}
nmon2influxdb dashboard simple_dual_vios.toml
Dashboard uploaded to grafana
{{< /highlight >}}

# dual vio template

This template is more advanced. It's using templating to create variables to select dynamically vio servers.

{{< highlight toml >}}

title = "templated dual vio"
[templates]
  [[templates.template]]

    # name of the template variable
    name = "vios"
    # the default nmon2influxdb database in InfluxDB
    datasource = "nmon_reports"
    # query used to retrieve values in InfluxDB
    query =  "SHOW TAG VALUES FROM CPU_ALL WITH KEY = host"
    # Regular expression used to filter the query result
    regex = "vios"
[[row]]
title = "LPAR"
    [[row.panel]]
    title = "shared processor"
        [[row.panel.metric]]
            measurement = "LPAR"
            # use template variable $vios
            hosts = ["$vios"]
            fields = ["PhysicalCPU", "virtualCPUs", "entitled"]
        [[row.panel.override]]
            alias = "$tag_host PhysicalCPU"
            stack = true
            fill = 1
        [row.panel.tooltip]
            value_type = "individual"

    [[row.panel]]
        title = "shared processor pool"
        [[row.panel.metric]]
            measurement = "LPAR"
            # use template variable $vios
            hosts = ["$vios"]
            fields = ["PoolIdle", "poolCPUs"]

[[row]]
title = "SEA"
    [[row.panel]]
    title = "SEA WRITE"
    stack = true
    fill = 1
        [[row.panel.metric]]
            measurement = "SEA"
            # use template variable $vios
            hosts = ["$vios"]
            # regular expressions are used on fields
            fields = ["write-KB"]
    [[row.panel]]
    title = "SEA READ"
    stack = true
    fill = 1
        [[row.panel.metric]]
            measurement = "SEA"
            # use template variable $vios
            hosts = ["$vios"]
            # regular expressions are used on fields
            fields = ["read-KB"]
[time]
  from = "2015-07-06T13:03:56.000Z"
  to = "2015-07-07T14:23:44.000Z"
{{< /highlight >}}

Upload the dashboard:
{{< highlight toml >}}
nmon2influxdb dashboard templated_dual_vios.toml
Dashboard uploaded to grafana
{{< /highlight >}}

# Other templates

Other templates are available in the [github repository](https://github.com/adejoux/nmon2influxdb/tree/master/templates).
