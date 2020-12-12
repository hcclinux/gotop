package analyzers

import (
	"fmt"

	"gotop/utils"
)

// Analyzers 分析师接口
type Analyzers interface {
	NotifyAnalyzer(utils.Order)
	PrintResult()
}

// TradeAnalyzer 交易分析功能
type TradeAnalyzer struct {
	count 			int
	buyCount		int
	sellCount 		int
	profitCount 	int
	lossCount 		int
	maxProfit		float64
	maxLoss			float64
	avgProfit		float64
	avgLoss			float64
	totalProfit		float64
	totalLoss		float64
	totalComm		float64
	buyOrder 		utils.Order
}

// NotifyAnalyzer 通知分析师订单
func (ta *TradeAnalyzer) NotifyAnalyzer(order utils.Order) {
	ta.Increase(order.Side)
	if order.Side == "buy" {
		ta.buyOrder = order
		return
	}
	ta.Calculate(order)
}

// PrintResult 打印结果
func (ta TradeAnalyzer) PrintResult() {
	str := "盈利次数:%d 亏损次数:%d 最大盈利:%.5f 最大亏损:%.5f 平均盈利:%.5f 平均亏损:%.5f 盈利金额:%.5f 亏损金额:%.5f 总手续费:%.5f 本次盈亏金额:%.5f\n"
	fmt.Printf(
		str,
		ta.profitCount,
		ta.lossCount,
		ta.maxProfit,
		ta.maxLoss,
		ta.avgProfit,
		ta.avgLoss,
		ta.totalProfit,
		ta.totalLoss,
		ta.totalComm,
		ta.totalProfit+ta.totalLoss,
	)
}

// Increase 增加count计数
func (ta *TradeAnalyzer) Increase(bos string) {
	ta.count++
	switch bos {
	case "buy":
		ta.buyCount++
	case "sell":
		ta.sellCount++
	}
}

// Calculate 计算盈亏
func (ta *TradeAnalyzer) Calculate(sellOrder utils.Order) {
	sellTotal := sellOrder.ExecPrice * sellOrder.Quantity
	buyTotal := ta.buyOrder.ExecPrice * sellOrder.Quantity
	ta.totalComm += (ta.buyOrder.Comm + sellOrder.Comm)
	result := sellTotal - buyTotal - (ta.buyOrder.Comm + sellOrder.Comm)
	if result < 0 {
		ta.plIncrease("loss")
		if result < ta.maxLoss {
			ta.maxLoss = result
		}
		ta.totalLoss += result
		ta.avgLoss = ta.totalLoss / float64(ta.lossCount)
	} else {
		ta.plIncrease("profit")
		if result > ta.maxProfit {
			ta.maxProfit = result
		}
		ta.totalProfit += result
		ta.avgProfit = ta.totalProfit / float64(ta.profitCount)
	}
}

func (ta *TradeAnalyzer) plIncrease(pol string) {
	switch pol {
	case "profit":
		ta.profitCount++
	case "loss":
		ta.lossCount++
	}
}