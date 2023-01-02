
## Settings for VSCode

<details>

  - `.vscode/settings.json`

  ```json
  {
      "protoc": {
          "options": [
              "--proto_path=./proto"
          ]
      },
      "files.exclude": {
          "vendor": true
      },
      "editor.tabSize": 4,
      "[go]": {
          "editor.tabSize": 4
      }
  }
  ```

</details>


## gRPC developer's info

<details>

## protobuf compiler

> https://grpc.io/docs/protoc-installation/

```
sudo apt install -y protobuf-compiler
```

```
brew install protobuf
```

- protoc-gen-go plugin

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## GRPC Gateway compiler

- [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

```
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

```
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

- [grpc-gateway with gin](https://blog.logrocket.com/guide-to-grpc-gateway/#using-grpc-gateway-with-gin)

### protobuf struct from/to json

```go
  buf, _ := ioutil.ReadAll(c.Request.Body)
  req := &protos.LoginRequest{}
  protojson.Unmarshal(buf, req)

  rsp, _ := s.Login(context.Background(), req)
  buf, _ = protojson.Marshal(rsp)

  c.Data(http.StatusOK, gin.MIMEJSON, buf)
```

</details>
