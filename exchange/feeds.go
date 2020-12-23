package exchange

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

// Exchange is an abstract read k-line data interface.
type Exchange interface {
	
}