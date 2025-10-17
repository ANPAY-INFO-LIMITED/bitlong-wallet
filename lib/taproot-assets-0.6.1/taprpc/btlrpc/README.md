# btlrpc

## 安装依赖

- protoc

https://github.com/protocolbuffers/protobuf/releases

- go工具
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

- 需要REST反向代理时*可选*
```bash
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

## 生成proto文件

- 此前请先进入目录
```bash
cd taprpc
```

- 生成go文件

```bash
protoc --go_out . --go_opt paths=source_relative --go-grpc_out require_unimplemented_servers=false:. --go-grpc_opt paths=source_relative btlrpc/btl.proto
```

- 关于`require_unimplemented_servers`请参考

https://stackoverflow.com/questions/65079032/grpc-with-mustembedunimplemented-method

## 注意事项

- 为防止`lightning-terminal/subservers/taproot-assets_btlapi.go`调用`btl_grpc.pb.go`的`RegisterBtlServer`中，`t.testEmbeddedByValue()`时出现由空指针导致的`panic`，请在生成文件后手动将其注释，而只保留`s.RegisterService(&Btl_ServiceDesc, srv)`，此为目前暂时的解决方法。

## 另一些示例命令

- 生成 protos
```bash
protoc --go_out . --go_opt paths=source_relative --go-grpc_out . --go-grpc_opt paths=source_relative btlrpc/btl.proto
```

- 生成 REST 反向代理
```bash
protoc --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt grpc_api_configuration=btlrpc/btl.yaml btlrpc/btl.proto
```

- 生成详细描述 REST API 的 swagger 文件
```bash
protoc --openapiv2_out .  --openapiv2_opt logtostderr=true --openapiv2_opt grpc_api_configuration=${annotationsFile} --openapiv2_opt json_names_for_fields=false btlrpc/btl.proto
```
