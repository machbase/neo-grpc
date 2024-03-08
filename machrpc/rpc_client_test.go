package machrpc_test

import (
	context "context"
	"testing"
	"time"

	"github.com/d5/tengo/v2/require"
	"github.com/machbase/neo-grpc/machrpc"
	"github.com/machbase/neo-grpc/mock"
	spi "github.com/machbase/neo-spi"
)

func TestMain(m *testing.M) {
	svr := &mock.MockServer{}

	err := svr.Start()
	if err != nil {
		panic(err)
	}

	m.Run()

	svr.Stop()
}

func newClient(t *testing.T) spi.DatabaseClient {
	t.Helper()
	opts := []machrpc.Option{
		machrpc.WithServer(mock.MockServerAddr),
		machrpc.WithQueryTimeout(3 * time.Second),
		machrpc.WithAppendTimeout(3 * time.Second),
	}
	cli, err := machrpc.NewClient(opts...)
	if err != nil {
		t.Fatalf("new client: %s", err.Error())
	}
	return cli
}

func newConn(t *testing.T) spi.Conn {
	t.Helper()
	cli := newClient(t)
	conn, err := cli.Connect(context.TODO(), machrpc.WithPassword("sys", "manager"))
	require.Nil(t, err)
	require.NotNil(t, conn)
	return conn
}

func TestAuth(t *testing.T) {
	opts := []machrpc.Option{
		machrpc.WithServer(mock.MockServerAddr),
	}
	cli, err := machrpc.NewClient(opts...)
	if err != nil {
		t.Fatalf("new client: %s", err.Error())
	}
	auth, ok := cli.(spi.DatabaseAuth)
	if !ok {
		t.Fatalf("client is not implments spi.DatabaseAuth")
	}

	ok, err = auth.UserAuth("sys", "mm")
	require.NotNil(t, err)
	require.Equal(t, "invalid username or password", err.Error())
	require.False(t, ok)

	ok, err = auth.UserAuth("sys", "manager")
	if err != nil {
		t.Fatalf("UserAuth failed: %s", err.Error())
	}
	require.True(t, ok)
}

func TestNewClient(t *testing.T) {
	var cli spi.DatabaseClient
	var err error

	// no server address
	cli, err = machrpc.NewClient()
	require.NotNil(t, err, "no error without server addr, want error")
	require.Nil(t, cli, "new client should fail")
	require.Equal(t, "server address is not specified", err.Error())

	// success creating client
	opts := []machrpc.Option{
		machrpc.WithServer(mock.MockServerAddr),
	}
	cli, err = machrpc.NewClient(opts...)
	if err != nil {
		t.Fatalf("new client: %s", err.Error())
	}

	ctx := context.TODO()

	// empty username, password
	conn, err := cli.Connect(ctx)
	require.NotNil(t, err)
	require.Equal(t, "invalid username or password", err.Error())
	require.Nil(t, conn)

	// wrong password
	conn, err = cli.Connect(ctx, machrpc.WithPassword("sys", "mm"))
	require.NotNil(t, err)
	require.Equal(t, "invalid username or password", err.Error())
	require.Nil(t, conn)

	// correct username, password
	conn, err = cli.Connect(ctx, machrpc.WithPassword("sys", "manager"))
	require.Nil(t, err)
	require.NotNil(t, conn)

	conn.Close()
}

func TestGetServerInfo(t *testing.T) {
	aux, ok := newClient(t).(spi.DatabaseAux)
	if !ok {
		t.Fatal("client is not implements spi.DatabaseAux")
	}
	nfo, err := aux.GetServerInfo()
	if err != nil {
		t.Fatalf("GetPostflights error: %s", err.Error())
	}
	require.Nil(t, err)
	require.NotNil(t, nfo)
}

func TestGetInflights(t *testing.T) {
	aux, ok := newClient(t).(spi.DatabaseAux)
	if !ok {
		t.Fatal("client is not implements spi.DatabaseAux")
	}
	postflights, err := aux.GetInflights()
	if err != nil {
		t.Fatalf("GetInflights error: %s", err.Error())
	}
	require.Nil(t, err)
	require.NotNil(t, postflights)
	require.Equal(t, 0, len(postflights))
}

