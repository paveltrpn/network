package dpma

import (
	"sort"
)

type ResultT struct {
	Status bool       `json:"status"`
	Data   ResultSetT `json:"data"`
	Error  string     `json:"error"`
}

type ResultSetT struct {
	SMS       [][]SMSData              `json:"sms"`
	MMS       [][]MMSData              `json:"mms"`
	VoiceCall []VoiceCallData          `json:"voice_call"`
	Email     map[string][][]EmailData `json:"email"`
	Billing   BillingData              `json:"billing"`
	Support   []int                    `json:"support"`
	Incidents []IncidentData           `json:"incident"`
}

const (
	assetsDir string = "../../assets/simulator/"
)

func CollectSMSData(fname string) ([][]SMSData, error) {
	var (
		smsDataFullCountryName []SMSData
		rt                     [][]SMSData
	)

	smsDataRaw, err := ReadSMSData(fname)
	if err != nil {
		return [][]SMSData{}, err
	}

	for _, sms := range smsDataRaw {
		smsDataFullCountryName = append(smsDataFullCountryName,
			SMSData{Country: alpha2map[sms.Country], Bandwidth: sms.Bandwidth, ResponseTime: sms.ResponseTime, Provider: sms.Provider})
	}

	sort.Slice(smsDataFullCountryName, func(i, j int) bool {
		return smsDataFullCountryName[i].Country < smsDataFullCountryName[j].Country
	})

	sort.Slice(smsDataRaw, func(i, j int) bool {
		return smsDataRaw[i].Provider < smsDataRaw[j].Provider
	})

	rt = append(rt, smsDataFullCountryName, smsDataRaw)

	return rt, nil
}

func CollectMMSData(url string) ([][]MMSData, error) {
	var (
		mmsDataFullCountryName []MMSData
		rt                     [][]MMSData
	)

	mmsDataRaw, err := RequestMMSData(url)
	if err != nil {
		return [][]MMSData{}, err
	}

	for _, mms := range mmsDataRaw {
		mmsDataFullCountryName = append(mmsDataFullCountryName,
			MMSData{Country: alpha2map[mms.Country], Bandwidth: mms.Bandwidth, ResponseTime: mms.ResponseTime, Provider: mms.Provider})
	}

	sort.Slice(mmsDataFullCountryName, func(i, j int) bool {
		return mmsDataFullCountryName[i].Country < mmsDataFullCountryName[j].Country
	})

	sort.Slice(mmsDataRaw, func(i, j int) bool {
		return mmsDataRaw[i].Provider < mmsDataRaw[j].Provider
	})

	rt = append(rt, mmsDataFullCountryName, mmsDataRaw)

	return rt, nil
}

func CollectEmailData(fname string) (map[string][][]EmailData, error) {
	email, err := ReadEmailData(fname)
	if err != nil {
		return nil, err
	}

	// preliminarily sort email data slice
	sort.Slice(email, func(i, j int) bool {
		return email[i].Country < email[j].Country
	})

	byCountry := make(map[string][][]EmailData)
	var providers []EmailData

	countryCode := email[0].Country

	// Arrange providers by countries. Make a map with
	// key of country code and value of slice with EmailData
	// strucures contents of providers that work in corresponding country
	for _, e := range email {
		if e.Country == countryCode {
			providers = append(providers, e)
		} else {
			sort.Slice(providers, func(i, j int) bool {
				return providers[i].DeliveryTime < providers[j].DeliveryTime
			})

			fastest := providers[:2]
			slowest := providers[len(providers)-3:]

			provData := make([][]EmailData, 0)

			provData = append(provData, fastest)
			provData = append(provData, slowest)

			byCountry[countryCode] = provData
			countryCode = e.Country

			providers = make([]EmailData, 0)
		}
	}

	return byCountry, nil
}

func CollectSupportData(url string) ([]int, error) {
	var (
		rt []int
	)

	support, err := RequestSupportData(url)
	if err != nil {
		return []int{}, err
	}

	var (
		totalTickets int
		load         int
	)

	for _, sup := range support {
		totalTickets += sup.ActiveTickets
	}

	switch {
	case totalTickets < 9:
		load = 1
	case totalTickets > 16:
		load = 3
	default:
		load = 2
	}
	rt = append(rt, load)
	rt = append(rt, (60/18)*totalTickets)

	return rt, nil
}

func CollectIncidentData(url string) ([]IncidentData, error) {
	// 8484 port and incident in description, accendent in simulator
	incident, err := RequestIncidentData(url)
	if err != nil {
		return []IncidentData{}, err
	}

	sort.Slice(incident, func(i, j int) bool {
		return incident[i].Status < incident[j].Status
	})

	return incident, nil
}
