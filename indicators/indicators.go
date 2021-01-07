package indicators

import (
	"github.com/shopspring/decimal"
)

// EMA EMA指标
type EMA struct {
	lastEMA 		*decimal.Decimal
}

// Calculate .
func (e *EMA) Calculate(price float64, period int32) float64 {
	if e.lastEMA == nil {
		tmp :=  decimal.NewFromInt32(period)
		e.lastEMA = &tmp
		return price
	}
	one := decimal.NewFromInt32(1)
	two := decimal.NewFromInt32(2)
	p1 := decimal.NewFromFloat(price)
	p2 := decimal.NewFromInt32(period)
	r := p1.Mul(two).Add(e.lastEMA.Mul(p2.Sub(one))).Div(p2.Add(one))
	e.lastEMA = &r
	result, _ := r.Round(4).Float64()
	return result
}

// MACD MACD indicator
type MACD struct {
	MACD       	[]float64
	DEA 		[]float64
	DIF 		[]float64
	emaShort	decimal.Decimal
	emaLong  	decimal.Decimal
	short 		decimal.Decimal
	long 		decimal.Decimal
	mid 		decimal.Decimal
	lastDEA     decimal.Decimal
}

// Init .
func (m *MACD) Init(short, long, mid int32) {
	m.short = decimal.NewFromInt32(short)
	m.long = decimal.NewFromInt32(long)
	m.mid = decimal.NewFromInt32(mid)
}

// Next .
func (m *MACD) Next(c float64) {
	one := decimal.NewFromFloat(1)
	two := decimal.NewFromFloat(2)
	close := decimal.NewFromFloat(c)
	if len(m.MACD) == 0 {
		m.emaShort, m.emaLong = close, close
		m.DEA = append(m.DEA, 0)
		m.DIF = append(m.DIF, 0)
		m.MACD = append(m.MACD, 0)
		return
	}
	m.emaShort = two.Mul(close).Add(m.short.Sub(one).Mul(m.emaShort)).Div(m.short.Add(one))
	m.emaLong = two.Mul(close).Add(m.long.Sub(one).Mul(m.emaLong)).Div(m.long.Add(one))
	dif := m.emaShort.Sub(m.emaLong)
	m.lastDEA = two.Mul(dif).Add(m.mid.Sub(one).Mul(m.lastDEA)).Div(m.mid.Add(one))
	macd, _ := dif.Sub(m.lastDEA).Round(4).Float64()
	tmp1, _ := dif.Round(4).Float64()
	tmp2, _ := m.lastDEA.Round(4).Float64()
	m.DIF = append(m.DIF, tmp1)
	m.DEA = append(m.DEA, tmp2)
	m.MACD = append(m.MACD, macd)
}

// MA MA指标
func MA(slice []float64) float64 {
	if len(slice) == 0 {
		return 0
	}
	var total decimal.Decimal
	l := decimal.NewFromInt(int64((len(slice))))
	for _, value := range slice {
		val := decimal.NewFromFloat(value)
		total = total.Add(val)
	}
	lenght, _ := total.Div(l).Round(4).Float64()
	return lenght
}

// Round 保留小数位
func Round(v float64) float64 {
	x := decimal.NewFromFloat(v).Round(4)
	result, _ := x.Float64()
	return result
}