package bars

const (
	// Minute type time frame.
	Minute = iota + 1
	// Hour type time frame.
	Hour
	// Day type time frame.
	Day
	// Week type time frame.
	Week
	// Month type time frame.
	Month
	// Year type time frame.
	Year
)

// Bars ...
type Bars interface {
	Init(... Option) error
}

// Bar ohlcvt data
type Bar struct {
	Open 		float64
	High 		float64
	Low 		float64
	Close 		float64
	Volume 		float64
	Timestamp	int64
}
