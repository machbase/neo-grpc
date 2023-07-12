package machrpc

import "time"

type Option func(*Client)

func WithServer(addr string) Option {
	return func(c *Client) {
		c.serverAddr = addr
	}
}

func WithServerCert(path string) Option {
	return func(c *Client) {
		c.serverCert = path
	}
}

func WithQueryTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.queryTimeout = timeout
	}
}

func WithCloseTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.closeTimeout = timeout
	}
}

func WithAppendTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.appendTimeout = timeout
	}
}
