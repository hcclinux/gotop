package exchange

import (
	"time"

	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/builder"
)

// Exchange .
type Exchange struct {
	api 		goex.API
	opts 		Options
}

// NewExchange .
func NewExchange(options ...Option) (ex *Exchange) {
	ex = &Exchange{
		opts: Options{},
	}

	for _, o := range options {
		o(&ex.opts)
	}
	return
}

// Init .
func (ex *Exchange) Init(opts ...Option) error {
	for _, o := range opts {
		o(&ex.opts)
	}

	build := builder.NewAPIBuilder().HttpTimeout(10 * time.Second).HttpProxy(ex.opts.Proxy)
	ex.api = build.APIKey(ex.opts.APIKey).APISecretkey(ex.opts.SecretKey).Build(ex.opts.Name)

	return nil
}

// GetBalance .
func (ex *Exchange) GetBalance() (*goex.Account, error) {
	return ex.api.GetAccount()
}