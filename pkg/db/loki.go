package db

import (
	"encoding/json"
	"fmt"
	"github.com/elvin-tacirzade/log-alerting/pkg/config"
	"github.com/elvin-tacirzade/log-alerting/pkg/models"
	"io"
	"net/http"
	"strconv"
	"time"
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
	req, reqErr := http.NewRequest("GET", l.URL, nil)
	if reqErr != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", reqErr)
	}

	q := req.URL.Query()
	q.Add("query", query)
	q.Add("start", startTimeStr)
	req.URL.RawQuery = q.Encode()

	// send the HTTP request
	client := &http.Client{}

	resp, respErr := client.Do(req)
	if respErr != nil {
		return nil, fmt.Errorf("error sending HTTP request: %v", respErr)
	}

	defer resp.Body.Close()

	// read the response body
	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		return nil, fmt.Errorf("failed to read response body: %v", bodyErr)
	}

	// check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status code: %v, got: %v, error: %v", http.StatusOK, resp.StatusCode, string(body))
	}

	// unmarshal json
	var lokiLogEntry models.LokiLogEntry
	if unmarshalErr := json.Unmarshal(body, &lokiLogEntry); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", unmarshalErr)
	}

	return &lokiLogEntry, nil
}
