# Etherpad Metrics Exporter

[![Integration](https://github.com/systemli/prometheus-etherpad-exporter/actions/workflows/integration.yaml/badge.svg)](https://github.com/systemli/prometheus-etherpad-exporter/actions/workflows/integration.yaml) [![Quality](https://github.com/systemli/prometheus-etherpad-exporter/actions/workflows/quality.yaml/badge.svg)](https://github.com/systemli/prometheus-etherpad-exporter/actions/workflows/quality.yaml) [![Release](https://github.com/systemli/prometheus-etherpad-exporter/actions/workflows/release.yaml/badge.svg)](https://github.com/systemli/prometheus-etherpad-exporter/actions/workflows/release.yaml)

Prometheus Exporter for Etherpad written in Go.

## Usage

```
go get github.com/systemli/prometheus-etherpad-exporter
go install github.com/systemli/prometheus-etherpad-exporter
$GOPATH/bin/prometheus-etherpad-exporter
```

### Docker

```
docker run -p 9011:9011 systemli/prometheus-etherpad-exporter:latest -etherpad.url http://localhost:9001 
```

## Metrics

```
# HELP etherpad_memory_usage
# TYPE etherpad_memory_usage gauge
etherpad_memory_usage{type="total"} 102801408
etherpad_memory_usage{type="heap"} 30452280
# HELP etherpad_total_users
# TYPE etherpad_total_users gauge
etherpad_total_users 1
# HELP etherpad_active_pads
# TYPE etherpad_active_pads gauge
etherpad_active_pads 1
# HELP etherpad_http_requests
# TYPE etherpad_http_requests counter
etherpad_http_requests 92
# HELP etherpad_connects
# TYPE etherpad_connects gauge
etherpad_connects 1
# HELP etherpad_edits
# TYPE etherpad_edits gauge
etherpad_connects 3
```

## License

GPLv3
