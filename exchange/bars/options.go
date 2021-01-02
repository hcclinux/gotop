package bars

import (
	"time"

	"github.com/nntaoli-project/goex"
)
// Options .
type Options struct {
	// Compression For example, TimeFrame = Minute, then Compression = 5 means that the data is 5-minute bar data
	Compression		uint8
	// TimeFrame type, Minute, Day
	TimeFrame		uint8
	// From start time
	From			time.Time
	// To end time
	To				time.Time
	// Symbol currency pair
	Symbol			goex.CurrencyPair
	APIKey			string
	SecretKey		string
	// For example: "socks5://127.0.0.1:1080"
	Proxy			string
	// Previous bar
	DropNew 		bool
	Path 			string
}

// NewOptions ...
func NewOptions(options ...Option) Options {
	opts := Options{
		Compression: 5,
		TimeFrame: Minute,
		From: time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local),
		To: time.Now(),
		Symbol: goex.BTC_USDT,
		Proxy: "socks5://127.0.0.1:1080",
		DropNew: true,
		Path: "./bars.csv",
	}

	for _, o := range options {
		o(&opts)
	}

	return opts
}