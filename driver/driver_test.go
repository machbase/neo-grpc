package driver_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/machbase/neo-grpc/mock"
	"github.com/stretchr/testify/require"
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

func connect(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("machbase", fmt.Sprintf("tcp://sys:manager@%s", mock.MockServerAddr))
	if err != nil {
		t.Fatalf("db connection failure %q", err.Error())
	}
	return db
}

func TestConnectFail(t *testing.T) {
	db, err := sql.Open("machbase", fmt.Sprintf("tcp://%s", mock.MockServerAddr))
	require.NotNil(t, err)
	require.Equal(t, "invalid username or password", err.Error())
	require.Nil(t, db)
}

func TestQuery(t *testing.T) {
	db := connect(t)
	defer db.Close()

	rows, err := db.Query(`select * from example where name = ?`, "query1")
	require.Nil(t, err)
	require.NotNil(t, rows)
	rows.Close()

	conn, err := db.Conn(context.TODO())
	require.Nil(t, err)
	require.NotNil(t, conn)

	rows, err = conn.QueryContext(context.TODO(), `select * from example where name = ?`, "query1")
	require.Nil(t, err)
	require.NotNil(t, rows)
	rows.Close()

	conn.Close()
}
