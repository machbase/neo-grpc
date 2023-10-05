package mock

import (
	context "context"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/machbase/neo-grpc/machrpc"
	"google.golang.org/grpc"
)

type MockServer struct {
	machrpc.MachbaseServer
	svr *grpc.Server

	sessionCounter int32
}

var _ machrpc.MachbaseServer = &MockServer{}

var MockServerAddr = "127.0.0.1:5655"

func (ms *MockServer) Start() error {
	svrOptions := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(int(1 * 1024 * 1024 * 1024)),
		grpc.MaxSendMsgSize(int(1 * 1024 * 1024 * 1024)),
	}
	ms.svr = grpc.NewServer(svrOptions...)
	machrpc.RegisterMachbaseServer(ms.svr, ms)
	lsnr, err := net.Listen("tcp", "127.0.0.1:0")
	MockServerAddr = lsnr.Addr().String()
	if err != nil {
		return err
	}
	go ms.svr.Serve(lsnr)
	return nil
}

func (ms *MockServer) Stop() {
	ms.svr.Stop()

	if ms.sessionCounter != 0 {
		panic(fmt.Errorf("WARN!!!! connection leak!!! There are %d sessions remained", ms.sessionCounter))
	}
}

func (ms *MockServer) CountSession() int {
	return int(ms.sessionCounter)
}

func (ms *MockServer) UserAuth(ctx context.Context, req *machrpc.UserAuthRequest) (*machrpc.UserAuthResponse, error) {
	auth := true
	reason := "success"
	if req.LoginName != "sys" || req.Password != "manager" {
		auth = false
		reason = "invalid username or password"
	}
	return &machrpc.UserAuthResponse{
		Success: auth,
		Reason:  reason,
		Elapse:  "1ms.",
	}, nil
}

func (ms *MockServer) Ping(ctx context.Context, req *machrpc.PingRequest) (*machrpc.PingResponse, error) {
	return &machrpc.PingResponse{
		Success: true,
		Reason:  "success",
		Elapse:  "1ms.",
		Token:   req.Token,
	}, nil
}

func (ms *MockServer) Explain(ctx context.Context, req *machrpc.ExplainRequest) (*machrpc.ExplainResponse, error) {
	return nil, nil
}

func (ms *MockServer) Exec(ctx context.Context, req *machrpc.ExecRequest) (*machrpc.ExecResponse, error) {
	return nil, nil
}

func (ms *MockServer) QueryRow(ctx context.Context, req *machrpc.QueryRowRequest) (*machrpc.QueryRowResponse, error) {
	return nil, nil
}

func (ms *MockServer) Query(ctx context.Context, req *machrpc.QueryRequest) (*machrpc.QueryResponse, error) {
	ret := &machrpc.QueryResponse{
		Success:      false,
		Reason:       "not implmeneted",
		Elapse:       "1ms.",
		RowsHandle:   &machrpc.RowsHandle{},
		RowsAffected: 0,
	}
	params := machrpc.ConvertPbToAny(req.Params)
	switch req.Sql {
	case `select * from example where name = ?`:
		if len(params) == 1 && params[0] == "query1" {
			ret.RowsHandle = &machrpc.RowsHandle{
				Handle: "query1",
			}
			ret.Success, ret.Reason = true, "success"
		} else {
			ret.Reason = fmt.Sprintf("not implemented %+v", params)
		}
		if ret.Success {
			atomic.AddInt32(&ms.sessionCounter, 1)
		}
	}
	return ret, nil
}

func (ms *MockServer) Columns(ctx context.Context, rows *machrpc.RowsHandle) (*machrpc.ColumnsResponse, error) {
	return nil, nil
}

func (ms *MockServer) RowsFetch(ctx context.Context, rows *machrpc.RowsHandle) (*machrpc.RowsFetchResponse, error) {
	return nil, nil
}

func (ms *MockServer) RowsClose(ctx context.Context, rows *machrpc.RowsHandle) (*machrpc.RowsCloseResponse, error) {
	atomic.AddInt32(&ms.sessionCounter, -1)
	return &machrpc.RowsCloseResponse{
		Success: true,
		Reason:  "success",
		Elapse:  "1ms.",
	}, nil
}

func (ms *MockServer) Appender(ctx context.Context, req *machrpc.AppenderRequest) (*machrpc.AppenderResponse, error) {
	return nil, nil
}

func (ms *MockServer) Append(stream machrpc.Machbase_AppendServer) error {
	return nil
}

func (ms *MockServer) GetServerInfo(ctx context.Context, req *machrpc.ServerInfoRequest) (*machrpc.ServerInfo, error) {
	return nil, nil
}

func (ms *MockServer) GetServicePorts(ctx context.Context, req *machrpc.ServicePortsRequest) (*machrpc.ServicePorts, error) {
	return nil, nil
}
