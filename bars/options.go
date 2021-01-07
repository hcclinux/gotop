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
	// For example: "socks5://127.0.0.1:1080"
	Proxy			string
	// Previous bar
	DropNew 		bool
	Path 			string
	ExchangeName 	string
}

// Option used by the ExchangeBars and CSVBars
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

// WithCompression select compression
func WithCompression(c uint8) Option {
	return func(o *Options) {
		o.Compression = c
	}
}

// WithTimeFrame select time frame
func WithTimeFrame(tf uint8) Option {
	return func(o *Options) {
		o.TimeFrame = tf
	}
}

// WithFrom start time of bar
func WithFrom(t time.Time) Option {
	return func(o *Options) {
		o.From = t
	}
}

// WithTo end time of bar
func WithTo(t time.Time) Option {
	return func(o *Options) {
		o.To = t
	}
}

// WithSymbol currency pair
func WithSymbol(s goex.CurrencyPair) Option {
	return func(o *Options) {
		o.Symbol = s
	}
}

// WithProxy proxy address and port
func WithProxy(s string) Option {
	return func(o *Options) {
		o.Proxy = s
	}
}

// WithPath file path
func WithPath(s string) Option {
	return func(o *Options) {
		o.Path = s
	}
}

// WithName exchange name
func WithName(name string) Option {
	return func(o *Options) {
		o.ExchangeName = name
	}
}