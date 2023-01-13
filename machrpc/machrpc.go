package machrpc

import (
	context "context"
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Client struct {
	conn grpc.ClientConnInterface
	cli  MachbaseClient

	closeTimeout  time.Duration
	queryTimeout  time.Duration
	appendTimeout time.Duration
}

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
		}
	}
	return client
}

func (client *Client) Connect(serverAddr string) error {
	conn, err := MakeGrpcConn(serverAddr)
	if err != nil {
		return errors.Wrap(err, "NewClient")
	}
	client.conn = conn
	client.cli = NewMachbaseClient(conn)
	return nil
}

func (client *Client) Disconnect() {
	client.conn = nil
	client.cli = nil
}

func (client *Client) queryContext() (context.Context, context.CancelFunc) {
	if client.queryTimeout > 0 {
		return context.WithTimeout(context.Background(), client.queryTimeout)
	} else {
		ctx := context.Background()
		return ctx, func() {}
	}
}

func (client *Client) Exec(sqlText string, params ...any) error {
	ctx, cancelFunc := client.queryContext()
	defer cancelFunc()
	return client.ExecContext(ctx, sqlText, params...)
}

func (client *Client) ExecContext(ctx context.Context, sqlText string, params ...any) error {
	pbparams, err := ConvertAnyToPb(params)
	if err != nil {
		return err
	}
	req := &ExecRequest{
		Sql:    sqlText,
		Params: pbparams,
	}
	rsp, err := client.cli.Exec(ctx, req)
	if err != nil {
		return err
	}
	if !rsp.Success {
		return fmt.Errorf(rsp.Reason)
	}
	return nil
}

func (client *Client) Query(sqlText string, params ...any) (*Rows, error) {
	ctx, cancelFunc := client.queryContext()
	defer cancelFunc()
	return client.QueryContext(ctx, sqlText, params...)
}

func (client *Client) QueryContext(ctx context.Context, sqlText string, params ...any) (*Rows, error) {
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
		return &Rows{client: client, handle: rsp.RowsHandle}, nil
	} else {
		if len(rsp.Reason) > 0 {
			return nil, errors.New(rsp.Reason)
		}
		return nil, errors.New("unknown error")
	}
}

type Rows struct {
	client *Client
	handle *RowsHandle
	values []any
	err    error
}

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

func (rows *Rows) Columns() ([]*Column, error) {
	ctx, cancelFunc := rows.client.queryContext()
	defer cancelFunc()

	rsp, err := rows.client.cli.Columns(ctx, rows.handle)
	if err != nil {
		return nil, err
	}
	if rsp.Success {
		return rsp.Columns, nil
	} else {
		if len(rsp.Reason) > 0 {
			return nil, errors.New(rsp.Reason)
		} else {
			return nil, fmt.Errorf("fail to get columns info")
		}
	}
}

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

func (rows *Rows) Scan(cols ...any) error {
	if rows.err != nil {
		return rows.err
	}
	if rows.values == nil {
		return sql.ErrNoRows
	}
	return scan(rows.values, cols)
}

func (client *Client) QueryRow(sqlText string, params ...any) *Row {
	ctx, cancelFunc := client.queryContext()
	defer cancelFunc()
	return client.QueryRowContext(ctx, sqlText, params...)
}

func (client *Client) QueryRowContext(ctx context.Context, sqlText string, params ...any) *Row {
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

func (client *Client) Appender(tableName string) (*Appender, error) {
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

	appendClient, err := client.cli.Append(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "AppendClient")
	}

	return &Appender{
		client:       client,
		appendClient: appendClient,
		tableName:    tableName,
		handle:       openRsp.Handle,
	}, nil
}

type Appender struct {
	client       *Client
	appendClient Machbase_AppendClient
	tableName    string
	handle       string
}

func (appender *Appender) Close() error {
	if appender.appendClient == nil {
		return nil
	}

	client := appender.appendClient
	appender.appendClient = nil

	err := client.CloseSend()
	if err != nil {
		return err
	}
	return nil
}

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
