package brokers

import (
	"time"
	"strconv"
	"fmt"

	"gotop/utils"
	"gotop/analyzers"
	"github.com/nntaoli-project/goex"
)



// CCBroker 加密货币经纪人
type CCBroker struct {
	Name 				string
	KLinePeriod 		int 			
	CurrencyPair 		goex.CurrencyPair
	DropNew 			bool
	commission 			float64
	exchange 			goex.API
	kline 				*utils.KLine
	analyzers			[]analyzers.Analyzers
	index 				int
	cash 				float64
	value 				float64
	orders 				map[string]utils.Order
	timePeriod			int
}

// Init 初始化经纪人
func (ccb *CCBroker) Init() {
	ccb.orders = make(map[string]utils.Order)
}

// AddAnalyzer 添加分析师
func (ccb *CCBroker) AddAnalyzer(anal analyzers.Analyzers) {
	ccb.analyzers = append(ccb.analyzers, anal)
}

// NotifyOrder 通知订单
func (ccb *CCBroker) NotifyOrder(order utils.Order) {
	switch order.Status() {
	case utils.Completed:
		if order.IsBuy() {
			fmt.Println(fmt.Sprintf("买入价格:%.5f", order.ExecPrice))
		} else {
			fmt.Println(fmt.Sprintf("卖出价格:%.5f", order.ExecPrice))
		}
		for _, a := range ccb.analyzers {
			a.NotifyAnalyzer(order)
		}
	case utils.Rejected:
		fmt.Println("订单被拒绝")
	}
}


// SetCash 设置现金
func (ccb *CCBroker) SetCash(cash float64) {

}

// Print 打印结果
func (ccb CCBroker) Print() {
	for _, a := range ccb.analyzers {
		a.PrintResult()
	}
}

// GetKLine 获取所有k线数据
func (ccb CCBroker) GetKLine() []utils.Candle {
	return ccb.kline.OHLC
}

// GetOrder 获取所有订单
func (ccb CCBroker) GetOrder() map[string]utils.Order {
	return ccb.orders
}

// SetIndex 自增下标
func (ccb *CCBroker) SetIndex() {
	ccb.index++
}

func (ccb CCBroker) idx() int {
	return ccb.index
}

// SetKLine 设置K线
func (ccb *CCBroker) SetKLine(k *utils.KLine) {
	if ccb.DropNew {
		k.OHLC = k.OHLC[:len(k.OHLC)-1]
	}
	ccb.kline = k
}

// GetName 获取经纪人名字
func (ccb CCBroker) GetName() string {
	return ccb.Name
}

// GetBalance 获取交易对数量
func (ccb CCBroker) GetBalance() (float64, float64, error) {
	account, err := ccb.exchange.GetAccount()
	if err != nil {
		return 0, 0, err
	}
	cash := account.SubAccounts[ccb.CurrencyPair.CurrencyA].Amount
	val := account.SubAccounts[ccb.CurrencyPair.CurrencyB].Amount
	return cash, val, nil
}


// GetCommission 获取手续费
func (ccb CCBroker) GetCommission() float64 {
	return ccb.commission
}

// SetCommission 设置手续费
func (ccb *CCBroker) SetCommission(comm float64) {
	ccb.commission = comm
}


// Buy 买入
func (ccb CCBroker) Buy(o utils.Order) (utils.Order, error) {
	var respOrder *goex.Order
	var err error
	a := strconv.FormatFloat(o.Quantity,'f', -1, 64)
	p := strconv.FormatFloat(o.ExecPrice,'f', -1, 64)
	switch o.Type {
	case utils.Limit:
		respOrder, err = ccb.exchange.LimitBuy(a, p, ccb.CurrencyPair)
	case utils.Market:
		respOrder, err = ccb.exchange.MarketBuy(a, p, ccb.CurrencyPair)
	}
	if err != nil {
		return utils.Order{}, err
	}
	order := ccb.replacetOrder(respOrder)
	order.SetStatus(utils.Completed)
	return order, nil
}


