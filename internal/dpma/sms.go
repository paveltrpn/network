package dpma

import (
	"bufio"
	"os"
	"strings"
)

var (
	smsProviders = map[string]int{
		"Topolo": 1,
		"Rond":   2,
		"Kildy":  3,
	}
)

type SMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

func ReadSMSData(fname string) ([]SMSData, error) {
	var rt []SMSData

	f, err := os.Open(fname)
	if err != nil {
		return []SMSData{}, err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ";")

		if len(tokens) != 4 {
			continue
		}

		if _, ok := smsProviders[tokens[3]]; !ok {
			continue
		}

		if _, ok := alpha2map[tokens[0]]; !ok {
			continue
		}

		rt = append(rt, SMSData{Country: tokens[0], Bandwidth: tokens[1], ResponseTime: tokens[2], Provider: tokens[3]})
	}

	return rt, nil
}
