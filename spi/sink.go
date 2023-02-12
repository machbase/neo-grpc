package spi

type Sink interface {
	Write([]byte) (int, error)
	Flush() error
	Close() error
}
