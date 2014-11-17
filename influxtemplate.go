// nmon2influx
// import nmon report in Influxdb
// author: adejoux@djouctech.net

package main

import (
    "fmt"
    "bytes"
    "text/template"
)
const influxtempl = `
{
  "id": null,
  "title": "nmon reports",
  "originalTitle": "nmon reports",
  "tags": [],
  "style": "dark",
  "timezone": "browser",
  "editable": true,
  "hideControls": false,
  "rows": [
    {
      "title": "Row1",
      "height": "250px",
      "editable": true,
      "collapse": false,
      "panels": [
        {
          "error": false,
          "span": 6,
          "editable": true,
          "type": "graph",
          "id": 1,
          "datasource": null,
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
              "column": "User%",
              "series": "CPU_ALL",
              "query": "select mean(\"User%\") from \"CPU_ALL\" where $timeFilter group by time($interval) order asc",
              "alias": "User%",
              "rawQuery": true,
              "hide": false
            },
            {
              "function": "mean",
              "column": "User%",
              "series": "CPU_ALL",
              "query": "select mean(\"Sys%\") from \"CPU_ALL\" where $timeFilter group by time($interval) order asc",
              "alias": "Sys%",
              "rawQuery": true,
              "hide": false
            },
            {
              "function": "mean",
              "column": "User%",
              "series": "CPU_ALL",
              "query": "select mean(\"Idle%\") from \"CPU_ALL\" where $timeFilter group by time($interval) order asc",
              "alias": "Idle%",
              "rawQuery": true,
              "hide": false
            },
            {
              "function": "mean",
              "column": "User%",
              "series": "CPU_ALL",
              "query": "select mean(\"Wait%\") from \"CPU_ALL\" where $timeFilter group by time($interval) order asc",
              "alias": "Wait%",
              "rawQuery": true,
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
          "span": 6,
          "editable": true,
          "type": "graph",
          "id": 2,
          "datasource": null,
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
              "column": "EC_User%",
              "series": "LPAR",
              "query": "select mean(\"EC_User%\") from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "EC_User%",
              "rawQuery": true,
              "hide": false
            },
            {
              "function": "mean",
              "column": "usedPoolCPU%",
              "series": "LPAR",
              "query": "select mean(\"usedPoolCPU%\") from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "usedPoolCPU%",
              "rawQuery": true,
              "hide": false
            },
            {
              "function": "mean",
              "column": "EC_Idle%",
              "series": "LPAR",
              "query": "select mean(\"EC_Idle%\") from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "EC_Idle%",
              "rawQuery": true,
              "hide": false
            },
            {
              "function": "mean",
              "column": "VP_Sys%",
              "series": "LPAR",
              "query": "select mean(\"VP_Sys%\") from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "VP_Sys%",
              "rawQuery": true,
              "hide": false
            },
            {
              "function": "mean",
              "column": "VP_User%",
              "series": "LPAR",
              "query": "select mean(\"VP_User%\") from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "VP_User%",
              "rawQuery": true,
              "hide": false
            },
            {
              "function": "mean",
              "column": "VP_User%",
              "series": "LPAR",
              "query": "select mean(\"VP_Wait%\") from \"LPAR\" where $timeFilter group by time($interval) order asc",
              "alias": "VP_Wait%",
              "rawQuery": true,
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
      "title": "New row",
      "height": "250px",
      "editable": true,
      "collapse": false,
      "panels": [
        {
          "error": false,
          "span": 12,
          "editable": true,
          "type": "graph",
          "id": 3,
          "datasource": null,
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
          "lines": false,
          "fill": 0,
          "linewidth": 1,
          "points": false,
          "pointradius": 5,
          "bars": true,
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
          {{ range $adapter := .GetIOADAPT }}
            {
              "function": "mean",
              "column": "{{$adapter}}",
              "series": "IOADAPT",
              "query": "select mean(\"{{$adapter}}\") from \"IOADAPT\" where $timeFilter group by time($interval) order asc",
              "rawQuery": true,
              "alias": "{{$adapter}}"
            },

          {{end}}
          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "title": "IOADAPT",
          "leftYAxisLabel": "KB/s"
        }
      ]
    },
    {
      "title": "New row",
      "height": "250px",
      "editable": true,
      "collapse": false,
      "panels": [
        {
          "error": false,
          "span": 6,
          "editable": true,
          "type": "graph",
          "id": 4,
          "datasource": null,
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
              "column": "%minperm",
              "series": "MEMUSE",
              "query": "select mean(\"%minperm\") from \"MEMUSE\" where $timeFilter group by time($interval) order asc",
              "rawQuery": true,
              "alias": "%minperm"
            },
            {
              "function": "mean",
              "column": "%minperm",
              "series": "MEMUSE",
              "query": "select mean(\"%maxperm\") from \"MEMUSE\" where $timeFilter group by time($interval) order asc",
              "rawQuery": true,
              "alias": "%maxperm"
            },
            {
              "function": "mean",
              "column": "%numperm",
              "series": "MEMUSE",
              "query": "select mean(\"%numperm\") from \"MEMUSE\" where $timeFilter group by time($interval) order asc",
              "rawQuery": true,
              "alias": "%numperm"
            },
            {
              "function": "mean",
              "column": "%numperm",
              "series": "MEMUSE",
              "query": "select mean(\"%numclient\") from \"MEMUSE\" where $timeFilter group by time($interval) order asc",
              "rawQuery": true,
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
          "span": 6,
          "editable": true,
          "type": "graph",
          "id": 5,
          "datasource": null,
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
              "column": "Real free(MB)",
              "series": "MEM",
              "query": "select mean(\"Real free(MB)\") from \"MEM\" where $timeFilter group by time($interval) order asc",
              "rawQuery": true,
              "alias": "Real free(MB)"
            },
            {
              "function": "mean",
              "column": "Virtual free(MB)",
              "series": "MEM",
              "query": "select mean(\"Virtual free(MB)\") from \"MEM\" where $timeFilter group by time($interval) order asc",
              "rawQuery": true,
              "alias": "Virtual free(MB)"
            },
            {
              "function": "mean",
              "column": "Real total(MB)",
              "series": "MEM",
              "query": "select mean(\"Real total(MB)\") from \"MEM\" where $timeFilter group by time($interval) order asc",
              "rawQuery": true,
              "alias": "Real total(MB)"
            },
            {
              "function": "mean",
              "column": "Virtual total(MB)",
              "series": "MEM",
              "query": "select mean(\"Virtual total(MB)\") from \"MEM\" where $timeFilter group by time($interval) order asc",
              "rawQuery": true,
              "alias": "Virtual total(MB)"
            }
          ],
          "aliasColors": {},
          "seriesOverrides": [],
          "title": "MEM",
          "leftYAxisLabel": "MB"
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
    "from": "2014-09-23T04:51:17.122Z",
    "to": "2014-09-23T23:09:50.129Z",
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
}%
`

func (influx *Influx) GetTemplate() {

    tmpl := template.New("influxtempl")
    tmpl.Parse(influxtempl)
    fmt.Printf("template:\n")
    result := new(bytes.Buffer)
    err := tmpl.Execute(result, influx)
    if err != nil {
        panic(err)
    }
    fmt.Println(result)

}

func (influx *Influx) GetIOADAPT() ([]string) {
   return influx.DataSeries["IOADAPT"].Columns
}
