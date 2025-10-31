package dpma

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"` // возможные статусы active и closed
}

func RequestIncidentData(url string) ([]IncidentData, error) {
	var rt []IncidentData

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewReader([]byte{}))
	if err != nil {
		return []IncidentData{}, err
	}
	defer req.Body.Close()

	req.Header = http.Header{
		"Content-Type": {"application/json"},
		"Charset":      {"utf-8"},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []IncidentData{}, err
	}

	if res.StatusCode != http.StatusOK {
		return []IncidentData{}, fmt.Errorf("RequestIncidentData(): Warning! response - %v", res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return []IncidentData{}, err
	}

	err = json.Unmarshal(resBody, &rt)
	if err != nil {
		return []IncidentData{}, err
	}

	return rt, nil
}
