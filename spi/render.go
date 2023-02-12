package spi

import (
	"context"
	"time"
)

type RowsRendererContext struct {
	Sink         Sink
	Rownum       bool
	Heading      bool
	TimeLocation *time.Location
	TimeFormat   string
	Precision    int
	HeaderHeight int
	ColumnNames  []string
	ColumnTypes  []string
}

type RowsRenderer interface {
	OpenRender(ctx *RowsRendererContext) error
	CloseRender()
	RenderRow(values []any) error
	PageFlush(heading bool)
}

type SeriesData struct {
	Name   string
	Values []float64
	Labels []string
}

type SeriesRenderer interface {
	Render(ctx context.Context, sink Sink, data []*SeriesData) error
}
