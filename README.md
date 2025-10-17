# wallet

## pack api

```shell
pwsh ./pack.api.ps1
```

## bitlong_pc

```shell
./build.pc.bat
```

## box

```shell
./build.box.bat
```

## build linux executable program

```shell
$env:CGO_ENABLED="0"; $env:GOOS="linux"; $env:GOARCH="amd64"; go build -ldflags="-s -w" -tags "litd autopilotrpc signrpc walletrpc chainrpc invoicesrpc watchtowerrpc neutrinorpc peersrpc btlapi" -o /path/to/output /path/to/main.go
```
