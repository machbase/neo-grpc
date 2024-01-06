
# neo-grpc

- `proto` gRPC .proto file specifies API of machbase-neo
- `machrpc` reference implementation of gRPC client of machbase-neo with helper functions.
- `driver` implmenetation of go sql/driver based gRPC

## gRPC client for machbase-neo

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

## Go Driver for machbase-neo

### Install

Tp add driver into your `go.mod`, execute this command on your working directory.

```sh
go get -u github.com/machbase/neo-grpc/driver
```

### Import

Add import statement in your source file.

```go
import (
    "database/sql"
    _ "github.com/machbase/neo-grpc/driver"
)
```

> The package name `github.com/machbase/neo-grpc/driver` implies that the driver is implemented over gRPC API of machbase-neo. See [gRPC API](/docs/api-grpc/) for more about it.

Let's load sql driver and connect to server.

```go
db, err := sql.Open("machbase", "127.0.0.1:5655")
if err != nil {
    panic(err)
}
defer db.Close()
```

`sql.Open()` is called, `machbase` as driverName and second argument is machbase-neo's gRPC address.

### Insert and Query

Since we get `*sql.DB` successfully, write and read data can be done by Go standard sql package.

```go
var tag = "tag01"
_, err = db.Exec("INSERT INTO EXAMPLE (name, time, value) VALUES(?, ?, ?)", tag, time.Now(), 1.234)
```

Use `db.QueryRow()` to count total records of same tag name.

```go
var count int
row := db.QueryRow("SELECT count(*) FROM EXAMPLE WHERE name = ?", tag)
row.Scan(&count)
```

To iterate result of query, use `db.Query()` and get `*sql.Rows`.

{{< callout type="warning" >}}
It is import not to forget releasing `rows` as soon as possible when you finish the work to prevent leaking resource.<br/>
General pattern is using `defer rows.Close()`.
{{< /callout >}}

```go
rows, err := db.Query("SELECT name, time, value FROM EXAMPLE WHERE name = ? ORDER BY TIME DESC", tag)
if err != nil {
    panic(err)
}
defer rows.Close()

for rows.Next() {
    var name string
    var ts time.Time
    var value float64
    rows.Scan(&name, &ts, &value)
    fmt.Println("name:", name, "time:", ts.Local().String(), "value:", value)
}
```

### Full source code 

```go
package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	driver "github.com/machbase/neo-grpc/driver"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	serverAddr := "127.0.0.1:5655"
	serverCert := filepath.Join(homeDir, ".config", "machbase", "cert", "machbase_cert.pem")
	// This example substitute server's key & cert for the client's key, cert.
	// It is just for the briefness of sample code
	// Client applicates **SHOULD** issue a certificate for each one.
	// Please refer to the "API Authentication" section of the documents.
	clientKey := filepath.Join(homeDir, ".config", "machbase", "cert", "machbase_key.pem")
	clientCert := filepath.Join(homeDir, ".config", "machbase", "cert", "machbase_cert.pem")

	// register machbase-neo datasource
	driver.RegisterDataSource("neo", &driver.DataSource{
		ServerAddr: serverAddr,
		ServerCert: serverCert,
		ClientKey:  clientKey,
		ClientCert: clientCert,
	})

	db, err := sql.Open("machbase", "neo")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// INSERT with Exec
	tag := "tag01"
	_, err = db.Exec("INSERT INTO EXAMPLE (name, time, value) VALUES(?, ?, ?)", tag, time.Now(), 1.234)
	if err != nil {
		panic(err)
	}

	// QueryRow
	row := db.QueryRow("SELECT count(*) FROM EXAMPLE WHERE name = ?", tag)
	if row.Err() != nil {
		panic(row.Err())
	}
	var count int
	if err = row.Scan(&count); err != nil {
		panic(err)
	}
	fmt.Println("count:", count)

	// Query
	rows, err := db.Query("SELECT name, time, value FROM EXAMPLE WHERE name = ? ORDER BY TIME DESC", tag)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var ts time.Time
		var value float64
		rows.Scan(&name, &ts, &value)
		fmt.Println("name:", name, "time:", ts.Local().String(), "value:", value)
	}
}
```