package spi

import (
	"context"
)

type Database interface {
	// GetServerInfo gets ServerInfo
	GetServerInfo() (*ServerInfo, error)

	// Explain retrieves execution plan of the given SQL statement.
	Explain(sqlText string) (string, error)

	// Exec executes SQL statements that does not return result
	// like 'ALTER', 'CREATE TABLE', 'DROP TABLE', ...
	Exec(sqlText string, params ...any) error

	// ExecContext executes SQL statements that does not return result
	// like 'ALTER', 'CREATE TABLE', 'DROP TABLE', ...
	ExecContext(ctx context.Context, sqlText string, params ...any) error

	// Query executes SQL statements that are expected multipe rows as result.
	// Commonly used to execute 'SELECT * FROM <TABLE>'
	//
	// Rows returned by Query() must be closed to prevent leaking resources.
	//
	//	rows, err := client.Query("select * from my_table where name = ?", "my_name")
	//	if err != nil {
	//		panic(err)
	//	}
	//	defer rows.Close()
	Query(sqlText string, params ...any) (Rows, error)

	// Query executes SQL statements that are expected multipe rows as result.
	// Commonly used to execute 'SELECT * FROM <TABLE>'
	//
	// Rows returned by QueryContext() must be closed to prevent server-side-resource leaks.
	//
	//	ctx, cancelFunc := context.WithTimeout(5*time.Second)
	//	defer cancelFunc()
	//
	//	rows, err := client.QueryContext(ctx, "select * from my_table where name = ?", my_name)
	//	if err != nil {
	//		panic(err)
	//	}
	//	defer rows.Close()
	QueryContext(ctx context.Context, sqlText string, params ...any) (Rows, error)

	// QueryRow executes a SQL statement that expects a single row result.
	//
	//	var cnt int
	//	row := client.QueryRow("select count(*) from my_table where name = ?", "my_name")
	//	row.Scan(&cnt)
	QueryRow(sqlText string, params ...any) Row

	// QueryRowContext executes a SQL statement that expects a single row result.
	//
	//	ctx, cancelFunc := context.WithTimeout(5*time.Second)
	//	defer cancelFunc()
	//
	//	var cnt int
	//	row := client.QueryRowContext(ctx, "select count(*) from my_table where name = ?", "my_name")
	//	row.Scan(&cnt)
	QueryRowContext(ctx context.Context, sqlText string, params ...any) Row

	// Appender creates a new Appender for the given table.
	// Appender should be closed as soon as finshing work, otherwise it may cause server side resource leak.
	//
	//	app, _ := client.Appender("MYTABLE")
	//	defer app.Close()
	//	app.Append("name", time.Now(), 3.14)
	Appender(tableName string) (Appender, error)
}

type ServerInfo struct {
	Version Version
	Runtime Runtime
}

type Version struct {
	Major          int32
	Minor          int32
	Patch          int32
	GitSHA         string
	BuildTimestamp string
	BuildCompiler  string
	Engine         string
}

type Runtime struct {
	OS             string
	Arch           string
	Pid            int32
	UptimeInSecond int64
	Processes      int32
	Goroutines     int32
	MemSys         uint64
	MemHeapSys     uint64
	MemHeapAlloc   uint64
	MemHeapInUse   uint64
	MemStackSys    uint64
	MemStackInUse  uint64
}

type Rows interface {
	// Next returns true if there are at least one more fetchable record remained.
	//
	//  rows, _ := db.Query("select name, value from my_table")
	//	for rows.Next(){
	//		var name string
	//		var value float64
	//		rows.Scan(&name, &value)
	//	}
	Next() bool

	// Scan retrieve values of columns in a row
	//
	//	for rows.Next(){
	//		var name string
	//		var value float64
	//		rows.Scan(&name, &value)
	//	}
	Scan(cols ...any) error

	// Close release all resources that assigned to the Rows
	Close() error

	// IsFetchable returns true if statement that produced this Rows was fetch-able (e.g was select?)
	IsFetchable() bool

	Message() string

	// Columns returns list of column info that consists of result of query statement.
	Columns() (Columns, error)
}

type Row interface {
	Success() bool
	Err() error
	Scan(cols ...any) error
	RowsAffected() int64
}

type Columns []*Column

type Column struct {
	Name   string
	Type   string
	Size   int32
	Length int32
}

type Appender interface {
	Append(cols ...any) error
	Close() error
}