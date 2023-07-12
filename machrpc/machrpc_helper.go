package machrpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func MakeGrpcInsecureConn(addr string) (grpc.ClientConnInterface, error) {
	return makeGrpcConn(addr, "")
}

func MakeGrpcTlsConn(addr string, caCertPath string) (grpc.ClientConnInterface, error) {
	return makeGrpcConn(addr, caCertPath)
}

func makeGrpcConn(addr string, caCertPath string) (grpc.ClientConnInterface, error) {
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

	if caCertPath == "" {
		return grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tlsCreds, err := loadTlsCreds(caCertPath)
		if err != nil {
			return nil, err
		}
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(tlsCreds))
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}

func loadTlsCreds(certPath string) (credentials.TransportCredentials, error) {
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		return nil, fmt.Errorf("fail to load server CA cert")
	}
	tlsConfig := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}
	return credentials.NewTLS(tlsConfig), nil
}
