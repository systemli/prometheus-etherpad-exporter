# Etherpad Metrics Exporter

[![Integration](https://github.com/systemli/prometheus-etherpad-exporter/actions/workflows/integration.yaml/badge.svg)](https://github.com/systemli/prometheus-etherpad-exporter/actions/workflows/integration.yaml) [![Quality](https://github.com/systemli/prometheus-etherpad-exporter/actions/workflows/quality.yaml/badge.svg)](https://github.com/systemli/prometheus-etherpad-exporter/actions/workflows/quality.yaml) [![Release](https://github.com/systemli/prometheus-etherpad-exporter/actions/workflows/release.yaml/badge.svg)](https://github.com/systemli/prometheus-etherpad-exporter/actions/workflows/release.yaml)

Prometheus Exporter for Etherpad written in Go.

## Usage

```
go install github.com/systemli/prometheus-etherpad-exporter@latest
$GOPATH/bin/prometheus-etherpad-exporter
```

### Commandline options

```
-web.listen-address ":9011" # Address on which to expose metrics and web interface.
-etherpad.url "http://localhost:9001" # URL to connect with Etherpad
-etherpad.api-token "" # "API Token for Etherpad"
```

With configured API Token the metrics `etherpad_total_pads`, `etherpad_total_sessions` and `etherpad_total_active_pads` will appended to the metrics

### Docker

```
docker run -p 9011:9011 systemli/prometheus-etherpad-exporter:latest -etherpad.url http://localhost:9001 
```

## Metrics

```
# HELP etherpad_total_pads
# TYPE etherpad_total_pads gauge
etherpad_total_pads 8
# HELP etherpad_total_sessions
# TYPE etherpad_total_sessions gauge
etherpad_total_sessions 0
# HELP etherpad_total_active_pads
# TYPE etherpad_total_active_pads gauge
etherpad_total_active_pads 0
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
# HELP etherpad_disconnects
# TYPE etherpad_disconnects gauge
etherpad_connects 0
# HELP etherpad_edits
# TYPE etherpad_edits gauge
etherpad_edits 3
# HELP etherpad_failed_changesets
# TYPE etherpad_failed_changesets gauge
etherpad_failed_changesets 0
# HELP etherpad_ueberdb_locks
# TYPE etherpad_ueberdb_locks gauge
etherpad_ueberdb_locks{state="awaits"} 0
etherpad_ueberdb_locks{state="acquires"} 0
etherpad_ueberdb_locks{state="releases"} 0
# HELP etherpad_ueberdb_reads
# TYPE etherpad_ueberdb_reads gauge
etherpad_ueberdb_reads 0
etherpad_ueberdb_reads{state="failed"} 0
etherpad_ueberdb_reads{state="finished"} 0
etherpad_ueberdb_reads{state="from_cache"} 0
etherpad_ueberdb_reads{state="from_db"} 0
etherpad_ueberdb_reads{state="from_db_failed"} 0
etherpad_ueberdb_reads{state="from_db_finished"} 0
# HELP etherpad_ueberdb_writes
# TYPE etherpad_ueberdb_writes gauge
etherpad_ueberdb_writes 0
etherpad_ueberdb_writes{state="failed"} 0
etherpad_ueberdb_writes{state="finished"} 0
etherpad_ueberdb_writes{state="obsoleted"} 0
etherpad_ueberdb_writes{state="to_db"} 0
etherpad_ueberdb_writes{state="to_db_failed"} 0
etherpad_ueberdb_writes{state="to_db_finished"} 0
```

## License

GPLv3
