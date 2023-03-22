package sasl

import (
	"fmt"

	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type Option func(*Options)

func WithConnectionTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.ConnectionTimeout = timeout
	}
}

func WithPollInterval(pollIntervalSeconds float64) Option {
	return func(o *Options) {
		o.PollIntervalSeconds = pollIntervalSeconds
	}
}

func WithBatchSize(bs int64) Option {
	return func(o *Options) {
		o.BatchSize = bs
	}
}

func WithGSSAPISaslTransport(service string) Option {
	return func(o *Options) {
		o.SaslTransportConfig = map[string]string{
			"mechanismName": "GSSAPI",
			"service":       service,
		}
	}
}

type Options struct {
	PollIntervalSeconds float64
	BatchSize           int64
	ConnectionTimeout   time.Duration
	SaslTransportConfig map[string]string
}

var (
	DefaultOptions = Options{PollIntervalSeconds: 0.1, BatchSize: 10000, ConnectionTimeout: 10000 * time.Millisecond}
)

type Connection struct {
	Transport       *TSaslTransport
	ProtocolFactory *thrift.TBinaryProtocolFactory
}

func Connect(host string, port int, opts ...Option) (*Connection, error) {
	var options = DefaultOptions
	for _, opt := range opts {
		opt(&options)
	}

	socket, err := thrift.NewTSocketTimeout(fmt.Sprintf("%s:%d", host, port), options.ConnectionTimeout)
	if err != nil {
		return nil, err
	}

	transport, err := NewTSaslTransport(socket, host, options.SaslTransportConfig)
	if err != nil {
		return nil, err
	}
	// conf := &thrift.TConfiguration{
	// 	ConnectTimeout:     options.ConnectionTimeout,
	// 	SocketTimeout:      options.ConnectionTimeout,
	// 	MaxFrameSize:       1024 * 1024 * 256,
	// 	TBinaryStrictRead:  thrift.BoolPtr(true),
	// 	TBinaryStrictWrite: thrift.BoolPtr(true),
	// }
	// protocolFactory := thrift.NewTBinaryProtocolFactoryConf(conf)
	protocolFactory := thrift.NewTBinaryProtocolFactory(false, true)
	if err := transport.Open(); err != nil {
		return nil, err
	}
	return &Connection{transport, protocolFactory}, nil
}

func (c *Connection) isOpen() bool {
	return c.Transport != nil
}

func (c *Connection) Close() error {
	if c.isOpen() {
		c.Transport.Close()
	}
	return nil
}
