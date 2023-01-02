package machrpc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func MakeGrpcConn(addr string) (grpc.ClientConnInterface, error) {
	var conn *grpc.ClientConn
	var err error

	pwd, _ := os.Getwd()

	if strings.HasPrefix(addr, "unix://../") {
		addr = fmt.Sprintf("unix:///%s", filepath.Join(filepath.Dir(pwd), addr[len("unix://../"):]))
	} else if strings.HasPrefix(addr, "../") {
		addr = fmt.Sprintf("unix:///%s", filepath.Join(filepath.Dir(pwd), addr[len("../"):]))
	} else if strings.HasPrefix(addr, "unix://./") {
		addr = fmt.Sprintf("unix:///%s", filepath.Join(pwd, addr[len("unix://./"):]))
	} else if strings.HasPrefix(addr, "./") {
		addr = fmt.Sprintf("unix:///%s", filepath.Join(pwd, addr[len("./"):]))
	} else if strings.HasPrefix(addr, "/") {
		addr = fmt.Sprintf("unix://%s", addr)
	} else {
		addr = strings.TrimPrefix(addr, "http://")
		addr = strings.TrimPrefix(addr, "tcp://")
	}
	conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return conn, err
}
