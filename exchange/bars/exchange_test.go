package bars

import (
	"testing"
)

func TestNewExchange(t *testing.T) {
	bars := NewExchangeBars(
		WithProxy("socks5://127.0.0.1:7890"),
	)
	if err := bars.Init(); err != nil {
		t.Log(err)
		return
	}
	t.Logf("kline length %v", len(bars.k))
}