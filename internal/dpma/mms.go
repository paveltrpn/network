package dpma

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

func RequestMMSData(url string) ([]MMSData, error) {
	var rt []MMSData

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewReader([]byte{}))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	req.Header = http.Header{
		"Content-Type": {"application/json"},
		"Charset":      {"utf-8"},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []MMSData{}, err
	}

	if res.StatusCode != http.StatusOK {
		return []MMSData{}, fmt.Errorf("RequestMMSData(): Warning! response - %v", res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return []MMSData{}, err
	}

	err = json.Unmarshal(resBody, &rt)
	if err != nil {
		return []MMSData{}, err
	}

	for i := range rt {
		if _, ok := alpha2map[rt[i].Country]; !ok {
			rt[i] = rt[len(rt)-1]
			rt = rt[:len(rt)-1]
		}
	}

	return rt, nil
}
