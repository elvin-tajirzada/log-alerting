package db

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/elvin-tajirzada/log-alerting/internal/config"
	"github.com/elvin-tajirzada/log-alerting/internal/models"
)

const query = "{job=\"log-exporter\"} | json | status =~ `^5\\d\\d$`"

type Loki struct {
	URL string
}

func NewLoki(conf *config.Config) *Loki {
	return &Loki{
		URL: conf.Loki.URL,
	}
}

func (l *Loki) Get(evaluationTime time.Duration) (*models.LokiLogEntry, error) {
	startTime := time.Now().Add(-evaluationTime)
	startTimeStr := strconv.FormatInt(startTime.UnixNano(), 10)

	// create the HTTP request
	req, err := http.NewRequest("GET", l.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	q := req.URL.Query()
	q.Add("query", query)
	q.Add("start", startTimeStr)
	req.URL.RawQuery = q.Encode()

	// send the HTTP request
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %v", err)
	}

	defer resp.Body.Close()

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status code: %v, got: %v, error: %v", http.StatusOK, resp.StatusCode, string(body))
	}

	// unmarshal json
	var lokiLogEntry models.LokiLogEntry
	if err := json.Unmarshal(body, &lokiLogEntry); err != nil {
		return nil, fmt.Errorf("unable to unmarshal json: %v", err)
	}

	return &lokiLogEntry, nil
}
