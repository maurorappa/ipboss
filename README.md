# IPboss
Linux daemon which manages IP based on local policies (process running, port open..)

### Build it
`CGO_ENABLED=0 GOOS=linux go build  -ldflags "-s"  -a -installsuffix cgo -o ipboss *.go`

### Dockerize it
`docker build .`

### Docs
see sample config file and/or read the code ;)

### Docker caveats
* you need the capability 'NET_ADMIN' to  managed interfaces, priviledged mode is overkill and unsafe
* you need to use networkmode 'host' to manage a server instance
