package brokers

import (
	"math/rand"
	"time"
	"errors"

	"gotop/analyzers"
	"gotop/utils"
)



// Broker 经纪人基础功能
type Broker interface {
	Next()										utils.Candle
	Buy(utils.Order)		(utils.Order, error)
	Sell(utils.Order)		(utils.Order, error)
	SetCommission(float64)
	GetCommission()								float64
	GetBalance()								(float64, float64, error)
	SetCash(float64)
	Cancel(string)								error
	GetName()									string
	SetKLine(*utils.KLine)
	NotifyOrder(utils.Order)
	AddAnalyzer(analyzers.Analyzers)
	Print()
	SetIndex()
	GetOrder()									map[string]utils.Order
	GetKLine()									[]utils.Candle
	TickerPrice()								(float64, error)
}

// DBroker 默认经纪人
type DBroker struct {
	Name 				string
	kline 				*utils.KLine
	commission 			float64
	index				int
	cash 				float64
	quantity 			float64
	orders 				map[string]utils.Order
	analyzers			[]analyzers.Analyzers
}

const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    
	letterIdxMask = 1 << letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)


func makeOrder(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// Init 初始化经纪人
func (db *DBroker) Init() {
	db.orders = make(map[string]utils.Order)
}

// GetKLine 获取所有k线数据
func (db DBroker) GetKLine() []utils.Candle {
	return db.kline.OHLC
}

// GetOrder 获取所有订单
func (db DBroker) GetOrder() map[string]utils.Order {
	return db.orders
}

// TickerPrice 返回最新价格
func (db DBroker) TickerPrice() (float64, error) {
	return 0, nil
}

// Print 打印结果
func (db DBroker) Print() {
	for _, a := range db.analyzers {
		a.PrintResult()
	}
}

// AddAnalyzer 添加分析师
func (db *DBroker) AddAnalyzer(anal analyzers.Analyzers) {
	db.analyzers = append(db.analyzers, anal)
}

// NotifyOrder 通知订单
func (db *DBroker) NotifyOrder(order utils.Order) {
	for _, a := range db.analyzers {
		a.NotifyAnalyzer(order)
	}
}


// Next 下一个K线
func (db DBroker) Next() utils.Candle {
	if db.idx() == len(db.kline.OHLC) {
		return utils.Candle{}
	}
	return db.kline.OHLC[db.idx()]
}

// SetIndex 自增
func (db *DBroker) SetIndex() {
	db.index++
}

func (db DBroker) idx() int {
	return db.index
}

// Buy 买入
func (db *DBroker) Buy(order utils.Order) (utils.Order, error) {
	if len(db.kline.OHLC) == db.idx() + 1 {
		order.Side = utils.Buy
		order.SetStatus(utils.Rejected)
		return order, errors.New("out of the k line range") 
	}
	open := db.kline.OHLC[db.idx()+1].Open
	cash, _, _ := db.GetBalance()
	if (cash / open) < order.Quantity {
		return order, errors.New("Exceeds the available amount")
	}
	order.Time = db.kline.OHLC[db.idx()+1].Date
	order.OrderID = makeOrder(32)
	order.Side = utils.Buy
	order.Comm = open * db.GetCommission()
	if order.Type == utils.Market {
		order.ExecPrice = open
	}
	db.AddQuantity(order.Quantity)
	order.SetStatus(utils.Completed)
	db.orders[order.OrderID] = order
	db.NotifyOrder(order)
    return order, nil
}

// Sell 卖出
func (db *DBroker) Sell(order utils.Order) (utils.Order, error) {
	if len(db.kline.OHLC) == db.idx() + 1 {
		order.Side = utils.Sell
		order.SetStatus(utils.Rejected)
		return order, errors.New("out of the k line range") 
	}
	_, quantity, _ := db.GetBalance()
	if order.Quantity > quantity {
		return order, errors.New("Excess holdings")
	}
	order.Time = db.kline.OHLC[db.idx()+1].Date
	order.OrderID = makeOrder(32)
	order.Side = utils.Sell
	order.Comm = db.kline.OHLC[db.idx()+1].Open * db.GetCommission()
	if order.Type == utils.Market {
		order.ExecPrice = db.kline.OHLC[db.idx()+1].Open
	}
	db.SubQuantity(order.Quantity)
	order.SetStatus(utils.Completed)
	db.orders[order.OrderID] = order
	db.NotifyOrder(order)
    return order, nil
}
// AddQuantity 增加数量
func (db *DBroker) AddQuantity(q float64) {
	db.quantity += q
}

// SubQuantity 减少数量
func (db *DBroker) SubQuantity(q float64) {
	db.quantity -= q
}

// SetCommission 设置手续费
func (db *DBroker) SetCommission(comm float64) {
    db.commission = comm
}

// GetCommission 获取手续费
func (db DBroker) GetCommission() float64 {
	return db.commission
}

// GetBalance 获取现金
func (db DBroker) GetBalance() (float64, float64, error) {
	return db.cash, db.quantity, nil
}

// SetCash 设置金额
func (db *DBroker) SetCash(cash float64) {
	db.cash = cash
}

// Cancel 取消订单
func (db *DBroker) Cancel(oid string) error {
	return nil
}

// GetName 获取经纪人名字
func (db DBroker) GetName() string {
    return db.Name
}

// SetKLine 设置K线
func (db *DBroker) SetKLine(k *utils.KLine) {
    db.kline = k
}