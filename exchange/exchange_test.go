package exchange

import (
	"testing"

	"github.com/nntaoli-project/goex"
)

func TestNewExchange(t *testing.T) {
	ex := NewExchange(
		WithExchangeName(goex.BINANCE),
		WithProxy("socks5://127.0.0.1:1080"),
		APIKey(""),
		SecretKey(""),
	)
	acc, err := ex.Balance()
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(acc)
}