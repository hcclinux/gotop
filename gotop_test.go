package gotop

import (
    "fmt"
    "time"
    "testing"

    "github.com/nntaoli-project/goex"
    "github.com/nntaoli-project/goex/builder"
)

func TestGoex(t *testing.T) {
    apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second).HttpProxy("socks5://127.0.0.1:7890")
	api := apiBuilder.Build(goex.BINANCE)
	n,_ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")
	start := int(n.Unix()*1000)
    fmt.Println(api.GetKlineRecords(goex.BTC_USDT, goex.KLINE_PERIOD_30MIN, 500, start))
}