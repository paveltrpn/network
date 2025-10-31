package dpma

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SupportData struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}

func RequestSupportData(url string) ([]SupportData, error) {
	var rt []SupportData

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewReader([]byte{}))
	if err != nil {
		return []SupportData{}, err
	}
	defer req.Body.Close()

	req.Header = http.Header{
		"Content-Type": {"application/json"},
		"Charset":      {"utf-8"},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []SupportData{}, err
	}

	if res.StatusCode != http.StatusOK {
		return []SupportData{}, fmt.Errorf("RequestSupportData(): Warning! response - %v", res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return []SupportData{}, err
	}

	err = json.Unmarshal(resBody, &rt)
	if err != nil {
		return []SupportData{}, err
	}

	return rt, nil
}
