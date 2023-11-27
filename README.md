
# neo-grpc

- `proto` gRPC .proto file specifies API of machbase-neo
- `machrpc` reference implementation of gRPC client of machbase-neo with helper functions.
- `driver` implmenetation of go sql/driver based gRPC

## For Gophers

### Install

```
go get -u github.com/machbase/neo-grpc
```

### How to use

```go
import "github.com/machbase/neo-grpc/machrpc"
```

```go
cli := machrpc.NewClient(
        machrpc.WithServer("127.0.0.1:5655"),
        machrpc.WithCertificate("/path/to/client_key.pem", "/path/to/client_cert.pem", "/path/to/server_cert.pem"),
        machrpc.WithQueryTimeout(5 * time.Second)
    )
if err := cli.Connect(); err != nil {
    panic(err)
}
defer cli.Disconnect()

sqlText := "insert into example (name, time, value) values (?, ?, ?)"
cli.Exec(sqlText, "tag-name", time.Now(), 3.1415)
```


## For gRPC Developers

### Settings for VSCode

  - `.vscode/settings.json`

  ```json
  {
      "protoc": {
          "options": [
              "--proto_path=./proto"
          ]
      }
  }
  ```

### protobuf compiler

> https://grpc.io/docs/protoc-installation/

- linux

```
sudo apt install -y protobuf-compiler
```

- macOS

```
brew install protobuf
```

- protoc-gen-go plugin

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
