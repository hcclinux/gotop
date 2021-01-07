package exchange

import (
	"time"

	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/builder"
)

// Exchange .
type Exchange struct {
	api 		goex.API
	proxy 		string
	secretKey 	string
	apiKey 		string
	name 		string
}

// NewExchange .
func NewExchange(options ...Option) (ex *Exchange) {
	ex = &Exchange{}

	for _, o := range options {
		o(ex)
	}

	build := builder.NewAPIBuilder().HttpTimeout(10 * time.Second).HttpProxy(ex.proxy)
	ex.api = build.APIKey(ex.apiKey).APISecretkey(ex.secretKey).Build(ex.name)
	return
}

// Balance .
func (e *Exchange) Balance() (*goex.Account, error) {
	return e.api.GetAccount()
}