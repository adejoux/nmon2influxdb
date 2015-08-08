// nmon2influxdb
// import nmon data in InfluxDB
// author: adejoux@djouxtech.net

package main

const influxtempl = `
{
  "id": null,
  "title": "{{.Hostname}} nmon report",
  "originalTitle": "{{.Hostname}} nmon report",
  "tags": [],
  "style": "light",
  "timezone": "browser",
  "editable": true,
  "hideControls": false,
  "rows": [
      {
      "title": "INFORMATION",
      "height": "250px",
      "editable": true,
      "collapse": true,
      "panels": [
        {
          "error": false,
          "span": 12,
          "editable": true,
          "type": "text",
          "id": 10,
          "mode": "html",
          "content": "<table>{{.TextContent}}</table>",
          "style": {},
          "title": "INFORMATION"
        }
      ]
    },
    {
      "title": "CPU",
      "height": "250px",
      "editable": true,
      "collapse": false,
      "panels": [
        {
          "error": false,
          "span": 4,
          "editable": true,
          "type": "graph",
          "id": 1,
          "datasource": "{{ .DataSource }}",
          "renderer": "flot",
          "x-axis": true,
          "y-axis": true,
          "scale": 1,
          "y_formats": [
            "short",
            "short"
          ],
          "grid": {
            "leftMax": null,
            "rightMax": null,
            "leftMin": null,
            "rightMin": null,
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "annotate": {
            "enable": false
          },
          "resolution": 100,
          "lines": true,
          "fill": 1,
          "linewidth": 1,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": true,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false
          },
          "percentage": false,
          "zerofill": true,
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "individual",
            "query_as_alias": true
          },
          "targets": [
            {
              "function": "mean",
              "column": "\"User%\"",
              "series": "CPU_ALL",
              "query": "select mean(\"User%\") from \"CPU_ALL\" where $timeFilter group by time($interval) order asc",
              "alias": "User%",
              "rawQuery": false,
              "hide": false
            },
            {
              "function": "mean",
              "column": "\"Sys%\"",
              "series": "CPU_ALL",
              "query": "select mean(\"Sys%\") from \"CPU_ALL\" where $timeFilter group by time($interval) order asc",
              "alias": "Sys%",
              "rawQuery": false,
              "hide": false
            },
            {
              "function": "mean",
              "column": "\"Wait%\"",
              "series": "CPU_ALL",
              "query": "select mean(\"Wait%\") from \"CPU_ALL\" where $timeFilter group by time($interval) order asc",
              "alias": "Wait%",
              "rawQuery": false,
              "hide": false
            },
            {
              "function": "mean",
              "column": "\"Idle%\"",
              "series": "CPU_ALL",
              "query": "select mean(\"Idle%\") from \"CPU_ALL\" where $timeFilter group by time($interval) order asc",
              "alias": "Idle%",
              "rawQuery": false,
              "hide": false
            }
          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "title": "CPU_ALL",
          "leftYAxisLabel": "%"
        },
        {
          "error": false,
          "span": 4,
          "editable": true,
          "type": "graph",
          "id": 2,
          "datasource": "{{ .DataSource }}",
          "renderer": "flot",
          "x-axis": true,
          "y-axis": true,
          "scale": 1,
          "y_formats": [
            "short",
            "short"
          ],
          "grid": {
            "leftMax": null,
            "rightMax": null,
            "leftMin": null,
            "rightMin": null,
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "annotate": {
            "enable": false
          },
          "resolution": 100,
          "lines": true,
          "fill": 1,
          "linewidth": 1,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": true,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false
          },
          "percentage": false,
          "zerofill": true,
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "cumulative",
            "query_as_alias": true
          },
          "targets": [
            {
              "function": "mean",
              "column": "\"EC_User%\"",
              "series": "LPAR",
              "query": "select mean(\"EC_User%\") from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "EC_User%",
              "rawQuery": false,
              "hide": false
            },
            {
              "function": "mean",
              "column": "\"EC_Sys%\"",
              "series": "LPAR",
              "query": "select mean(\"EC_Sys%\") from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "EC_Sys%",
              "rawQuery": false,
              "hide": false
            },
            {
              "function": "mean",
              "column": "\"EC_Wait%\"",
              "series": "LPAR",
              "query": "select mean(\"EC_Wait%\") from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "EC_Wait%",
              "rawQuery": false,
              "hide": false
            },
            {
              "function": "mean",
              "column": "\"EC_Idle%\"",
              "series": "LPAR",
              "query": "select mean(\"EC_Idle%\") from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "EC_Idle%",
              "rawQuery": false,
              "hide": false
            }
          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "title": "LPAR",
          "leftYAxisLabel": "%"
        },
        {
          "error": false,
          "span": 4,
          "editable": true,
          "type": "graph",
          "id": 2,
          "datasource": "{{ .DataSource }}",
          "renderer": "flot",
          "x-axis": true,
          "y-axis": true,
          "scale": 1,
          "y_formats": [
            "short",
            "short"
          ],
          "grid": {
            "leftMax": null,
            "rightMax": null,
            "leftMin": null,
            "rightMin": null,
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "annotate": {
            "enable": false
          },
          "resolution": 100,
          "lines": true,
          "fill": 1,
          "linewidth": 1,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": true,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false
          },
          "percentage": false,
          "zerofill": true,
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "cumulative",
            "query_as_alias": true
          },
          "targets": [
            {
              "function": "mean",
              "column": "\"VP_User%\"",
              "series": "LPAR",
              "query": "select mean(\"VP_User%\") from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "VP_User%",
              "rawQuery": false,
              "hide": false
            },
            {
              "function": "mean",
              "column": "\"VP_Sys%\"",
              "series": "LPAR",
              "query": "select mean(\"VP_Sys%\") from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "VP_Sys%",
              "rawQuery": false,
              "hide": false
            },
            {
              "function": "mean",
              "column": "\"VP_Wait%\"",
              "series": "LPAR",
              "query": "select mean(\"VP_Wait%\") from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "VP_Wait%",
              "rawQuery": false,
              "hide": false
            },
            {
              "function": "mean",
              "column": "\"VP_Idle%\"",
              "series": "LPAR",
              "query": "select mean(\"VP_Wait%\") from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "VP_Wait%",
              "rawQuery": false,
              "hide": false
            }
          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "title": "LPAR",
          "leftYAxisLabel": "%"
        }
      ]
    },
    {
      "title": "LPAR",
      "height": "250px",
      "editable": true,
      "collapse": false,
      "panels": [
        {
          "error": false,
          "span": 12,
          "editable": true,
          "type": "graph",
          "id": 6,
          "datasource": "{{ .DataSource }}",
          "renderer": "flot",
          "x-axis": true,
          "y-axis": true,
          "scale": 1,
          "y_formats": [
            "short",
            "short"
          ],
          "grid": {
            "leftMax": null,
            "rightMax": null,
            "leftMin": null,
            "rightMin": null,
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "annotate": {
            "enable": false
          },
          "resolution": 100,
          "lines": true,
          "fill": 0,
          "linewidth": 1,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": false,
          "legend": {
            "show": true,
            "values": true,
            "min": true,
            "max": true,
            "current": false,
            "total": false,
            "avg": true
          },
          "percentage": false,
          "zerofill": true,
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "cumulative",
            "query_as_alias": true
          },
          "targets": [
            {
              "function": "mean",
              "column": "PhysicalCPU",
              "series": "LPAR",
              "query": "select mean(PhysicalCPU) from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "PhysicalCPU"
            },
            {
              "function": "mean",
              "column": "entitled",
              "series": "LPAR",
              "query": "select mean(entitled) from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "entitled"
            },
            {
              "function": "mean",
              "column": "Folded",
              "series": "LPAR",
              "query": "select mean(Folded) from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "folded"
            },
            {
              "function": "mean",
              "column": "virtualCPUs",
              "series": "LPAR",
              "query": "select mean(virtualCPUs) from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "virtualCPUs"
            },
            {
              "function": "mean",
              "column": "logicalCPUs",
              "series": "LPAR",
              "query": "select mean(logicalCPUs) from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "logicalCPUs"
            },
            {
              "function": "mean",
              "column": "PoolIdle",
              "series": "LPAR",
              "query": "select mean(PoolIdle) from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "PoolIdle"
            },
            {
              "function": "mean",
              "column": "poolCPUs",
              "series": "LPAR",
              "query": "select mean(poolCPUs) from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "poolCPUs"
            }
          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "title": "LPAR",
          "leftYAxisLabel": "cpu"
        }
      ]
    },
    {
      "title": "MEMORY",
      "height": "250px",
      "editable": true,
      "collapse": true,
      "panels": [
        {
          "error": false,
          "span": 4,
          "editable": true,
          "type": "graph",
          "id": 4,
          "datasource": "{{ .DataSource }}",
          "renderer": "flot",
          "x-axis": true,
          "y-axis": true,
          "scale": 1,
          "y_formats": [
            "short",
            "short"
          ],
          "grid": {
            "leftMax": null,
            "rightMax": null,
            "leftMin": null,
            "rightMin": null,
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "annotate": {
            "enable": false
          },
          "resolution": 100,
          "lines": true,
          "fill": 0,
          "linewidth": 1,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": false,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false
          },
          "percentage": false,
          "zerofill": true,
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "cumulative",
            "query_as_alias": true
          },
          "targets": [
            {
              "function": "mean",
              "column": "\"%minperm\"",
              "series": "MEMUSE",
              "query": "select mean(\"%minperm\") from \"MEMUSE\" where $timeFilter group by time($interval) order asc",
              "rawQuery": false,
              "alias": "%minperm"
            },
            {
              "function": "mean",
              "column": "\"%maxperm\"",
              "series": "MEMUSE",
              "query": "select mean(\"%maxperm\") from \"MEMUSE\" where $timeFilter group by time($interval) order asc",
              "rawQuery": false,
              "alias": "%maxperm"
            },
            {
              "function": "mean",
              "column": "\"%numperm\"",
              "series": "MEMUSE",
              "query": "select mean(\"%numperm\") from \"MEMUSE\" where $timeFilter group by time($interval) order asc",
              "rawQuery": false,
              "alias": "%numperm"
            },
            {
              "function": "mean",
              "column": "\"%numclient\"",
              "series": "MEMUSE",
              "query": "select mean(\"%numclient\") from \"MEMUSE\" where $timeFilter group by time($interval) order asc",
              "rawQuery": false,
              "alias": "%numclient"
            }
          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "title": "MEMUSE",
          "leftYAxisLabel": "%"
        },
        {
          "error": false,
          "span": 4,
          "editable": true,
          "type": "graph",
          "id": 5,
          "datasource": "{{ .DataSource }}",
          "renderer": "flot",
          "x-axis": true,
          "y-axis": true,
          "scale": 1,
          "y_formats": [
            "short",
            "short"
          ],
          "grid": {
            "leftMax": null,
            "rightMax": null,
            "leftMin": null,
            "rightMin": null,
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "annotate": {
            "enable": false
          },
          "resolution": 100,
          "lines": true,
          "fill": 1,
          "linewidth": 1,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": false,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false
          },
          "percentage": false,
          "zerofill": true,
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "cumulative",
            "query_as_alias": true
          },
          "targets": [
            {
              "function": "mean",
              "column": "\"Real free(MB)\"",
              "series": "MEM",
              "query": "select mean(\"Real free(MB)\") from \"MEM\" where $timeFilter group by time($interval) order asc",
              "rawQuery": false,
              "alias": "Real free(MB)"
            },
            {
              "function": "mean",
              "column": "\"Virtual free(MB)\"",
              "series": "MEM",
              "query": "select mean(\"Virtual free(MB)\") from \"MEM\" where $timeFilter group by time($interval) order asc",
              "rawQuery": false,
              "alias": "Virtual free(MB)"
            },
            {
              "function": "mean",
              "column": "\"Real total(MB)\"",
              "series": "MEM",
              "query": "select mean(\"Real total(MB)\") from \"MEM\" where $timeFilter group by time($interval) order asc",
              "rawQuery": false,
              "alias": "Real total(MB)"
            },
            {
              "function": "mean",
              "column": "\"Virtual total(MB)\"",
              "series": "MEM",
              "query": "select mean(\"Virtual total(MB)\") from \"MEM\" where $timeFilter group by time($interval) order asc",
              "rawQuery": false,
              "alias": "Virtual total(MB)"
            }
          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "title": "MEM",
          "leftYAxisLabel": "MB"
        },
        {
          "error": false,
          "span": 4,
          "editable": true,
          "type": "graph",
          "id": 5,
          "datasource": "{{ .DataSource }}",
          "renderer": "flot",
          "x-axis": true,
          "y-axis": true,
          "scale": 1,
          "y_formats": [
            "short",
            "short"
          ],
          "grid": {
            "leftMax": null,
            "rightMax": null,
            "leftMin": null,
            "rightMin": null,
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "annotate": {
            "enable": false
          },
          "resolution": 100,
          "lines": true,
          "fill": 1,
          "linewidth": 1,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": true,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false
          },
          "percentage": false,
          "zerofill": true,
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "cumulative",
            "query_as_alias": true
          },
          "targets": [
            {
              "function": "mean",
              "column": "\"Process%\"",
              "series": "MEMNEW",
              "query": "select mean(\"Process%\") from \"MEMNEW\" where $timeFilter group by time($interval) order asc",
              "rawQuery": false,
              "alias": "Process%"
            },
            {
              "function": "mean",
              "column": "\"FScache%\"",
              "series": "MEMNEW",
              "query": "select mean(\"FScache%\") from \"MEMNEW\" where $timeFilter group by time($interval) order asc",
              "rawQuery": false,
              "alias": "FScache%"
            },
            {
              "function": "mean",
              "column": "\"System%\"",
              "series": "MEMNEW",
              "query": "select mean(\"System%\") from \"MEMNEW\" where $timeFilter group by time($interval) order asc",
              "rawQuery": false,
              "alias": "System%"
            }

          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "title": "MEMNEW",
          "leftYAxisLabel": "%"
        }
      ]
    },
    {
      "title": "IOADAPT",
      "height": "250px",
      "editable": true,
      "collapse": true,
      "panels": [
        {
          "error": false,
          "span": 12,
          "editable": true,
          "type": "graph",
          "id": 3,
          "datasource": "{{ .DataSource }}",
          "renderer": "flot",
          "x-axis": true,
          "y-axis": true,
          "scale": 1,
          "y_formats": [
            "short",
            "short"
          ],
          "grid": {
            "leftMax": null,
            "rightMax": null,
            "leftMin": null,
            "rightMin": null,
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "annotate": {
            "enable": false
          },
          "resolution": 100,
          "lines": true,
          "fill": 4,
          "linewidth": 1,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": true,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false
          },
          "percentage": false,
          "zerofill": true,
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "individual",
            "query_as_alias": true
          },
          "targets": [
            {{ range $index, $adapter := .GetFilteredColumns "IOADAPT"  "-KB/s"}}{{ if $index}},{{end}}
                {
                  "function": "mean",
                  "column": "\"{{$adapter}}\"",
                  "series": "IOADAPT",
                  "query": "select mean(\"{{$adapter}}\") from \"IOADAPT\" where $timeFilter group by time($interval) order asc",
                  "rawQuery": false,
                  "alias": "{{$adapter}}"
                }{{end}}
          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "title": "IOADAPT",
          "leftYAxisLabel": "KB/s"
        }
      ]
    },
    {{ if .GetColumns "SEA"}}
      {
        "title": "SEA",
        "height": "250px",
        "editable": true,
        "collapse": true,
        "panels": [
          {
            "error": false,
            "span": 12,
            "editable": true,
            "type": "graph",
            "id": 3,
            "renderer": "flot",
            "x-axis": true,
            "y-axis": true,
            "scale": 1,
            "y_formats": [
              "short",
              "short"
            ],
            "grid": {
              "leftMax": null,
              "rightMax": null,
              "leftMin": null,
              "rightMin": null,
              "threshold1": null,
              "threshold2": null,
              "threshold1Color": "rgba(216, 200, 27, 0.27)",
              "threshold2Color": "rgba(234, 112, 112, 0.22)"
            },
            "annotate": {
              "enable": false
            },
            "resolution": 100,
            "lines": true,
            "fill": 0,
            "linewidth": 1,
            "points": false,
            "pointradius": 5,
            "bars": false,
            "stack": false,
            "legend": {
              "show": true,
              "values": false,
              "min": false,
              "max": false,
              "current": false,
              "total": false,
              "avg": false
            },
            "percentage": false,
            "zerofill": true,
            "nullPointMode": "connected",
            "steppedLine": false,
            "tooltip": {
              "value_type": "individual",
              "query_as_alias": true
            },
            "targets": [
              {{ range $index, $adapter := .GetColumns "SEA" }}{{ if $index}},{{end}}
                  {
                    "function": "mean",
                    "column": "\"{{$adapter}}\"",
                    "series": "SEA",
                    "query": "select mean(\"{{$adapter}}\") from \"SEA\" where $timeFilter group by time($interval) order asc",
                    "rawQuery": false,
                    "alias": "{{$adapter}}"
                  }{{end}}
            ],
            "aliasColors": {},
            "seriesOverrides": [],
            "title": "SEA",
            "leftYAxisLabel": "KB/s"
          }
        ]
      },
    {{end}}
    {{ if .GetColumns "NPIV"}}
      {
        "title": "NPIV",
        "height": "250px",
        "editable": true,
        "collapse": true,
        "panels": [
          {
            "error": false,
            "span": 12,
            "editable": true,
            "type": "graph",
            "id": 3,
            "datasource": "{{ .DataSource }}",
            "renderer": "flot",
            "x-axis": true,
            "y-axis": true,
            "scale": 1,
            "y_formats": [
              "short",
              "short"
            ],
            "grid": {
              "leftMax": null,
              "rightMax": null,
              "leftMin": null,
              "rightMin": null,
              "threshold1": null,
              "threshold2": null,
              "threshold1Color": "rgba(216, 200, 27, 0.27)",
              "threshold2Color": "rgba(234, 112, 112, 0.22)"
            },
            "annotate": {
              "enable": false
            },
            "resolution": 100,
            "lines": true,
            "fill": 0,
            "linewidth": 1,
            "points": false,
            "pointradius": 5,
            "bars": false,
            "stack": false,
            "legend": {
              "show": true,
              "values": false,
              "min": false,
              "max": false,
              "current": false,
              "total": false,
              "avg": false
            },
            "percentage": false,
            "zerofill": true,
            "nullPointMode": "connected",
            "steppedLine": false,
            "tooltip": {
              "value_type": "individual",
              "query_as_alias": true
            },
            "targets": [
              {{ range $index, $adapter := .GetFilteredColumns "NPIV" "e-KB"}}{{ if $index}},{{end}}
                  {
                    "function": "mean",
                    "column": "\"{{$adapter}}\"",
                    "series": "NPIV",
                    "query": "select mean(\"{{$adapter}}\") from \"NPIV\" where $timeFilter group by time($interval) order asc",
                    "rawQuery": false,
                    "alias": "{{$adapter}}"
                  }{{end}}
            ],
            "aliasColors": {},
            "seriesOverrides": [],
            "title": "NPIV",
            "leftYAxisLabel": "KB/s"
          }
        ]
      },
    {{end}}
        {{ if .GetColumns "SEACLITRAFFIC"}}
      {
        "title": "SEACLITRAFFIC",
        "height": "250px",
        "editable": true,
        "collapse": true,
        "panels": [
          {
            "error": false,
            "span": 12,
            "editable": true,
            "type": "graph",
            "id": 3,
            "datasource": "{{ .DataSource }}",
            "renderer": "flot",
            "x-axis": true,
            "y-axis": true,
            "scale": 1,
            "y_formats": [
              "short",
              "short"
            ],
            "grid": {
              "leftMax": null,
              "rightMax": null,
              "leftMin": null,
              "rightMin": null,
              "threshold1": null,
              "threshold2": null,
              "threshold1Color": "rgba(216, 200, 27, 0.27)",
              "threshold2Color": "rgba(234, 112, 112, 0.22)"
            },
            "annotate": {
              "enable": false
            },
            "resolution": 100,
            "lines": true,
            "fill": 3,
            "linewidth": 1,
            "points": false,
            "pointradius": 5,
            "bars": false,
            "stack": true,
            "legend": {
              "show": true,
              "values": false,
              "min": false,
              "max": false,
              "current": false,
              "total": false,
              "avg": false
            },
            "percentage": false,
            "zerofill": true,
            "nullPointMode": "connected",
            "steppedLine": false,
            "tooltip": {
              "value_type": "individual",
              "query_as_alias": true
            },
            "targets": [
              {{ range $index, $adapter := .GetFilteredColumns "SEACLITRAFFIC" "e-KB"}}{{ if $index}},{{end}}
                  {
                    "function": "mean",
                    "column": "\"{{$adapter}}\"",
                    "series": "SEACLITRAFFIC",
                    "query": "select mean(\"{{$adapter}}\") from \"SEACLITRAFFIC\" where $timeFilter group by time($interval) order asc",
                    "rawQuery": false,
                    "alias": "{{$adapter}}"
                  }{{end}}
            ],
            "aliasColors": {},
            "seriesOverrides": [],
            "title": "SEACLITRAFFIC",
            "leftYAxisLabel": "KB/s"
          }
        ]
      },
    {{end}}
    {
      "title": "NET",
      "height": "250px",
      "editable": true,
      "collapse": true,
      "panels": [
        {
          "error": false,
          "span": 12,
          "editable": true,
          "type": "graph",
          "id": 7,
          "datasource": "{{ .DataSource }}",
          "renderer": "flot",
          "x-axis": true,
          "y-axis": true,
          "scale": 1,
          "y_formats": [
            "short",
            "short"
          ],
          "grid": {
            "leftMax": null,
            "rightMax": null,
            "leftMin": null,
            "rightMin": null,
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "annotate": {
            "enable": false
          },
          "resolution": 100,
          "lines": true,
          "fill": 1,
          "linewidth": 1,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": true,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false
          },
          "percentage": false,
          "zerofill": true,
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "cumulative",
            "query_as_alias": true
          },
          "targets": [
            {{ range $index, $adapter := .GetColumns "NET" }}{{ if $index}},{{end}}
                {
                    "function": "mean",
                    "column": "\"{{$adapter}}\"",
                    "series": "NET",
                    "query": "select mean(\"{{$adapter}}\") from \"NET\" where $timeFilter group by time($interval) order asc",
                    "rawQuery": false,
                    "alias": "{{$adapter}}"
                }{{end}}
            ],
          "aliasColors": {},
          "seriesOverrides": [],
          "title": "NET"
        }
      ]
    },
    {
      "title": "PAGE",
      "height": "250px",
      "editable": true,
      "collapse": true,
      "panels": [
        {
          "error": false,
          "span": 12,
          "editable": true,
          "type": "graph",
          "id": 8,
          "datasource": "{{ .DataSource }}",
          "renderer": "flot",
          "x-axis": true,
          "y-axis": true,
          "scale": 1,
          "y_formats": [
            "short",
            "short"
          ],
          "grid": {
            "leftMax": null,
            "rightMax": null,
            "leftMin": null,
            "rightMin": null,
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "annotate": {
            "enable": false
          },
          "resolution": 100,
          "lines": true,
          "fill": 0,
          "linewidth": 1,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": false,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false
          },
          "percentage": false,
          "zerofill": true,
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "cumulative",
            "query_as_alias": true
          },
          "targets": [
            {
              "function": "mean",
              "column": "pgin",
              "series": "PAGE",
              "query": "select mean(pgin) from \"PAGE\" where $timeFilter group by time($interval) order asc",
              "alias": "pgin"
            },
            {
              "function": "mean",
              "column": "pgout",
              "series": "PAGE",
              "query": "select mean(pgout) from \"PAGE\" where $timeFilter group by time($interval) order asc",
              "alias": "pgout"
            },
            {
              "function": "mean",
              "column": "scans",
              "series": "PAGE",
              "query": "select mean(scans) from \"PAGE\" where $timeFilter group by time($interval) order asc",
              "alias": "scans"
            },
            {
              "function": "mean",
              "column": "cycles",
              "series": "PAGE",
              "query": "select mean(cycles) from \"PAGE\" where $timeFilter group by time($interval) order asc",
              "alias": "cycles"
            },
            {
              "function": "mean",
              "column": "faults",
              "series": "PAGE",
              "query": "select mean(faults) from \"PAGE\" where $timeFilter group by time($interval) order asc",
              "alias": "faults"
            },
            {
              "function": "mean",
              "column": "pgsin",
              "series": "PAGE",
              "query": "select mean(pgsin) from \"PAGE\" where $timeFilter group by time($interval) order asc",
              "alias": "pgsin"
            },
            {
              "function": "mean",
              "column": "pgsout",
              "series": "PAGE",
              "query": "select mean(pgsout) from \"PAGE\" where $timeFilter group by time($interval) order asc",
              "alias": "pgsout"
            },
            {
              "function": "mean",
              "column": "reclaims",
              "series": "PAGE",
              "query": "select mean(reclaims) from \"PAGE\" where $timeFilter group by time($interval) order asc",
              "alias": "reclaims"
            }
          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "title": "PAGE"
        }
      ]
    }
  ],
  "nav": [
    {
      "type": "timepicker",
      "enable": true,
      "status": "Stable",
      "time_options": [
        "5m",
        "15m",
        "1h",
        "6h",
        "12h",
        "24h",
        "2d",
        "7d",
        "30d"
      ],
      "refresh_intervals": [
        "5s",
        "10s",
        "30s",
        "1m",
        "5m",
        "15m",
        "30m",
        "1h",
        "2h",
        "1d"
      ],
      "now": false,
      "collapse": false,
      "notice": false
    }
  ],
  "time": {
    "from": "{{.StartTime}}",
    "to": "{{.StopTime}}",
    "now": false
  },
  "templating": {
    "list": []
  },
  "annotations": {
    "list": [],
    "enable": false
  },
  "refresh": false,
  "version": 6
}
`
