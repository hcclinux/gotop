package exchange

// Options .
type Options struct {
	Proxy 		string
	SecretKey 	string
	APIKey 		string
	Name 		string
}

// Option used by the Exchange
type Option func(*Options)

// WithExchangeName .
func WithExchangeName(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// WithProxy .
func WithProxy(p string) Option {
	return func(o *Options) {
		o.Proxy = p
	}
}

// SecretKey .
func SecretKey(s string) Option {
	return func(o *Options) {
		o.SecretKey = s
	}
}

// APIKey .
func APIKey(a string) Option {
	return func(o *Options) {
		o.APIKey = a
	}
}