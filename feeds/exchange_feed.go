package feeds

import (
	"time"

	"gotop/brokers"
	"gotop/utils"

	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/builder"
)

/*
ExchangeFeed 处理交易所数据
Compression: 压缩比，例如TimeFrame=Minutes,那么Compression=5表示数据是5分钟k线数据。
TimeFrame: 时间帧类型，Minutes、Days等类型
FromDate: 起始时期
ToDate: 结束日期
Symbol: 交易对
APIKey: 密钥
SecretKey: 密钥
Proxy: 代理 "socks5://127.0.0.1:1080"
*/
type ExchangeFeed struct {
	Compression		uint8
	TimeFrame		uint8
	FromDate		time.Time
	ToDate			time.Time
	Symbol			goex.CurrencyPair
	APIKey			string
	SecretKey		string
	Proxy			string
	DropNew 		bool
}

// New 初始化一些跟交易所的对接准备，然后返回Kline
// exName 交易所名字
// bName broker name
func (ex *ExchangeFeed) New(exName, bName string) (brokers.CCBroker, error) {
	broker := brokers.CCBroker{Name: bName, CurrencyPair: ex.Symbol, DropNew: ex.DropNew}
	apiBuilder := builder.NewAPIBuilder().HttpTimeout(30 * time.Second).HttpProxy(ex.Proxy)
	api := apiBuilder.APIKey(ex.APIKey).APISecretkey(ex.SecretKey).Build(exName)
	broker.SetExchange(api)
	broker.KLinePeriod = ex.getPeriod()
	timePeriod := ex.handlePeriodTime(ex.getPeriod())
	broker.SetTimePeriod(timePeriod)
	kline, err := ex.handleKline(api)
	if err != nil {
		return broker, nil
	}
	broker.SetKLine(kline)
	return broker, nil
}

func (ex ExchangeFeed) handlePeriodTime(period int) int {
	switch period {
	case goex.KLINE_PERIOD_1MIN:
		return 1
	case goex.KLINE_PERIOD_3MIN:
		return 3
	case goex.KLINE_PERIOD_5MIN:
		return 5
	case goex.KLINE_PERIOD_15MIN:
		return 15
	case goex.KLINE_PERIOD_30MIN:
		return 30
	case goex.KLINE_PERIOD_60MIN:
		return 60
	case goex.KLINE_PERIOD_1H:
		return 60
	case goex.KLINE_PERIOD_2H:
		return 60 * 2
	case goex.KLINE_PERIOD_4H:
		return 60 * 4
	case goex.KLINE_PERIOD_6H:
		return 60 * 6
	case goex.KLINE_PERIOD_8H:
		return 60 * 8
	case goex.KLINE_PERIOD_12H:
		return 60 * 12
	case goex.KLINE_PERIOD_1DAY:
		return 60 * 24
	case goex.KLINE_PERIOD_3DAY:
		return 60 * 24 * 3
	case goex.KLINE_PERIOD_1WEEK:
		return 60 * 24 * 7
	case goex.KLINE_PERIOD_1MONTH:
		return 60 * 24 * 30
	case goex.KLINE_PERIOD_1YEAR:
		return 60 * 24 * 356
	}
	return 60
}


// HandleKline 处理K线数据
func (ex *ExchangeFeed) handleKline(api goex.API) (*utils.KLine, error) {
	var kline utils.KLine
	period := ex.getPeriod()
	lastTime := ex.handleTime()
	for {
		respData, err := api.GetKlineRecords(ex.Symbol, period, 1000, int(lastTime.Unix()*1000))
		if err != nil {
			return &utils.KLine{}, err
		}
		if len(respData) == 0 {
			break
		}
		if lastTime.Equal(time.Unix(respData[len(respData)-1].Timestamp, 0)) {
			break
		}
		lastTime = time.Unix(respData[len(respData)-1].Timestamp, 0)
		for _, k := range respData {
			c := utils.Candle{
				Date: time.Unix(k.Timestamp, 0),
				Open: k.Open,
				Close: k.Close,
				High: k.High,
				Low: k.Low,
				Volume: k.Vol,
			}
			kline.OHLC = append(kline.OHLC, c)
		}
	}
	return &kline, nil
}

func (ex ExchangeFeed) handleTime() time.Time {
	now := time.Now()
	timestamp := now.Unix() - int64(now.Second()) - int64((60 * now.Minute()))
	timestamp -= (3600 * 1000)
	if !ex.FromDate.IsZero() && ex.ToDate.IsZero() {
		return ex.FromDate
	}
	return time.Unix(timestamp, 0)
}


func (ex ExchangeFeed) getPeriod() int {
	switch ex.TimeFrame {
	case Minute:
		switch ex.Compression {
		case 1:
			return goex.KLINE_PERIOD_1MIN
		case 3:
			return goex.KLINE_PERIOD_3MIN
		case 5:
			return goex.KLINE_PERIOD_5MIN
		case 15:
			return goex.KLINE_PERIOD_15MIN
		case 30:
			return goex.KLINE_PERIOD_30MIN
		case 60:
			return goex.KLINE_PERIOD_60MIN
		}
	case Hour:
		switch ex.Compression {
		case 1:
			return goex.KLINE_PERIOD_1H
		case 2:
			return goex.KLINE_PERIOD_2H
		case 4:
			return goex.KLINE_PERIOD_4H
		case 6:
			return goex.KLINE_PERIOD_6H
		case 8:
			return goex.KLINE_PERIOD_8H
		case 12:
			return goex.KLINE_PERIOD_12H
		}
	case Day:
		switch ex.Compression {
		case 1:
			return goex.KLINE_PERIOD_1DAY
		case 3:
			return goex.KLINE_PERIOD_3DAY
		}
	case Week:
		return goex.KLINE_PERIOD_1WEEK
	case Month:
		return goex.KLINE_PERIOD_1MONTH
	case Year:
		return goex.KLINE_PERIOD_1YEAR
	}
	return 1
}