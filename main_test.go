package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type constHandler struct {
	s string
}

func (h constHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.s))
}

func TestGetMetrics(t *testing.T) {
	tcs := []struct {
		statsJson string
		expected  string
	}{
		{
			statsJson: `{"httpStartTime":1616252622240,"memoryUsage":194924544,"memoryUsageHeap":36845824,"totalUsers":1,"httpRequests":{"meter":{"mean":0.19830325431417883,"count":15,"currentRate":0.27914103077088026,"1MinuteRate":0.15421117090001737,"5MinuteRate":0.04501349471844342,"15MinuteRate":0.01608537818128321},"histogram":{"min":0.8717880100011826,"max":12.836927011609077,"sum":41.26530301570892,"variance":10.038097927388511,"mean":2.7510202010472615,"stddev":3.1682957449374123,"count":15,"median":1.572950005531311,"p75":2.5995219945907593,"p95":12.836927011609077,"p99":12.836927011609077,"p999":12.836927011609077}},"connects":{"mean":0.1699089576647754,"count":1,"currentRate":0.16990908200350413,"1MinuteRate":0.015991117074135343,"5MinuteRate":0.0033057092356765017,"15MinuteRate":0.001108030399020654},"activePads":1,"pendingEdits":0,"edits":{"meter":{"mean":1.0212670241700796,"count":3,"currentRate":1.0212673798266378,"1MinuteRate":0,"5MinuteRate":0,"15MinuteRate":0},"histogram":{"min":1.6878540068864822,"max":2.3621450066566467,"sum":5.915806010365486,"variance":0.12211450650073814,"mean":1.9719353367884953,"stddev":0.34944886106659173,"count":3,"median":1.8658069968223572,"p75":2.3621450066566467,"p95":2.3621450066566467,"p99":2.3621450066566467,"p999":2.3621450066566467}}}`,
			expected: `# HELP etherpad_memory_usage
# TYPE etherpad_memory_usage gauge
etherpad_memory_usage{type="total"} 194924544
etherpad_memory_usage{type="heap"} 36845824
# HELP etherpad_total_users
# TYPE etherpad_total_users gauge
etherpad_total_users 1
# HELP etherpad_active_pads
# TYPE etherpad_active_pads gauge
etherpad_active_pads 1
# HELP etherpad_http_requests
# TYPE etherpad_http_requests counter
etherpad_http_requests 15
# HELP etherpad_connects
# TYPE etherpad_connects gauge
etherpad_connects 1
# HELP etherpad_edits
# TYPE etherpad_edits gauge
etherpad_connects 3
`,
		},
	}

	for _, tc := range tcs {
		srv := httptest.NewServer(constHandler{tc.statsJson})

		h := handler{
			etherpadURL: srv.URL,
		}
		req, err := http.NewRequest("GET", "/metrics", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if rr.Body.String() != tc.expected {
			t.Errorf("Response does not match the expected string:\n%s", cmp.Diff(rr.Body.String(), tc.expected))
		}

		srv.Close()
	}
}
