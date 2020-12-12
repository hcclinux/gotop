package utils


import (
	"time"
)

const (
	// Error 订单异常
	Error = 	iota
	// Canceled 订单取消状态
	Canceled 		
	// Rejected 订单拒绝状态
	Rejected
	// Submitted 订单提交状态
	Submitted
	// Accepted 订单接受状态
	Accepted
	// Completed 订单完成状态
	Completed
	// Expired 订单过期状态
	Expired
	// Partial 订单部分成交状态
	Partial
	// Margin 订单余量状态
	Margin
)

const (
	// Market 市价单
	Market		=	"MARKET"
	// Limit 限价单
	Limit		= 	"LIMIT"
	// Long 多头持仓方向
	Long		= 	"LONG"
	// Short 空头持仓方向
	Short		= 	"SHORT"
	// Buy 买入
	Buy 		= 	"BUY"
	// Sell 卖出
	Sell 		= 	"SELL"
)

/*
KLine 用于实现K线的显示和数据的操作
*/
type KLine struct {
	OHLC 		[]Candle
}



// Append 添加candle数据
func (k *KLine) Append(c Candle) {
	k.OHLC = append(k.OHLC, c)
}

/*
Candle kline的基础数据
*/
type Candle struct {
	Date		time.Time					`json:"datetime" csv:"date"`
	Open		float64						`json:"open" csv:"open"`
	High 		float64						`json:"high" csv:"high"`
	Low			float64						`json:"low" csv:"low"`
	Close 		float64						`json:"close" csv:"close"`
	Volume		float64						`json:"volume" csv:"volume"`
}

// IndexCandle 1
type IndexCandle struct {
	Date		time.Time					`json:"datetime" csv:"date"`
	Open		float64						`json:"open" csv:"open"`
	High 		float64						`json:"high" csv:"high"`
	Low			float64						`json:"low" csv:"low"`
	Close 		float64						`json:"close" csv:"close"`
	Volume		float64						`json:"volume" csv:"volume"`
	Index 		int							`json:"index" csv:"index"`
}

// Order 订单详情
type Order struct {
	Symbol				string
	Time 				time.Time   		`json:"time"`
	OrderID 			string
	ExecPrice			float64				`json:"price"`
	Type 				string
	Quantity 			float64
	Comm				float64	
	Side 				string				`json:"type"`
	PositionSide 		string
	status 				int
}

// Status 订单状态
func (o Order) Status() int {
	return o.status
}

// SetStatus 设置订单状态
func (o *Order) SetStatus(s int) {
	o.status = s
}

// IsBuy 订单是否是买
func (o Order) IsBuy() bool {
	if o.Side != Buy {
		return false
	}
	return true
}


// IsSell 订单是否是卖
func (o Order) IsSell() bool {
	if o.Side != Sell {
		return false
	}
	return true
}