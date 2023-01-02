package machrpc

import "time"

type ClientOption interface {
	clientoption()
}

type queryTimeoutOption struct {
	timeout time.Duration
}

func (o *queryTimeoutOption) clientoption() {}

type closeTimeoutOption struct {
	timeout time.Duration
}

func (o *closeTimeoutOption) clientoption() {}

func CloseTimeout(timeout time.Duration) ClientOption {
	return &closeTimeoutOption{timeout: timeout}
}

func QueryTimeout(timeout time.Duration) ClientOption {
	return &queryTimeoutOption{timeout: timeout}
}
