---
date: 2016-04-19T13:59:05+02:00
title: json
menu:
  main:
    parent: Templates
    identifier: /templates/json
    weight: 20
---

# dashboard export

Grafana allows the export of any dashboard
{{< gallery image="json_export.png" >}}

It will be a json file containing all the parameters of your dashboard.

# dashboard upload

**nmon2influxdb** can upload this file directly in your grafana instance:
{{< highlight batch >}}
nmon2influxdb dashboard myfile.json
Dashboard uploaded to grafana
{{< /highlight >}}
