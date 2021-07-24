package phemex

import (
	"fmt"
	"math"
)

// type ClosedPosition struct {
// 	// Dollar
// 	EntryPrice float64
// 	// crypto
// 	EntryValue float64 // Size * EntryPrice

// 	// Dollar
// 	ExitPrice float64
// 	// crypto
// 	ExitValue float64 // Size * ExitPrice

// 	// crypto
// 	Proceeds float64 // value that is earned (gain/loss) // Earnings, Revenue  exit(entryValue, exitValue)[long -> entryValue - exitValue, short -> exitValue - entryValue]

// 	// proceedsRate -> (exitValue - entryValue) / entryValue

// 	// crypto
// 	Gross float64 // entryValue + proceeds

// 	// Dollar
// 	proceedsAmount float64 // Dollar value at exit price
// 	// Dollar
// 	grossAmount float64 // Dollar value at exit price

// 	// rate
// 	GrossYield float64 // on value -> (gross - entryValue) / entryValue
// 	// Dollar
// 	priceDifference float64 // price Change -> [long -> exitPrice - entryPrice, short -> entryPrice - exitPrice]
// 	// rate
// 	pricePercentage float64 // -> [long -> price difference / entryPrice, short -> price difference / exitPrice]
// 	// Dollar
// 	Size int64 // -> Quantity from contract closed

// 	// Dollar
// 	PNL float64 // proceeds in Dollar amount
// 	// Dollar
// 	Total float64 // gross in Dollar amount
// 	// rate
// 	TotalYield float64 // on Size -> (grossAmount - size) / size

// 	ExchangeFee float64
// 	FundingFee  float64
// 	Net         float64
// }

type ClosedPosition struct {
	Calc
	ContractType string
	Quote        struct { // Contract, Dollar?, USD?
		Entry       float64
		Exit        float64
		PriceChange float64
		Size        int64
		PNL         float64 // not very relevent -> should be called yield?
		Earnings    float64
		Total       float64 // Gross? // need net?
	}
	Value struct { // Settled, Crypto?, Base?
		Entry       float64
		Exit        float64
		PNL         float64 // relevent -> under Revenue struct?
		Earnings    float64 // relevent -> under Revenue struct?
		FundingFee  float64 // not relevent, only revelent when entire position is closed.
		ExchangeFee float64 // ExitFee -> under Revenue struct?
	}
	Rate struct {
		PriceChange float64
		PNL         float64 // relevent
		Yield       float64
		Total       float64
	}
}

type Calc struct {
	*ClosedPosition
}

func (c *Calc) Value() {
	c.ClosedPosition.Value.Entry = value(c.Quote.Size, c.Quote.Entry)
	c.ClosedPosition.Value.Exit = value(c.Quote.Size, c.Quote.Exit)
}

func (c *Calc) Long() {
	c.Calc.Value()
	c.ClosedPosition.Value.PNL = close(c.ClosedPosition.Value.Entry, c.ClosedPosition.Value.Exit)
	c.Close()
}

func (c *Calc) Short() {
	c.Calc.Value()
	c.ClosedPosition.Value.PNL = close(c.ClosedPosition.Value.Exit, c.ClosedPosition.Value.Entry)
	c.Close()
}

func (c *Calc) Close() {
	p := c.ClosedPosition
	v := &p.Value
	q := &p.Quote
	r := &p.Rate

	q.PriceChange = truncate((q.Exit - q.Entry), math.Pow(10, 2))
	v.Earnings = truncate((v.PNL + v.Entry), math.Pow(10, 8))

	q.PNL = truncate((q.Exit * v.PNL), math.Pow(10, 2))
	q.Earnings = truncate(((q.Exit * v.Earnings) - float64(q.Size)), math.Pow(10, 2))

	q.Total = truncate((q.Exit * v.Earnings), math.Pow(10, 2))

	r.PriceChange = truncate((q.PriceChange / q.Entry), math.Pow(10, 5))
	r.PNL = truncate((v.PNL / v.Entry), math.Pow(10, 5)) // [value] or this should be yield?

	r.Yield = truncate((q.PNL / float64(q.Size)), math.Pow(10, 5))                                     // [quote value] // or this should PNL?
	r.Total = truncate((((v.Earnings * q.Exit) - float64(q.Size)) / float64(q.Size)), math.Pow(10, 5)) // dollar value

}

func value(size int64, price float64) float64 {
	factor := math.Pow(10, 8)
	fmt.Println((float64(size) / price))
	return truncate((float64(size) / price), factor)
}

func Value(size int64, price float64) float64 {
	return value(size, price)
}

func close(a, b float64) float64 {
	factor := math.Pow(10, 8)
	return truncate((a - b), factor)
}

func truncate(a, b float64) float64 {
	// truncate
	return math.Floor(a*b) / b
}

func Close(entryPrice, exitPrice float64, size int64) float64 {
	trade := new(ClosedPosition)
	trade.Calc.ClosedPosition = trade

	trade.Quote.Entry = entryPrice
	trade.Quote.Exit = exitPrice
	trade.Quote.Size = size
	trade.Calc.Long()
	// trade.Calc.Short()

	fmt.Print(trade)
	return trade.Value.PNL
}

// Inverse Contracks [BTC/USD]
// Value Amount
//	-> EntryValue 	(Size * Entry)
// 	-> ExitValue 	(Size * Exit)
//	-> PNL			(long -> Entry - Exit, short -> Exit - Entry)  [Amount Accumalated from trade] + should include fees?? maybe not
//	-> Earnings 	(PNL + Entry) [Total -> Total Amount From Trade ,close value] + include ExchangeFee, FundingFee will be tricky if trade was executed on multpile fill orders @ different prices
//	-> FundingFee 	(variable every 8 hours, +0.0100 base) inputs [maybe condsider not part of trade but calculated to balance seperatly as deductions]
//	-> ExchangeFee 	(Maker -0.25, +0.75) inputs / calculated [maybe condsider not part of trade but calculated to balance seperatly, or this is part of trade]
// Rates
//	-> PriceChange 	(long -> (Exit - Entry) / Entry, short -> (Entry - Exit) / Exit)
//	-> PNL 			(long -> value.. (Entry - Exit) / Entry, short -> value.. (Exit - Entry) / Entry)
//	-> Yield 		(Quote.PNL / Quote.Size) [parralle to PriceChange?]
//	-> Total		((value.Earnings * Exit) - Size) / Size [Parralle to PNL?]
// Quote
//	-> EntryPrice (input)
//	-> ExitPrice (input)
//	-> PriceChange (long -> Exit - Entry, short -> Entry - Exit)
//	-> Size (input -> amount at Entry and at exit)
//	-> PNL (Exit * value.PNL) [yield?]
//	-> Earnings  [(((value.PNL + value.Entry) * Quote.Exit) - Quote.Size)] [Total USD Value] [difference of size and PNL] [Earnings? Reveune? Returns?]
//	-> Total (Exit * value.Earnings -> value at exit [amount]) [Total of PNL and Earnings]
