package bars

import (
	"time"

	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/builder"
)


// ExchangeBars ...
type ExchangeBars struct {
	opts 	Options
	api 	*builder.APIBuilder
}

// NewExchangeBars init exchange interface.
func NewExchangeBars(opt ...Option) *ExchangeBars {
	opts := NewOptions(opt...)
	api := builder.NewAPIBuilder().HttpTimeout(10 * time.Second).HttpProxy(opts.Proxy)
	// api := apiBuilder.APIKey(ex.apiKey).APISecretkey(ex.secretKey).Build(exName)

	return &ExchangeBars{
		opts: opts,
		api: api,
	}
}

func (ex *ExchangeBars) handlePeriodTime(period int) int {
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
func (ex *ExchangeBars) handleKline(api goex.API) ([]*Bar, error) {
	kline := make([]*Bar, 0)
	period := ex.getPeriod()
	lastTime := ex.handleTime()
	for {
		resp, err := api.GetKlineRecords(ex.opts.Symbol, period, 1000, int(lastTime.Unix()*1000))
		if err != nil {
			return kline, err
		}
		if len(resp) == 0 {
			break
		}
		if lastTime.Equal(time.Unix(resp[len(resp)-1].Timestamp, 0)) {
			break
		}
		lastTime = time.Unix(resp[len(resp)-1].Timestamp, 0)
		for _, k := range resp {
			b := &Bar{
				Timestamp: k.Timestamp,
				Open: k.Open,
				Close: k.Close,
				High: k.High,
				Low: k.Low,
				Volume: k.Vol,
			}
			kline = append(kline, b)
		}
	}
	return kline, nil
}

func (ex *ExchangeBars) handleTime() time.Time {
	now := time.Now()
	timestamp := now.Unix() - int64(now.Second()) - int64((60 * now.Minute()))
	timestamp -= (3600 * 1000)
	if !ex.opts.From.IsZero() && ex.opts.To.IsZero() {
		return ex.opts.From
	}
	return time.Unix(timestamp, 0)
}


func (ex *ExchangeBars) getPeriod() int {
	switch ex.opts.TimeFrame {
	case Minute:
		switch ex.opts.Compression {
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
		switch ex.opts.Compression {
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
		switch ex.opts.Compression {
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