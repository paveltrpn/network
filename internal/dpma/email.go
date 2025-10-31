package dpma

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

var (
	emailProviders = map[string]int{
		"Gmail":       1,
		"Yahoo":       2,
		"Hotmail":     3,
		"MSN":         4,
		"Orange":      5,
		"Comcast":     6,
		"AOL":         7,
		"Live":        8,
		"RediffMail":  9,
		"GMX":         10,
		"Proton Mail": 11,
		"Yandex":      12,
		"Mail.ru":     13,
	}
)

type EmailData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	DeliveryTime int    `json:"delivery_time"`
}

func ReadEmailData(fname string) ([]EmailData, error) {
	var rt []EmailData

	f, err := os.Open(fname)
	if err != nil {
		return []EmailData{}, err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ";")

		if len(tokens) != 3 {
			continue
		}

		if _, ok := emailProviders[tokens[1]]; !ok {
			continue
		}

		if _, ok := alpha2map[tokens[0]]; !ok {
			continue
		}

		dTime, _ := strconv.Atoi(tokens[2])

		rt = append(rt, EmailData{Country: tokens[0], Provider: tokens[1], DeliveryTime: dTime})
	}

	return rt, nil
}
