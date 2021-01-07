package exchange

import (
	"github.com/hcclinux/gotop/exchange/bars"
)

// Exchange .
type Exchange interface {
	b 		bars.Bars
}