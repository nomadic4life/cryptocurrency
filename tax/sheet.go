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

func formatCurrency(value float64, quote string) string {

	if quote == "USD" {
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
	properties := make([][]string, 2)
	sheet := make([][][]string, 2)
	header := make([]table, 2)
	fields := make([]table, 2)

	properties[0] = []string{
		"Transaction ID",
		// "Order Date",
		"Order Pair",
		"Order Type",
		"Order Price",
		"Order Quantity",
		"Order Amount",
		"USD Price Value"}

	properties[1] = []string{
		"Transaction ID -> Credit",
		"Transaction ID -> Debit",
		"Quote Price -> Entry",
		"Quote Price -> Exit",
		"USD Price -> Entry",
		"USD Price -> Exit",
		"Allocation -> Quantity",
		"Allocation -> Amount",
		"Allocation -> Value",
		"Balance -> Quantity",
		"Balance -> Amount",
		"Balance -> Value",
		"Holdings -> Balance",
		// "Holdings -> unrealized",
		"PNL -> Amount",
		"PNL -> Total"}

	header[0] = map[string]string{
		"Transaction ID":  "TRANSACTION",
		"Order Date":      "ORDER",
		"Order Pair":      "ORDER",
		"Order Type":      "ORDER",
		"Order Price":     "ORDER",
		"Order Quantity":  "ORDER",
		"Order Amount":    "ORDER",
		"USD Price Value": "USD PRICE",
		"Fee Amount":      "FEE"}

	header[1] = map[string]string{
		"Transaction ID -> Credit": "ID",
		"Transaction ID -> Debit":  "ID",
		"Quote Price -> Entry":     "QUOTE PRICE",
		"Quote Price -> Exit":      "QUOTE PRICE",
		"USD Price -> Entry":       "USD PRICE",
		"USD Price -> Exit":        "USD PRICE",
		"Allocation -> Quantity":   "ALLOCATION",
		"Allocation -> Amount":     "ALLOCATION",
		"Allocation -> Value":      "ALLOCATION",
		"Balance -> Quantity":      "BALANCE",
		"Balance -> Amount":        "BALANCE",
		"Balance -> Value":         "BALANCE",
		"Holdings -> Balance":      "HOLDINGS",
		"Holdings -> Unrealized":   "HOLDINGS",
		"PNL -> Amount":            "PNL",
		"PNL -> Total":             "PNL"}

	fields[0] = map[string]string{
		"Transaction ID":  "ID",
		"Order Pair":      "PAIR",
		"Order Type":      "TYPE",
		"Order Price":     "PRICE",
		"Order Quantity":  "QAUNTITY",
		"Order Amount":    "AMOUNT",
		"USD Price Value": "VALUE",
		"Fee Amount":      "FEE"}

	fields[1] = map[string]string{
		"Transaction ID -> Credit": "CREDIT",
		"Transaction ID -> Debit":  "DEBIT",
		"Quote Price -> Entry":     "ENTRY",
		"Quote Price -> Exit":      "EXIT",
		"USD Price -> Entry":       "ENTRY",
		"USD Price -> Exit":        "EXIT",
		"Allocation -> Quantity":   "QUANTITY",
		"Allocation -> Amount":     "AMOUNT",
		"Allocation -> Value":      "VALUE",
		"Balance -> Quantity":      "QUANTITY",
		"Balance -> Amount":        "AMOUNT",
		"Balance -> Value":         "VALUE",
		"Holdings -> Balance":      "BALANCE",
		"Holdings -> Unrealized":   "UNREALIZED",
		"PNL -> Amount":            "AMOUNT",
		"PNL -> Total":             "TOTAL"}

	fmt.Print("\n")
	fmt.Println("STATEMENT:", a.Statement)
	fmt.Println("ASSETHOLDINGS:", a.Assets.Holdings)
	fmt.Println("LEDGER -> TRANSACTION:")

	transactions := make([][]string, 0, len(a.Ledger.Transactions))

	// filter transactions
	for i := 0; i < len(a.Ledger.Transactions); i++ {
		transactions = append(transactions, a.Ledger.Transactions[i].filter(properties[0]))
	}

	// sheet[0] = createSheet(properties[0], transactions)
	// fmt.Println(header[0])
	sheet[0] = createSheet(header[0].filter(properties[0]), fields[0].filter(properties[0]), transactions)

	// display transactions
	fmt.Println()
	for i := 0; i < len(sheet[0]); i++ {
		fmt.Println(sheet[0][i])
		if i%3 == 0 && i != 0 {
			fmt.Println()
		}
	}

	fmt.Println()
	fmt.Println("LEDGER -> COST BASIS:")

	costBases := make([][]string, 0, len(a.Ledger.CostBases))

	// filter costBases
	for i := 0; i < len(a.Ledger.CostBases); i++ {
		costBases = append(costBases, a.Ledger.CostBases[i].filter(properties[1]))
	}

	// sheet[0] = createSheet(properties[0], transactions)
	sheet[1] = createSheet(header[1].filter(properties[1]), fields[1].filter(properties[1]), costBases)

	// display costBases
	fmt.Println()
	for i := 0; i < len(sheet[1]); i++ {
		fmt.Println(sheet[1][i])
		if i%3 == 0 && i != 0 {
			fmt.Println()
		}
	}

	fmt.Println()
}

func createSheet(header, fields []string, data [][]string) [][]string {

	// OPTIONS:
	// Headder Offset -> 0
	// Data Offset -> 1
	// field Size -> 20

	sheet := make([][]string, 0, len(data)+2)

	// main header
	if len(header) == 1 {
		sheet = append(sheet, createRow(header, 200, 80))

	} else {
		sheet = append(sheet, createRow(header, 16, 0))
	}

	// sub header
	sheet = append(sheet, createRow(fields, 16, 0))

	// body
	for i := 0; i < len(data); i++ {
		sheet = append(sheet, createRow(data[i], 16, 1))
	}

	return sheet
}

func pad(v string, num int) string {
	return strings.Repeat(v, num)
}

var propertyWidths map[string]map[string]int = map[string]map[string]int{
	// 268
	"Transaction ID": {
		"From": 6,
		"To":   6},
	"Quote Price": {
		"Entry": 14,
		"Exit":  14},
	"USD Price": {
		"Entry": 15,
		"Exit":  15},
	"Allocation": {
		"Quantity": 18,
		"Amount":   18,
		"Value":    14},
	"Balance": {
		"Quantity": 18,
		"Amount":   18,
		"Value":    14},
	"Holdings": {
		"Balance":    20,
		"Unrealized": 20},
	"PNL": {
		"Amount": 14,
		"Total":  20},
	// 137
	"Transactions": {
		"ID":       8,
		"Date":     20,
		"Pair":     11,
		"Type":     8,
		"Price":    18,
		"Quantity": 18,
		"Amount":   18,
		"Value":    18,
		"Fee":      18}}
