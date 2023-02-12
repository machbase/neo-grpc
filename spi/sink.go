package spi

type Sink interface {
	Write(p []byte) (n int, err error)
	Flush() error
	Close() error
}
