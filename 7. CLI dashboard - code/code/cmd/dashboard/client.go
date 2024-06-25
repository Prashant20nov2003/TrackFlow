package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gotracker"
	"io"
	"net/http"
)

func getMetric(what gotracker.QueryType) ([]gotracker.Metric, error) {
	data := gotracker.MetricData{
		What:   what,
		SiteID: siteID,
		Start:  uint32(start),
		End:    uint32(end),
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "http://localhost:9876/stats", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var metrics []gotracker.Metric
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		fmt.Println("error from API: ", string(b))
		return nil, err
	}

	return metrics, nil
}
