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
	ExchangeName 	string
}

// Option used by the Bars
type Option func(*Options)


func newOptions(options ...Option) Options {
	// Default options
	opts := Options{
		Compression: 5,
		TimeFrame: Minute,
		From: time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local),
		Symbol: goex.BTC_USDT,
		DropNew: true,
		ExchangeName: goex.BINANCE,
	}

	for _, o := range options {
		o(&opts)
	}

	return opts
}

// WithCompression .
func WithCompression(c uint8) Option {
	return func(o *Options) {
		o.Compression = c
	}
}

// WithTimeFrame .
func WithTimeFrame(tf uint8) Option {
	return func(o *Options) {
		o.TimeFrame = tf
	}
}

// WithFrom .
func WithFrom(t time.Time) Option {
	return func(o *Options) {
		o.From = t
	}
}

// WithTo .
func WithTo(t time.Time) Option {
	return func(o *Options) {
		o.To = t
	}
}

// WithSymbol .
func WithSymbol(s goex.CurrencyPair) Option {
	return func(o *Options) {
		o.Symbol = s
	}
}

// WithProxy .
func WithProxy(s string) Option {
	return func(o *Options) {
		o.Proxy = s
	}
}

// WithPath .
func WithPath(s string) Option {
	return func(o *Options) {
		o.Path = s
	}
}

// WithName .
func WithName(name string) Option {
	return func(o *Options) {
		o.ExchangeName = name
	}
}