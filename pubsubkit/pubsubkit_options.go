package pubsubkit

// Options ...
type Options struct {
	checkExists bool
	autoCreate  bool
}

func newDefaultOptions() *Options {
	return &Options{
		checkExists: true,
		autoCreate:  false,
	}
}

// Option sets options for connect pubsub.
type Option func(*Options)

// WithAutoCreate will create pubsub resource when it's not exists yet.
func WithAutoCreate() Option {
	return func(o *Options) {
		o.autoCreate = true
	}
}

// WithoutCheckExistance bypass pubsub resource existence validations.
func WithoutCheckExistance() Option {
	return func(o *Options) {
		o.checkExists = false
	}
}
