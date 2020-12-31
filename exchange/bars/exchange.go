package bars

import (
	"time"

	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/builder"
)


// ExchangeBars ...
type ExchangeBars struct {
	// For example, timeFrame = minutes, then compression = 5 means that the data is 5-minute bar data
	compression		uint8
	// Time frame type, Minutes, Days
	timeFrame		uint8
	// Start time
	from			time.Time
	// End time
	to				time.Time
	// Currency pair
	symbol			goex.CurrencyPair
	apiKey			string
	secretKey		string
	// For example: "socks5://127.0.0.1:1080"
	proxy			string
	// Previous bar
	dropNew 		bool
}

// NewExchangeBars init exchange interface.
func NewExchangeBars(exName, bName string) error {
	apiBuilder := builder.NewAPIBuilder().HttpTimeout(30 * time.Second).HttpProxy(ex.proxy)
	api := apiBuilder.APIKey(ex.apiKey).APISecretkey(ex.secretKey).Build(exName)

	return nil
}

func (ex ExchangeBars) handlePeriodTime(period int) int {
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
func (ex *ExchangeBars) handleKline(api goex.API) (*utils.KLine, error) {
	var kline utils.KLine
	period := ex.getPeriod()
	lastTime := ex.handleTime()
	for {
		respData, err := api.GetKlineRecords(ex.symbol, period, 1000, int(lastTime.Unix()*1000))
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

func (ex ExchangeBars) handleTime() time.Time {
	now := time.Now()
	timestamp := now.Unix() - int64(now.Second()) - int64((60 * now.Minute()))
	timestamp -= (3600 * 1000)
	if !ex.from.IsZero() && ex.to.IsZero() {
		return ex.from
	}
	return time.Unix(timestamp, 0)
}


func (ex ExchangeBars) getPeriod() int {
	switch ex.timeFrame {
	case Minute:
		switch ex.compression {
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
		switch ex.compression {
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
		switch ex.compression {
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