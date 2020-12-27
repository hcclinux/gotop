package order

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

// Order 订单详情
type Order struct {
	Symbol				string
	DateTime 			time.Time   		`json:"time"`
	OrderID 			string
	ExecPrice			float64				`json:"price"`
	Type 				string
	Quantity 			float64
	Comm				float64	
	Side 				string				`json:"type"`
	PositionSide 		string
	status 				int
}