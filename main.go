package main

import (
    "fmt"
    "time"

    "github.com/nntaoli-project/goex"
    "github.com/nntaoli-project/goex/builder"
)

func main() {
    apiBuilder := builder.NewAPIBuilder().HttpTimeout(5 * time.Second).HttpProxy("socks5://127.0.0.1:1080")
	api := apiBuilder.APIKey("").APISecretkey("").Build(goex.BINANCE)
	t,_ := time.Parse("2006-01-02 15:04:05", "2018-04-11 00:00:00")
	start := int(t.Unix()*1000)
    fmt.Println(api.GetKlineRecords(goex.BTC_USDT, goex.KLINE_PERIOD_30MIN, 500, start))
}