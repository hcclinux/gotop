package bars

import (
	"testing"
)

func TestNewCSV(t *testing.T) {
	bars := NewCSVBars(
		WithPath("btc_1m.csv"),
	)
	if err := bars.Init(); err != nil {
		t.Log(err)
		return
	}
	t.Logf("kline length %v", len(bars.k))
}