func TestGetPostflights(t *testing.T) {
	aux, ok := newClient(t).(spi.DatabaseAux)
	if !ok {
		t.Fatal("client is not implements spi.DatabaseAux")
	}
	inflights, err := aux.GetPostflights()
	if err != nil {
		t.Fatalf("GetPostflights error: %s", err.Error())
	}
	require.Nil(t, err)
	require.NotNil(t, inflights)
	require.Equal(t, 0, len(inflights))
}

func TestGetServicePorts(t *testing.T) {
	aux, ok := newClient(t).(spi.DatabaseAux)
	if !ok {
		t.Fatal("client is not implements spi.DatabaseAux")
	}
	ports, err := aux.GetServicePorts("grpc")
	if err != nil {
		t.Fatalf("GetServicePorts error: %s", err.Error())
	}
	require.NotNil(t, ports)
	require.Equal(t, 0, len(ports))
}

func TestPing(t *testing.T) {
	conn := newConn(t)
	pinger, ok := conn.(spi.Pinger)
	if !ok {
		t.Fatalf("Connection is not implments spi.Pinger")
	}
	dur, err := pinger.Ping()
	require.Nil(t, err)
	if dur == 0 {
		t.Fatal("ping failure")
	}
	conn.Close()
}

func TestExplainTest(t *testing.T) {
	conn := newConn(t)
	defer conn.Close()
	exp, ok := conn.(spi.Explainer)
	if !ok {
		t.Fatal("client is not implements spi.Explainer")
	}
	result, err := exp.Explain(context.TODO(), "select * from dummy", true)
	if err != nil {
		t.Fatalf("Explain error: %s", err.Error())
	}
	require.Equal(t, "explain dummy result", result)
}

func TestExec(t *testing.T) {
	conn := newConn(t)
	defer conn.Close()
	result := conn.Exec(context.TODO(), "insert into example (name, time, value) values(?, ?, ?)", 1, 2, 3)
	require.NotNil(t, result)
	require.Nil(t, result.Err())
	require.Equal(t, int64(1), result.RowsAffected())
}

func TestQueryRow(t *testing.T) {
	conn := newConn(t)
	defer conn.Close()
	row := conn.QueryRow(context.TODO(), "select count(*) from example where name = ?", "query1")
	require.NotNil(t, row)

	require.True(t, row.Success())
	require.Nil(t, row.Err())
	require.Equal(t, int64(1), row.RowsAffected())
	require.Equal(t, "a row selected.", row.Message())
	require.Equal(t, 1, len(row.Values()))

	var val int
	if err := row.Scan(&val); err != nil {
		t.Fatalf("row scan fail; %s", err.Error())
	}
	require.Equal(t, 123, val)
}

func TestQuery(t *testing.T) {
	conn := newConn(t)
	defer conn.Close()
	rows, err := conn.Query(context.TODO(), "select * from example where name = ?", "query1")
	if err != nil {
		t.Fatalf("query fail, %q", err.Error())
	}
	defer rows.Close()

	require.True(t, rows.IsFetchable())
	require.Equal(t, int64(0), rows.RowsAffected())
	require.Equal(t, "success", rows.Message())

	columns, err := rows.Columns()
	if err != nil {
		t.Fatalf("columns error, %s", err.Error())
	}
	require.Equal(t, 3, len(columns))

	var name string
	var ts time.Time
	var value float64
	for rows.Next() {
		err := rows.Scan(&name, &ts, &value)
		if err != nil {
			t.Fatalf("rows scan error, %s", err.Error())
		}
	}

	require.Equal(t, "tag", name)
	require.Equal(t, time.Unix(0, 1).Nanosecond(), ts.Nanosecond())
	require.Equal(t, 3.14, value)
}

func TestAppend(t *testing.T) {
	conn := newConn(t)
	defer conn.Close()
	appender, err := conn.Appender(context.TODO(), "example")
	if err != nil {
		t.Fatalf("appender error, %s", err.Error())
	}
	require.NotNil(t, appender)

	for i := 0; i < 10; i++ {
		err := appender.Append(i)
		if err != nil {
			t.Fatalf("append fail, %s", err.Error())
		}
	}

	succ, fail, err := appender.Close()
	if err != nil {
		t.Errorf("appender close error, %s", err.Error())
	}
	require.Equal(t, int64(10), succ)
	require.Equal(t, int64(0), fail)
}
