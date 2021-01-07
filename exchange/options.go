package exchange

// Option used by the Exchange
type Option func(*Exchange)

// WithExchangeName .
func WithExchangeName(n string) Option {
	return func(o *Exchange) {
		o.name = n
	}
}

// WithProxy .
func WithProxy(p string) Option {
	return func(o *Exchange) {
		o.proxy = p
	}
}

// SecretKey .
func SecretKey(s string) Option {
	return func(o *Exchange) {
		o.secretKey = s
	}
}

// APIKey .
func APIKey(a string) Option {
	return func(o *Exchange) {
		o.apiKey = a
	}
}