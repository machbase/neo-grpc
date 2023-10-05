package machrpc_test

import (
	context "context"
	"testing"

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

func TestNewClient(t *testing.T) {
	opts := []machrpc.Option{
		machrpc.WithServer(mock.MockServerAddr),
	}
	cli, err := machrpc.NewClient(opts...)
	if err != nil {
		t.Fatalf("new client: %s", err.Error())
	}

	ctx := context.TODO()

	// empty username, password
	conn, err := cli.Connect(ctx)
	require.NotNil(t, err)
	require.Equal(t, "no db user specified", err.Error())
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

func TestPing(t *testing.T) {
	cli := newClient(t)
	pinger, ok := cli.(spi.Pinger)
	if !ok {
		t.Fatal("client is not implements spi.Pinger")
	}
	dur, err := pinger.Ping()
	require.Nil(t, err)
	if dur == 0 {
		t.Fatal("ping failure")
	}
	cli.Close()
}

func TestQuery(t *testing.T) {
	conn := newConn(t)
	rows, err := conn.Query(context.TODO(), "select * from example where name = ?", "query1")
	if err != nil {
		t.Fatalf("query fail, %q", err.Error())
	}
	require.NotNil(t, rows)
	rows.Close()
	conn.Close()
}
