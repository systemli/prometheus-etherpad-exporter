package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"text/template"
)

var (
	addr        = flag.String("web.listen-address", ":9011", "Address on which to expose metrics and web interface.")
	etherpadURL = flag.String("etherpad.url", "http://localhost:9001", "URL to connect with Etherpad")
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

var tpl = template.Must(template.New("stats").Parse(`# HELP etherpad_memory_usage
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
etherpad_connects {{.Edits.Meter.Count}}
`))

type handler struct {
	etherpadURL string
}

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := http.Get(h.etherpadURL + "/stats")
	if err != nil {
		log.Printf("scrape error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var stats etherpadStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		log.Printf("json decoding error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	err = tpl.Execute(w, &stats)
	if err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	http.Handle("/metrics", handler{etherpadURL: *etherpadURL})
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}

	log.Println("Started Etherpad Metrics Exporter")
}
