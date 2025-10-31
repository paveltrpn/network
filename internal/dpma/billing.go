package dpma

import (
	"bufio"
	"math"
	"os"
	"strconv"
)

type BillingData struct {
	CreateCustomer bool
	Purchase       bool
	Payout         bool
	Recurring      bool
	FraudControl   bool
	CheckoutPage   bool
}

func checkBit(bits uint8, atPos uint) bool {
	val := bits & (1 << atPos)
	return (val > 0)
}

func ReadBillingData(fname string) (BillingData, error) {
	var (
		accum uint8
		rt    BillingData
	)

	f, err := os.Open(fname)
	if err != nil {
		return BillingData{}, err
	}

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	line := scanner.Text()

	for i := 5; i > -1; i-- {
		bit, _ := strconv.Atoi(string(line[i]))
		accum += uint8(math.Pow(2, float64(5-i))) * uint8(bit)
	}

	rt.CreateCustomer = checkBit(accum, 0)
	rt.Purchase = checkBit(accum, 1)
	rt.Payout = checkBit(accum, 2)
	rt.Recurring = checkBit(accum, 3)
	rt.FraudControl = checkBit(accum, 4)
	rt.CheckoutPage = checkBit(accum, 5)

	return rt, nil
}
