package bars

import (
	"testing"
)

func TestNewExchange(t *testing.T) {
	bars := NewExchangeBars(
		WithProxy("127.0.0.1:7890"),
	)
	bars.Init()
}