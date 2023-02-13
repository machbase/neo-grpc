// package machrpc is a refrence implementation of client
// that interwork with machbase-neo server via gRPC.
package machrpc

import (
	context "context"
	"database/sql"
	"fmt"
	"net"
	"time"

	spi "github.com/machbase/neo-spi"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func init() {
	spi.RegisterFactory("grpc", func() (spi.Database, error) {
		return NewClient(), nil
	})
}

// Client is a convenient data type represents client side of machbase-neo.
//
//	client := machrpc.NewClient()
type Client struct {
	conn grpc.ClientConnInterface
	cli  MachbaseClient

	closeTimeout  time.Duration
	queryTimeout  time.Duration
	appendTimeout time.Duration
}

// NewClient creates new instance of Client.
func NewClient(options ...ClientOption) *Client {
	client := &Client{
		closeTimeout:  3 * time.Second,
		queryTimeout:  0,
		appendTimeout: 3 * time.Second,
	}
	for _, opt := range options {
		switch o := opt.(type) {
		case *queryTimeoutOption:
			client.queryTimeout = o.timeout
		case *closeTimeoutOption:
			client.closeTimeout = o.timeout
		case *appendTimeoutOption:
			client.appendTimeout = o.timeout
		}
	}
	return client
}

// Connect make a connection to the server with the given address.
//
// serverAddr can be tcp://ipaddr:port or unix://path.
// The path of unix domain socket can be absolute/releative path.
func (client *Client) Connect(serverAddr string) error {
	conn, err := MakeGrpcConn(serverAddr)
	if err != nil {
		return errors.Wrap(err, "NewClient")
	}
	client.conn = conn
	client.cli = NewMachbaseClient(conn)
	return nil
}

// Disconnect release connection to the server.
func (client *Client) Disconnect() {
	client.conn = nil
	client.cli = nil
}

// GetServerInfo invoke gRPC call to get ServerInfo
func (client *Client) GetServerInfo() (*spi.ServerInfo, error) {
	ctx, cancelFunc := client.queryContext()
	defer cancelFunc()
	req := &ServerInfoRequest{}
	rsp, err := client.cli.GetServerInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	if !rsp.Success {
		return nil, errors.New(rsp.Reason)
	}
	return toSpiServerInfo(rsp), nil
}

func toSpiServerInfo(info *ServerInfo) *spi.ServerInfo {
	return &spi.ServerInfo{
		Version: spi.Version{
			Major:          info.Version.Major,
			Minor:          info.Version.Minor,
			Patch:          info.Version.Patch,
			GitSHA:         info.Version.GitSHA,
			BuildTimestamp: info.Version.BuildTimestamp,
			BuildCompiler:  info.Version.BuildCompiler,
			Engine:         info.Version.Engine,
		},
		Runtime: spi.Runtime{
			OS:             info.Runtime.OS,
			Arch:           info.Runtime.Arch,
			Pid:            info.Runtime.Pid,
			UptimeInSecond: info.Runtime.UptimeInSecond,
			Processes:      info.Runtime.Processes,
			Goroutines:     info.Runtime.Goroutines,
			MemSys:         info.Runtime.MemSys,
			MemHeapSys:     info.Runtime.MemHeapSys,
			MemHeapAlloc:   info.Runtime.MemHeapAlloc,
			MemHeapInUse:   info.Runtime.MemHeapInUse,
			MemStackSys:    info.Runtime.MemStackSys,
			MemStackInUse:  info.Runtime.MemStackInUse,
		},
	}
}

// Explain retrieve execution plan of the given SQL statement.
func (client *Client) Explain(sqlText string) (string, error) {
	ctx, cancelFunc := client.queryContext()
	defer cancelFunc()
	req := &ExplainRequest{Sql: sqlText}
	rsp, err := client.cli.Explain(ctx, req)
	if err != nil {
		return "", err
	}
	if !rsp.Success {
		return "", fmt.Errorf(rsp.Reason)
	}
	return rsp.Plan, nil
}

func (client *Client) queryContext() (context.Context, context.CancelFunc) {
	if client.queryTimeout > 0 {
		return context.WithTimeout(context.Background(), client.queryTimeout)
	} else {
		ctx := context.Background()
		return ctx, func() {}
	}
}

// Exec executes SQL statements that does not return result
// like 'ALTER', 'CREATE TABLE', 'DROP TABLE', ...
func (client *Client) Exec(sqlText string, params ...any) spi.Result {
	ctx, cancelFunc := client.queryContext()
	defer cancelFunc()
	return client.ExecContext(ctx, sqlText, params...)
}

// Exec executes SQL statements that does not return result
// like 'ALTER', 'CREATE TABLE', 'DROP TABLE', ...
func (client *Client) ExecContext(ctx context.Context, sqlText string, params ...any) spi.Result {
	pbparams, err := ConvertAnyToPb(params)
	if err != nil {
		return &ExecResult{err: err}
	}
	req := &ExecRequest{
		Sql:    sqlText,
		Params: pbparams,
	}
	rsp, err := client.cli.Exec(ctx, req)
	if err != nil {
		return &ExecResult{err: err}
	}
	if !rsp.Success {
		return &ExecResult{err: errors.New(rsp.Reason), message: rsp.Reason}
	}
	return &ExecResult{message: rsp.Reason, rowsAffected: rsp.RowsAffected}
}

type ExecResult struct {
	err          error
	rowsAffected int64
	message      string
}

func (r *ExecResult) Err() error {
	return r.err
}

func (r *ExecResult) RowsAffected() int64 {
	return r.rowsAffected
}

func (r *ExecResult) Message() string {
	return r.message
}

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
func (client *Client) Query(sqlText string, params ...any) (spi.Rows, error) {
	ctx, cancelFunc := client.queryContext()
	defer cancelFunc()
	return client.QueryContext(ctx, sqlText, params...)
}

// Query executes SQL statements that are expected multipe rows as result.
// Commonly used to execute 'SELECT * FROM <TABLE>'
//
// Rows returned by QueryContext() must be closed to prevent leaking resources.
//
//	ctx, cancelFunc := context.WithTimeout(5*time.Second)
//	defer cancelFunc()
//
//	rows, err := client.QueryContext(ctx, "select * from my_table where name = ?", my_name)
//	if err != nil {
//		panic(err)
//	}
//	defer rows.Close()
func (client *Client) QueryContext(ctx context.Context, sqlText string, params ...any) (spi.Rows, error) {
	pbparams, err := ConvertAnyToPb(params)
	if err != nil {
		return nil, err
	}

	req := &QueryRequest{Sql: sqlText, Params: pbparams}
	rsp, err := client.cli.Query(ctx, req)
	if err != nil {
		return nil, err
	}

	if rsp.Success {
		return &Rows{client: client, rowsAffected: rsp.RowsAffected, message: rsp.Reason, handle: rsp.RowsHandle}, nil
	} else {
		if len(rsp.Reason) > 0 {
			return nil, errors.New(rsp.Reason)
		}
		return nil, errors.New("unknown error")
	}
}

type Rows struct {
	client       *Client
	message      string
	rowsAffected int64
	handle       *RowsHandle
	values       []any
	err          error
}

// Close release all resources that assigned to the Rows
func (rows *Rows) Close() error {
	var ctx context.Context
	if rows.client.closeTimeout > 0 {
		ctx0, cancelFunc := context.WithTimeout(context.Background(), rows.client.closeTimeout)
		defer cancelFunc()
		ctx = ctx0
	} else {
		ctx = context.Background()
	}
	_, err := rows.client.cli.RowsClose(ctx, rows.handle)
	return err
}

// IsFetchable returns true if statement that produced this Rows was fetch-able (e.g was select?)
func (rows *Rows) IsFetchable() bool {
	return rows.handle != nil
}

func (rows *Rows) Message() string {
	return rows.message
}

func (rows *Rows) RowsAffected() int64 {
	return rows.rowsAffected
}

// Columns returns list of column info that consists of result of query statement.
func (rows *Rows) Columns() (spi.Columns, error) {
	ctx, cancelFunc := rows.client.queryContext()
	defer cancelFunc()

	rsp, err := rows.client.cli.Columns(ctx, rows.handle)
	if err != nil {
		return nil, err
	}
	if rsp.Success {
		return toSpiColumns(rsp.Columns), nil
	} else {
		if len(rsp.Reason) > 0 {
			return nil, errors.New(rsp.Reason)
		} else {
			return nil, fmt.Errorf("fail to get columns info")
		}
	}
}

func toSpiColumns(cols []*Column) spi.Columns {
	rt := make([]*spi.Column, len(cols))
	for i, c := range cols {
		rt[i] = &spi.Column{
			Name:   c.Name,
			Type:   c.Type,
			Size:   int(c.Size),
			Length: int(c.Length),
		}
	}
	return rt
}

// Next returns true if there are at least one more record that can be fetchable
// rows, _ := client.Query("select name, value from my_table")
//
//	for rows.Next(){
//		var name string
//		var value float64
//		rows.Scan(&name, &value)
//	}
func (rows *Rows) Next() bool {
	if rows.err != nil {
		return false
	}
	ctx, cancelFunc := rows.client.queryContext()
	defer cancelFunc()
	rsp, err := rows.client.cli.RowsFetch(ctx, rows.handle)
	if err != nil {
		rows.err = err
		return false
	}
	if rsp.Success {
		if rsp.HasNoRows {
			return false
		}
		rows.values = ConvertPbToAny(rsp.Values)
	} else {
		if len(rsp.Reason) > 0 {
			rows.err = errors.New(rsp.Reason)
		}
		rows.values = nil
	}
	return !rsp.HasNoRows
}

// Scan retrieve values of columns
//
//	for rows.Next(){
//		var name string
//		var value float64
//		rows.Scan(&name, &value)
//	}
func (rows *Rows) Scan(cols ...any) error {
	if rows.err != nil {
		return rows.err
	}
	if rows.values == nil {
		return sql.ErrNoRows
	}
	return scan(rows.values, cols)
}

// QueryRow executes a SQL statement that expects a single row result.
//
//	var cnt int
//	row := client.QueryRoq("select count(*) from my_table where name = ?", "my_name")
//	row.Scan(&cnt)
func (client *Client) QueryRow(sqlText string, params ...any) spi.Row {
	ctx, cancelFunc := client.queryContext()
	defer cancelFunc()
	return client.QueryRowContext(ctx, sqlText, params...)
}

func (client *Client) QueryRowContext(ctx context.Context, sqlText string, params ...any) spi.Row {
	pbparams, err := ConvertAnyToPb(params)
	if err != nil {
		return &Row{success: false, err: err}
	}

	req := &QueryRowRequest{Sql: sqlText, Params: pbparams}
	rsp, err := client.cli.QueryRow(ctx, req)
	if err != nil {
		return &Row{success: false, err: err}
	}

	var row = &Row{}
	row.success = rsp.Success
	row.rowsAffected = rsp.RowsAffected
	row.err = nil
	if !rsp.Success && len(rsp.Reason) > 0 {
		row.err = errors.New(rsp.Reason)
	}
	row.values = ConvertPbToAny(rsp.Values)
	return row
}

type Row struct {
	success bool
	err     error
	values  []any

	rowsAffected int64
}

func (row *Row) Success() bool {
	return row.success
}

func (row *Row) Err() error {
	return row.err
}

func (row *Row) Scan(cols ...any) error {
	if row.err != nil {
		return row.err
	}
	if !row.success {
		return sql.ErrNoRows
	}
	err := scan(row.values, cols)
	return err
}

func (row *Row) RowsAffected() int64 {
	return row.rowsAffected
}

func (row *Row) Message() string {
	// TODO
	return ""
}
func (row *Row) Values() []any {
	return row.values
}

func scan(src []any, dst []any) error {
	for i := range dst {
		if i >= len(src) {
			return fmt.Errorf("column %d is out of range %d", i, len(src))
		}
		if src[i] == nil {
			dst[i] = nil
			continue
		}
		var err error
		switch v := src[i].(type) {
		default:
			return fmt.Errorf("column %d is %T, not compatible with %T", i, v, dst[i])
		case *int:
			err = ScanInt32(int32(*v), dst[i])
		case *int16:
			err = ScanInt16(*v, dst[i])
		case *int32:
			err = ScanInt32(*v, dst[i])
		case *int64:
			err = ScanInt64(*v, dst[i])
		case *time.Time:
			err = ScanDateTime(*v, dst[i])
		case *float32:
			err = ScanFloat32(*v, dst[i])
		case *float64:
			err = ScanFloat64(*v, dst[i])
		case *net.IP:
			err = ScanIP(*v, dst[i])
		case *string:
			err = ScanString(*v, dst[i])
		case []byte:
			err = ScanBytes(v, dst[i])
		case int:
			err = ScanInt32(int32(v), dst[i])
		case int16:
			err = ScanInt16(v, dst[i])
		case int32:
			err = ScanInt32(v, dst[i])
		case int64:
			err = ScanInt64(v, dst[i])
		case time.Time:
			err = ScanDateTime(v, dst[i])
		case float32:
			err = ScanFloat32(v, dst[i])
		case float64:
			err = ScanFloat64(v, dst[i])
		case net.IP:
			err = ScanIP(v, dst[i])
		case string:
			err = ScanString(v, dst[i])
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Appender creates a new Appender for the given table.
// Appender should be closed otherwise it may cause server side resource leak.
//
//	app, _ := client.Appender("MYTABLE")
//	defer app.Close()
//	app.Append("name", time.Now(), 3.14)
func (client *Client) Appender(tableName string) (spi.Appender, error) {
	var ctx0 context.Context
	if client.appendTimeout > 0 {
		_ctx, _cf := context.WithTimeout(context.Background(), client.appendTimeout)
		defer _cf()
		ctx0 = _ctx
	} else {
		ctx0 = context.Background()
	}

	openRsp, err := client.cli.Appender(ctx0, &AppenderRequest{TableName: tableName})
	if err != nil {
		return nil, errors.Wrap(err, "Appender")
	}

	if !openRsp.Success {
		return nil, errors.New(openRsp.Reason)
	}

	appendClient, err := client.cli.Append(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "AppendClient")
	}

	return &Appender{
		client:       client,
		appendClient: appendClient,
		tableName:    openRsp.TableName,
		tableType:    spi.TableType(openRsp.TableType),
		handle:       openRsp.Handle,
	}, nil
}

type Appender struct {
	client       *Client
	appendClient Machbase_AppendClient
	tableName    string
	tableType    spi.TableType
	handle       string
}

// Close releases all resources that allocated to the Appender
func (appender *Appender) Close() (int64, int64, error) {
	if appender.appendClient == nil {
		return 0, 0, nil
	}

	client := appender.appendClient
	appender.appendClient = nil

	err := client.CloseSend()
	if err != nil {
		return 0, 0, err
	}
	return 0, 0, nil
}

func (appender *Appender) TableName() string {
	return appender.tableName
}

func (appender *Appender) TableType() spi.TableType {
	return appender.tableType
}

func (appender *Appender) AppendWithTimestamp(ts time.Time, cols ...any) error {
	return appender.Append(append([]any{ts}, cols...))
}

// Append appends a new record of the table.
func (appender *Appender) Append(cols ...any) error {
	if appender.appendClient == nil {
		return sql.ErrTxDone
	}

	params, err := ConvertAnyToPb(cols)
	if err != nil {
		return err
	}
	err = appender.appendClient.Send(&AppendData{
		Handle: appender.handle,
		Params: params,
	})
	return err
}

func (appender *Appender) Columns() (spi.Columns, error) {
	// TODO implements
	return nil, nil
}
