package tax

import (
	"fmt"
	"strconv"
	"strings"
)

func (t *TransactionEntry) formatCurrency(value float64) string {

	if t.quote() == "USD" {
		return dollarFormat(value)
	}

	return cryptoFormat(value)
}

func dollarFormat(value float64) string {
	dollar := fmt.Sprintf("%.2f", value)
	split := strings.Split(dollar, ".")
	return "$ " + commaSep(split)
}

func cryptoFormat(value float64) string {
	a, err := strconv.ParseFloat(fmt.Sprintf("%.4f", value), 64)
	if err != nil {
		fmt.Println("error")
	}

	b, err := strconv.ParseFloat(fmt.Sprintf("%.8f", value), 64)
	if err != nil {
		fmt.Println("error")
	}

	num := strings.Split(fmt.Sprintf("%f", value), ".")[0]

	var crypto string

	if len(num) <= 5 && a != b {
		crypto = fmt.Sprintf("%.8f", value)
	} else if len(num) <= 7 || a == b {
		crypto = fmt.Sprintf("%.4f", value)
	} else {
		crypto = fmt.Sprintf("%.1f", value)
		fmt.Println(crypto)
		split := strings.Split(crypto, ".")
		split[1] = "0"
		crypto = strings.Join(split, ".")
	}

	return commaSep(strings.Split(crypto, "."))
}

func commaSep(value []string) string {
	offset := len(value[0]) / 3
	size := len(value[0]) + offset
	results := make([]byte, size, size)
	counter := 0

	// adding comma seperation
	for i := len(value[0]) - 1; i >= 0; i-- {
		if (len(value[0])-1-i)%3 == 0 && len(value[0])-1-i != 0 {
			results[i+offset-counter] = ','
			counter++
		}

		results[i-counter+offset] = value[0][i]
	}
	value[0] = string(results)

	return strings.TrimLeft(strings.Join(value, "."), string(byte(0)))
}

// func typeTransaction(t *TransactionEntry) ([]string, map[string]string) {
// 	properties := []string{
// 		"Transaction ID",
// 		// "Order Date",
// 		"Order Pair",
// 		"Order Type",
// 		"Order Price",
// 		"Order Quantity",
// 		"Order Amount",
// 		"USD Price Value"}
// 	// "Fee Amount",}

// 	headerFields := map[string]string{
// 		"Transaction ID":  "ID",
// 		"Order Date":      "Date",
// 		"Order Pair":      "Pair",
// 		"Order Type":      "Type",
// 		"Order Price":     "Price",
// 		"Order Quantity":  "Quantity",
// 		"Order Amount":    "Amount",
// 		"USD Price Value": "Price Value",
// 		"Fee Amount":      "Fee"}

// 	// filter header
// 	// -> logic

// 	// filter body
// 	data := t.filter(properties)

// 	}

// 	return fields, tinyFields
// }

func typeCostBasis() ([]string, map[string]string) {
	fields := []string{
		"Transaction ID",
		"Quote Price Entry",
		"Quote Price Exit",
		"USD Price Entry",
		"USD Price Exit",
		"Allocation -> Quantity",
		"Allocation -> Amount",
		"Allocation -> Value",
		"Balance -> Quantity",
		"Balance -> Amount",
		"Balance -> Value",
		"Holdings -> Balance",
		"Holdings -> unrealized",
		"PNL -> Amount",
		"PNL -> Total"}

	tinyFields := map[string]string{
		"Transaction ID":         "ID",
		"Quote Price Entry":      "Entry",
		"Quote Price Exit":       "Exit",
		"USD Price Entry":        "USD Entry",
		"USD Price Exit":         "USD Exit",
		"Allocation -> Quantity": "Allocate Q",
		"Allocation -> Amount":   "Allocate A",
		"Allocation -> Value":    "Allocate V",
		"Balance -> Quantity":    "Balance Q",
		"Balance -> Amount":      "Balance A",
		"Balance -> Value":       "Balance V",
		"Holdings -> Balance":    "Holdings B",
		"Holdings -> unrealized": "Holdings U",
		"PNL -> Amount":          "PNL A",
		"PNL -> Total":           "PNL T"}

	return fields, tinyFields
}

func createRow(fields []string, width, offset int) []string {

	row := make([]string, len(fields))

	for i := 0; i < len(fields); i++ {
		padRight := offset
		padLeft := width - len(fields[i]) - padRight

		cell := ""
		cell += pad(" ", padLeft)
		cell += fields[i]
		cell += pad(" ", padRight)

		row[i] = cell
	}

	return row
}

func (a *Account) display() {
	properties := []string{
		"Transaction ID",
		// "Order Date",
		"Order Pair",
		"Order Type",
		"Order Price",
		"Order Quantity",
		"Order Amount",
		"USD Price Value"}

	transactions := make([][]string, 0, len(a.Ledger.Transactions))
	// costBases := a.Ledger.CostBases

	// filter transactions
	for i := 0; i < len(a.Ledger.Transactions); i++ {
		transactions = append(transactions, a.Ledger.Transactions[i].filter(properties))
	}

	sheet := createSheet(properties, transactions)

	for i := 0; i < len(sheet); i++ {
		fmt.Println(sheet[i])
	}

	fmt.Println()

}

func createSheet(fields []string, data [][]string) [][]string {

	// OPTIONS:
	// Headder Offset -> 0
	// Data Offset -> 1
	// field Size -> 20

	sheet := make([][]string, 0, len(data)+1)

	// header
	sheet = append(sheet, createRow(fields, 20, 0))

	// body
	for i := 0; i < len(data); i++ {
		sheet = append(sheet, createRow(data[i], 20, 1))
	}

	return sheet
}

func pad(v string, num int) string {
	return strings.Repeat(v, num)
}
