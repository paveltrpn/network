package dpma

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

var (
	voiceCallProviders = map[string]int{
		"TransparentCalls": 1,
		"E-Voice":          2,
		"JustPhone":        3,
	}
)

type VoiceCallData struct {
	Country           string  `json:"country"`
	Bandwith          int     `json:"bandwidth"`
	AvgRespTime       int     `json:"response_time"`
	Provider          string  `json:"provider"`
	ConnStability     float32 `json:"connection_stability"`
	Frequency         int     `json:"ttfb"`
	VoicePurity       int     `json:"voice_purity"`
	ConnectionTimeMed int     `json:"median_of_calls_time"`
}

func ReadVoiceCallData(fname string) ([]VoiceCallData, error) {
	var rt []VoiceCallData

	f, err := os.Open(fname)
	if err != nil {
		return []VoiceCallData{}, err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ";")

		if len(tokens) != 8 {
			continue
		}

		if _, ok := voiceCallProviders[tokens[3]]; !ok {
			continue
		}

		if _, ok := alpha2map[tokens[0]]; !ok {
			continue
		}

		bandwith, err := strconv.Atoi(tokens[1])
		if err != nil {
			return []VoiceCallData{}, err
		}

		ansTime, err := strconv.Atoi(tokens[2])
		if err != nil {
			return []VoiceCallData{}, err
		}

		stab, err := strconv.ParseFloat(tokens[4], 32)
		if err != nil {
			return []VoiceCallData{}, err
		}

		freq, err := strconv.Atoi(tokens[5])
		if err != nil {
			return []VoiceCallData{}, err
		}

		voice, err := strconv.Atoi(tokens[6])
		if err != nil {
			return []VoiceCallData{}, err
		}

		contime, err := strconv.Atoi(tokens[7])
		if err != nil {
			return []VoiceCallData{}, err
		}

		rt = append(rt, VoiceCallData{Country: tokens[0], Bandwith: bandwith, AvgRespTime: ansTime, Provider: tokens[3],
			ConnStability: float32(stab), Frequency: freq, VoicePurity: voice, ConnectionTimeMed: contime})
	}

	return rt, nil
}
