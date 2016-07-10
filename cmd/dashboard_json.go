package cmd

const dashboardJSON = `
{
  "dashboard": {
      "id": null,
      "title": "Belt",
      "originalTitle": "Belt",
      "tags": [],
      "style": "dark",
      "timezone": "browser",
      "editable": false,
      "hideControls": true,
      "sharedCrosshair": false,
      "rows": [
        {
          "collapse": false,
          "editable": true,
          "height": "250px",
          "panels": [
            {
              "cacheTimeout": null,
              "colorBackground": false,
              "colorValue": false,
              "colors": [
                "rgba(245, 54, 54, 0.9)",
                "rgba(237, 129, 40, 0.89)",
                "rgba(50, 172, 45, 0.97)"
              ],
              "datasource": null,
              "editable": true,
              "error": false,
              "format": "percent",
              "gauge": {
                "maxValue": 100,
                "minValue": 0,
                "show": true,
                "thresholdLabels": false,
                "thresholdMarkers": true
              },
              "height": "",
              "id": 2,
              "interval": "1s",
              "isNew": true,
              "links": [],
              "maxDataPoints": 100,
              "nullPointMode": "connected",
              "nullText": null,
              "postfix": "",
              "postfixFontSize": "50%",
              "prefix": "",
              "prefixFontSize": "50%",
              "span": 3,
              "sparkline": {
                "fillColor": "rgba(31, 118, 189, 0.18)",
                "full": false,
                "lineColor": "rgb(31, 120, 193)",
                "show": false
              },
              "targets": [
                {
                  "dsType": "influxdb",
                  "groupBy": [
                    {
                      "params": [
                        "$interval"
                      ],
                      "type": "time"
                    },
                    {
                      "params": [
                        "null"
                      ],
                      "type": "fill"
                    }
                  ],
                  "measurement": "cpu",
                  "policy": "default",
                  "query": "SELECT mean(\"usage_user\") + mean(\"usage_system\") FROM \"cpu\" WHERE $timeFilter GROUP BY time($interval) fill(null)",
                  "rawQuery": true,
                  "refId": "A",
                  "resultFormat": "time_series",
                  "select": [
                    [
                      {
                        "params": [
                          "usage_user"
                        ],
                        "type": "field"
                      },
                      {
                        "params": [],
                        "type": "mean"
                      }
                    ]
                  ],
                  "tags": []
                }
              ],
              "thresholds": "",
              "title": "Cluster CPU Average",
              "type": "singlestat",
              "valueFontSize": "80%",
              "valueMaps": [
                {
                  "op": "=",
                  "text": "N/A",
                  "value": "null"
                }
              ],
              "valueName": "avg"
            },
            {
              "aliasColors": {},
              "bars": false,
              "datasource": null,
              "editable": true,
              "error": false,
              "fill": 1,
              "grid": {
                "threshold1": null,
                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                "threshold2": null,
                "threshold2Color": "rgba(234, 112, 112, 0.22)"
              },
              "id": 1,
              "isNew": true,
              "legend": {
                "alignAsTable": true,
                "avg": false,
                "current": false,
                "max": false,
                "min": false,
                "show": true,
                "total": false,
                "values": false
              },
              "lines": true,
              "linewidth": 2,
              "links": [],
              "nullPointMode": "connected",
              "percentage": false,
              "pointradius": 5,
              "points": false,
              "renderer": "flot",
              "seriesOverrides": [],
              "span": 9,
              "stack": false,
              "steppedLine": false,
              "targets": [
                {
                  "alias": "[[tag_host]]",
                  "dsType": "influxdb",
                  "groupBy": [
                    {
                      "params": [
                        "$interval"
                      ],
                      "type": "time"
                    },
                    {
                      "params": [
                        "host"
                      ],
                      "type": "tag"
                    },
                    {
                      "params": [
                        "null"
                      ],
                      "type": "fill"
                    }
                  ],
                  "measurement": "cpu",
                  "policy": "default",
                  "refId": "A",
                  "resultFormat": "time_series",
                  "select": [
                    [
                      {
                        "params": [
                          "usage_user"
                        ],
                        "type": "field"
                      },
                      {
                        "params": [],
                        "type": "mean"
                      }
                    ]
                  ],
                  "tags": []
                }
              ],
              "timeFrom": null,
              "timeShift": null,
              "title": "CPU usage",
              "tooltip": {
                "msResolution": false,
                "shared": true,
                "value_type": "cumulative"
              },
              "type": "graph",
              "xaxis": {
                "show": true
              },
              "yaxes": [
                {
                  "format": "short",
                  "label": null,
                  "logBase": 1,
                  "max": null,
                  "min": null,
                  "show": true
                },
                {
                  "format": "short",
                  "label": null,
                  "logBase": 1,
                  "max": null,
                  "min": null,
                  "show": true
                }
              ]
            }
          ],
          "title": "CPU"
        },
        {
          "collapse": false,
          "editable": true,
          "height": "250px",
          "panels": [
            {
              "aliasColors": {},
              "bars": false,
              "datasource": null,
              "editable": true,
              "error": false,
              "fill": 1,
              "grid": {
                "threshold1": null,
                "threshold1Color": "rgba(216, 200, 27, 0.27)",
                "threshold2": null,
                "threshold2Color": "rgba(234, 112, 112, 0.22)"
              },
              "id": 4,
              "isNew": true,
              "legend": {
                "alignAsTable": true,
                "avg": false,
                "current": false,
                "max": false,
                "min": false,
                "show": true,
                "total": false,
                "values": false
              },
              "lines": true,
              "linewidth": 2,
              "links": [],
              "nullPointMode": "connected",
              "percentage": false,
              "pointradius": 5,
              "points": false,
              "renderer": "flot",
              "seriesOverrides": [],
              "span": 9,
              "stack": false,
              "steppedLine": false,
              "targets": [
                {
                  "alias": "[[tag_host]]",
                  "dsType": "influxdb",
                  "groupBy": [
                    {
                      "params": [
                        "$interval"
                      ],
                      "type": "time"
                    },
                    {
                      "params": [
                        "host"
                      ],
                      "type": "tag"
                    },
                    {
                      "params": [
                        "null"
                      ],
                      "type": "fill"
                    }
                  ],
                  "measurement": "mem",
                  "policy": "default",
                  "query": "SELECT mean(\"usage_user\") FROM \"cpu\" WHERE $timeFilter GROUP BY time($interval), \"host\" fill(null)",
                  "rawQuery": false,
                  "refId": "A",
                  "resultFormat": "time_series",
                  "select": [
                    [
                      {
                        "params": [
                          "used"
                        ],
                        "type": "field"
                      },
                      {
                        "params": [],
                        "type": "mean"
                      }
                    ]
                  ],
                  "tags": []
                }
              ],
              "timeFrom": null,
              "timeShift": null,
              "title": "MEM usage",
              "tooltip": {
                "msResolution": false,
                "shared": true,
                "value_type": "cumulative"
              },
              "type": "graph",
              "xaxis": {
                "show": true
              },
              "yaxes": [
                {
                  "format": "bytes",
                  "label": null,
                  "logBase": 1,
                  "max": null,
                  "min": null,
                  "show": true
                },
                {
                  "format": "short",
                  "label": null,
                  "logBase": 1,
                  "max": null,
                  "min": null,
                  "show": true
                }
              ]
            },
            {
              "cacheTimeout": null,
              "colorBackground": false,
              "colorValue": false,
              "colors": [
                "rgba(245, 54, 54, 0.9)",
                "rgba(237, 129, 40, 0.89)",
                "rgba(50, 172, 45, 0.97)"
              ],
              "datasource": null,
              "editable": true,
              "error": false,
              "format": "percent",
              "gauge": {
                "maxValue": 100,
                "minValue": 0,
                "show": true,
                "thresholdLabels": false,
                "thresholdMarkers": true
              },
              "height": "",
              "id": 3,
              "interval": "1s",
              "isNew": true,
              "links": [],
              "maxDataPoints": 100,
              "nullPointMode": "connected",
              "nullText": null,
              "postfix": "",
              "postfixFontSize": "50%",
              "prefix": "",
              "prefixFontSize": "50%",
              "span": 3,
              "sparkline": {
                "fillColor": "rgba(31, 118, 189, 0.18)",
                "full": false,
                "lineColor": "rgb(31, 120, 193)",
                "show": false
              },
              "targets": [
                {
                  "dsType": "influxdb",
                  "groupBy": [
                    {
                      "params": [
                        "$interval"
                      ],
                      "type": "time"
                    },
                    {
                      "params": [
                        "null"
                      ],
                      "type": "fill"
                    }
                  ],
                  "measurement": "mem",
                  "policy": "default",
                  "query": "SELECT mean(\"used_percent\") FROM \"mem\" WHERE $timeFilter GROUP BY time($interval) fill(null)",
                  "rawQuery": false,
                  "refId": "A",
                  "resultFormat": "time_series",
                  "select": [
                    [
                      {
                        "params": [
                          "used_percent"
                        ],
                        "type": "field"
                      },
                      {
                        "params": [],
                        "type": "mean"
                      }
                    ]
                  ],
                  "tags": []
                }
              ],
              "thresholds": "",
              "title": "Cluster Memory Average",
              "type": "singlestat",
              "valueFontSize": "80%",
              "valueMaps": [
                {
                  "op": "=",
                  "text": "N/A",
                  "value": "null"
                }
              ],
              "valueName": "avg"
            }
          ],
          "title": "MEM"
        }
      ],
      "time": {
        "from": "now-30m",
        "to": "now"
      },
      "timepicker": {
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
        ]
      },
      "templating": {
        "list": []
      },
      "annotations": {
        "list": []
      },
      "refresh": "5s",
      "schemaVersion": 12,
      "version": 40,
      "links": []
  },
  "overwrite": true
}`
