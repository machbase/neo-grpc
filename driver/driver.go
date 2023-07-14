package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/machbase/neo-grpc/machrpc"
	spi "github.com/machbase/neo-spi"
)

func init() {
	sql.Register("machbase", &NeoDriver{})
}

const Name = "machbase"

var configReigstry = map[string]*DataSource{}

func RegisterDataSource(name string, conf *DataSource) {
	configReigstry[name] = conf
}

type DataSource struct {
	ServerAddr string
	ServerCert string
}

type NeoDriver struct {
}

var _ driver.Driver = &NeoDriver{}
var _ driver.DriverContext = &NeoDriver{}

func parseDataSourceName(name string) (addr string, certPath string) {
	u, err := url.Parse(name)
	if err != nil {
		return name, ""
	}
	addr = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)
	vals := u.Query()
	if serverCerts, ok := vals["server-cert"]; ok && len(serverCerts) > 0 {
		serverCert := serverCerts[0]
		return addr, serverCert
	}
	return addr, ""
}

func makeClient(dsn string) spi.DatabaseClient {
	var conf *DataSource
	if c, ok := configReigstry[dsn]; ok {
		conf = c
	} else {
		addr, certPath := parseDataSourceName(dsn)
		conf = &DataSource{ServerAddr: addr, ServerCert: certPath}
	}
	opts := []machrpc.Option{
		machrpc.WithServer(conf.ServerAddr, conf.ServerCert),
		machrpc.WithQueryTimeout(0),
	}
	return machrpc.NewClient(opts...)
}

func (d *NeoDriver) Open(name string) (driver.Conn, error) {
	client := makeClient(name)
	err := client.Connect()
	if err != nil {
		return nil, err
	}

	conn := &NeoConn{
		name:   name,
		client: client,
	}
	return conn, nil
}

func (d *NeoDriver) OpenConnector(name string) (driver.Connector, error) {
	client := makeClient(name)
	err := client.Connect()
	if err != nil {
		return nil, err
	}
	conn := &NeoConnector{
		name:   name,
		driver: d,
		client: client,
	}
	return conn, nil
}

type NeoConnector struct {
	driver.Connector
	name   string
	driver *NeoDriver
	client spi.DatabaseClient
}

func (cn *NeoConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn := &NeoConn{
		name:   cn.name,
		client: cn.client,
	}
	return conn, nil
}

func (cn *NeoConnector) Driver() driver.Driver {
	return cn.driver
}

type NeoConn struct {
	driver.Conn
	driver.Pinger
	driver.ConnBeginTx
	driver.QueryerContext
	driver.ExecerContext
	driver.ConnPrepareContext

	name   string
	client spi.DatabaseClient
}

func (c *NeoConn) Close() error {
	if c.client != nil {
		c.client.Disconnect()
		c.client = nil
	}
	return nil
}

func (c *NeoConn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *NeoConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return &NeoTx{}, nil
}

func (c *NeoConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	vals := make([]any, len(args))
	for i := range args {
		vals[i] = args[i].Value
	}
	rows, err := c.client.QueryContext(ctx, query, vals...)
	if err != nil {
		return nil, err
	}
	return &NeoRows{rows: rows}, nil
}

func (c *NeoConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	vals := make([]any, len(args))
	for i := range args {
		vals[i] = args[i].Value
	}
	row := c.client.QueryRowContext(ctx, query, vals...)
	if row.Err() != nil {
		return nil, row.Err()
	}
	return &NeoResult{row: row}, nil
}

func (c *NeoConn) Prepare(query string) (driver.Stmt, error) {
	return c.PrepareContext(context.Background(), query)
}

func (c *NeoConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	stmt := &NeoStmt{
		ctx:     ctx,
		conn:    c,
		sqlText: query,
	}
	return stmt, nil
}

func (c *NeoConn) Ping(ctx context.Context) error {
	// return driver.ErrBadConn
	return nil
}

type NeoTx struct {
}

func (tx *NeoTx) Commit() error {
	return nil
}

func (tx *NeoTx) Rollback() error {
	return errors.New("Rollback method is not supported")
}

type NeoStmt struct {
	driver.Stmt
	driver.StmtExecContext
	driver.StmtQueryContext

	ctx     context.Context
	conn    *NeoConn
	sqlText string
}

func (stmt *NeoStmt) Close() error {
	return nil
}

func (stmt *NeoStmt) NumInput() int {
	return -1
}

func (stmt *NeoStmt) Exec(args []driver.Value) (driver.Result, error) {
	vals := make([]any, len(args))
	for i := range args {
		vals[i] = args[i]
	}
	row := stmt.conn.client.QueryRow(stmt.sqlText, vals...)
	if row.Err() != nil {
		return nil, row.Err()
	}
	return &NeoResult{row: row}, nil
}

func (stmt *NeoStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	vals := make([]any, len(args))
	for i := range args {
		vals[i] = args[i].Value
	}
	row := stmt.conn.client.QueryRowContext(ctx, stmt.sqlText, vals...)
	if row.Err() != nil {
		return nil, row.Err()
	}
	return &NeoResult{row: row}, nil
}

func (stmt *NeoStmt) Query(args []driver.Value) (driver.Rows, error) {
	vals := make([]any, len(args))
	for i := range args {
		vals[i] = args[i]
	}
	rows, err := stmt.conn.client.Query(stmt.sqlText, vals...)
	if err != nil {
		return nil, err
	}
	return &NeoRows{rows: rows}, nil
}

func (stmt *NeoStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	vals := make([]any, len(args))
	for i := range args {
		vals[i] = args[i]
	}
	rows, err := stmt.conn.client.QueryContext(ctx, stmt.sqlText, vals...)
	if err != nil {
		return nil, err
	}
	return &NeoRows{rows: rows}, nil
}

type NeoResult struct {
	row spi.Row
}

func (r *NeoResult) LastInsertId() (int64, error) {
	return 0, errors.New("LastInsertId is not implemented")
}

func (r *NeoResult) RowsAffected() (int64, error) {
	if r.row == nil {
		return 0, nil
	}
	return r.row.RowsAffected(), nil
}

type NeoRows struct {
	rows spi.Rows
	cols spi.Columns
}

func (r *NeoRows) Columns() []string {
	if r.cols == nil {
		r.cols, _ = r.rows.Columns()
	}
	c := make([]string, len(r.cols))
	for i := range r.cols {
		c[i] = r.cols[i].Name
	}
	return c
}

func (r *NeoRows) Close() error {
	if r.rows == nil {
		return nil
	}
	err := r.rows.Close()
	if err != nil {
		return err
	}
	r.rows = nil
	return nil
}

func (r *NeoRows) Next(dest []driver.Value) error {
	if !r.rows.Next() {
		return io.EOF
	}
	vals := make([]any, len(dest))
	for i := range dest {
		vals[i] = &dest[i]
	}
	return r.rows.Scan(vals...)
}