// Sell 卖出
func (ccb CCBroker) Sell(o utils.Order) (utils.Order, error) {
	var respOrder *goex.Order
	var err error
	a := strconv.FormatFloat(o.Quantity,'f', -1, 64)
	p := strconv.FormatFloat(o.ExecPrice,'f', -1, 64)
	switch o.Type {
	case utils.Limit:
		respOrder, err = ccb.exchange.LimitSell(a, p, ccb.CurrencyPair)
	case utils.Market:
		respOrder, err = ccb.exchange.MarketSell(a, p, ccb.CurrencyPair)
	}
	if err != nil {
		return utils.Order{}, err
	}
	order := ccb.replacetOrder(respOrder)
	order.SetStatus(utils.Completed)
	return order, nil
}


// Next 下一个K线
func (ccb *CCBroker) Next() utils.Candle {
	if ccb.idx() < len(ccb.kline.OHLC) {
		return ccb.kline.OHLC[ccb.idx()]
	}
	return ccb.getData()
}

func (ccb *CCBroker) getData() utils.Candle {
	var c utils.Candle
	now := time.Now()
	dest := (now.Minute() + ccb.timePeriod) - ((now.Minute() + ccb.timePeriod) % ccb.timePeriod)
	next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), dest, 0, 0, now.Location())
	ticker := time.NewTicker(next.Sub(now))
	defer ticker.Stop()
	select {
	case <-ticker.C:
		time.Sleep(time.Second * 2)
		lastTime := ccb.kline.OHLC[len(ccb.kline.OHLC)-1].Date
		resp, err := ccb.exchange.GetKlineRecords(ccb.CurrencyPair, ccb.KLinePeriod, 100, int(lastTime.Unix()*1000))
		if err != nil {
			return utils.Candle{}
		}
		c = ccb.dropNew(resp, next.Unix())
		ccb.kline.OHLC = append(ccb.kline.OHLC, c)
	}
	return c
}

func (ccb CCBroker) dropNew(resp []goex.Kline, timestamp int64) utils.Candle {
	if ccb.DropNew {
		for _, v := range resp {
			if v.Timestamp == timestamp {
				return ccb.replaceCandle(v)
			}
		}
		return utils.Candle{}
	}
	if len(resp) == 0 {
		return utils.Candle{}
	}
	return ccb.replaceCandle(resp[len(resp)-1])
}

// Cancel 取消订单
func (ccb CCBroker) Cancel(order string) error {
	_, err := ccb.exchange.CancelOrder(order, ccb.CurrencyPair)
	if err != nil {
		return err
	}
	return nil
}

func (ccb CCBroker) replaceCandle(resp goex.Kline) utils.Candle {
	return utils.Candle{
		Date: time.Unix(resp.Timestamp, 0),
		High: resp.High,
		Low: resp.Low,
		Open: resp.Open,
		Close: resp.Close,
		Volume: resp.Vol,
	}
}

// ConvertOrder 转换订单类型
func (ccb CCBroker) replacetOrder(order *goex.Order) utils.Order {
	o := utils.Order{
		ExecPrice: order.Price,
		Quantity: order.Amount,
		Time: time.Unix(int64(order.OrderTime), 0),
		OrderID:order.OrderID2,
		Type: order.Type,
	}
	return o
}

// SetTimePeriod 设置时间帧
func (ccb *CCBroker) SetTimePeriod(tf int) {
	ccb.timePeriod = tf
}

// SetExchange 设置交易所
func (ccb *CCBroker) SetExchange(ex goex.API) {
	ccb.exchange = ex
}

// TickerPrice 返回最新价格
func (ccb CCBroker) TickerPrice() (float64, error) {
	price, err := ccb.exchange.GetTicker(ccb.CurrencyPair)
	if err != nil {
		return 0, err
	}
	return price.Last, err
}