package exchange

import (
	"testing"

	"github.com/nntaoli-project/goex"
)

func TestNewExchange(t *testing.T) {
	ex := NewExchange(
		WithExchangeName(goex.BINANCE),
		WithProxy("socks5://127.0.0.1:7890"),
		APIKey(""),
		SecretKey(""),
	)
	if err := ex.Init(); err != nil {
		t.Error(err)
		return
	}
	t.Log(ex.GetBalance())
}