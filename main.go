package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"text/template"

	log "github.com/sirupsen/logrus"
)

var (
	addr             = flag.String("web.listen-address", ":9011", "Address on which to expose metrics and web interface.")
	etherpadURL      = flag.String("etherpad.url", "http://localhost:9001", "URL to connect with Etherpad")
	etherpadAPIToken = flag.String("etherpad.api-token", "", "API Token for Etherpad")
)

type etherpadStats struct {
	MemoryUsage     int          `json:"memoryUsage"`
	MemoryUsageHeap int          `json:"memoryUsageHeap"`
	TotalUsers      int          `json:"totalUsers"`
	ActivePads      int          `json:"activePads"`
	PendingEdits    int          `json:"pendingEdits"`
	HttpRequests    httpRequests `json:"httpRequests"`
	Connects        meter        `json:"connects"`
	Edits           edits        `json:"edits"`
}

type etherpadAPIStats struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TotalPads       int `json:"totalPads"`
		TotalSessions   int `json:"totalSessions"`
		TotalActivePads int `json:"totalActivePads"`
	} `json:"data"`
}

type httpRequests struct {
	Meter     meter     `json:"meter"`
	Histogram histogram `json:"histogram"`
}

type edits struct {
	Meter     meter     `json:"meter"`
	Histogram histogram `json:"histogram"`
}

type meter struct {
	Mean              float64 `json:"mean"`
	Count             int     `json:"count"`
	CurrentRate       float64 `json:"currentRate"`
	OneMinuteRate     float64 `json:"1MinuteRate"`
	FiveMinuteRate    float64 `json:"5MinuteRate"`
	FifteenMinuteRate float64 `json:"15MinuteRate"`
}

type histogram struct {
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Sum      float64 `json:"sum"`
	Variance float64 `json:"variance"`
	Mean     float64 `json:"mean"`
	Stddev   float64 `json:"stddev"`
	Count    int     `json:"count"`
	Median   float64 `json:"median"`
	P75      float64 `json:"p75"`
	P95      float64 `json:"p95"`
	P99      float64 `json:"p99"`
	P999     float64 `json:"p999"`
}

var statsTpl = template.Must(template.New("stats").Parse(`# HELP etherpad_memory_usage
# TYPE etherpad_memory_usage gauge
etherpad_memory_usage{type="total"} {{.MemoryUsage}}
etherpad_memory_usage{type="heap"} {{.MemoryUsageHeap}}
# HELP etherpad_total_users
# TYPE etherpad_total_users gauge
etherpad_total_users {{.TotalUsers}}
# HELP etherpad_active_pads
# TYPE etherpad_active_pads gauge
etherpad_active_pads {{.ActivePads}}
# HELP etherpad_http_requests
# TYPE etherpad_http_requests counter
etherpad_http_requests {{.HttpRequests.Meter.Count}}
# HELP etherpad_connects
# TYPE etherpad_connects gauge
etherpad_connects {{.Connects.Count}}
# HELP etherpad_edits
# TYPE etherpad_edits gauge
etherpad_edits {{.Edits.Meter.Count}}
`))

var apiStatsTpl = template.Must(template.New("apiStats").Parse(`# HELP etherpad_total_pads
# TYPE etherpad_total_pads gauge
etherpad_total_pads {{.Data.TotalPads}}
# HELP etherpad_total_sessions
# TYPE etherpad_total_sessions gauge
etherpad_total_sessions {{.Data.TotalSessions}}
# HELP etherpad_total_active_pads
# TYPE etherpad_total_active_pads gauge
etherpad_total_active_pads {{.Data.TotalActivePads}}
`))

type handler struct {
	etherpadURL      string
	etherpadAPIToken string
}

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	res, err := http.Get(h.etherpadURL + "/stats")
	if err != nil {
		log.WithError(err).Error("error while fetching /stats from etherpad")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	var stats etherpadStats
	if err := json.NewDecoder(res.Body).Decode(&stats); err != nil {
		log.WithError(err).Error("error while decoding json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if h.etherpadAPIToken != "" {
		res, err := http.Get(h.etherpadURL + fmt.Sprintf("/api/1.2.14/getStats?apikey=%s", h.etherpadAPIToken))
		if err != nil {
			log.WithError(err).Error("error while fetching /getStats from etherpad api")
		}
		defer res.Body.Close()

		var apiStats etherpadAPIStats
		if err := json.NewDecoder(res.Body).Decode(&apiStats); err != nil {
			log.WithError(err).Error("error while decoding json")
		}

		err = apiStatsTpl.Execute(w, &apiStats)
		if err != nil {
			log.WithError(err).Error("error while executing template for apiStats")
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	err = statsTpl.Execute(w, &stats)
	if err != nil {
		log.WithError(err).Error("error while executing template for stats")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	flag.Parse()

	log.Info("Started Etherpad Metrics Exporter")
	http.Handle("/metrics", handler{etherpadURL: *etherpadURL, etherpadAPIToken: *etherpadAPIToken})
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.WithError(err).Error("error while try to start exporter")
	}
}
