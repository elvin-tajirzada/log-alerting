![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/elvin-tacirzade/log-alerting?logo=go)
[![Go Reference](https://pkg.go.dev/badge/github.com/elvin-tacirzade/log-alerting.svg)](https://pkg.go.dev/github.com/elvin-tacirzade/log-alerting)
![Docker Pulls](https://img.shields.io/docker/pulls/elvintacirzade/log-alerting?logo=docker&logoColor=white)
![Docker Image Size (tag)](https://img.shields.io/docker/image-size/elvintacirzade/log-alerting/latest?logo=docker&logoColor=white)

# Log Alerting

Log Alerting makes it possible to send error logs to telegram.

## Overview

Log Alerting depends on [Log Exporter](https://github.com/elvin-tacirzade/log-exporter) service. The log structure
should be the same. Log Alerting sends a request to [Loki](https://grafana.com/oss/loki/) at each evaluation time, so it
retrieves logs where the status code equals 5xx and, it sends messages to telegram. It also sends a daily report to
telegram.

<table>
    <tr>
      <td>
        <img alt="Server Error" src="https://github.com/elvin-tacirzade/log-alerting/blob/main/photos/server_error.jpg?raw=true">
      </td>
      <td>
        <img alt="Daily Report" src="https://github.com/elvin-tacirzade/log-alerting/blob/main/photos/daily_report.jpg?raw=true">
      </td>
      <td>
        <img alt="Daily Report" src="https://github.com/elvin-tacirzade/log-alerting/blob/main/photos/daily_report_no_error.jpg?raw=true">
      </td>
    </tr>
</table>

## Usage

We use the following command to run it on [Docker](https://www.docker.com/).

```
docker run -d \
  --name log-alerting \
  --network main \
  --env LOKI_HOST=your_loki_host \
  --env LOKI_PORT=your_loki_port \
  --env TELEGRAM_TOKEN=your_telegram_token \
  --env TELEGRAM_CHAT_ID=your_telegram_chat_id \
  --env GRAFANA_PANEL_URL=your_grafana_panel_url \
  elvintacirzade/log-alerting:latest
```

[See](https://hub.docker.com/r/elvintacirzade/log-alerting) more information.




