# IPboss
Linux daemon which manages IP based on local policies (process running, port open..)
It is fully integrated with AWS EC2 and it can add secondary IPs to the instances to be able to route traffic

### Build it
`CGO_ENABLED=0 GOOS=linux go build  -ldflags "-s"  -a -installsuffix cgo -o ipboss *.go`

### Dockerize it
`docker build .`
the default config is harmless and useful for testing, it adds IP to localhost

### Docs
see sample config file and/or read the code ;)

### Docker caveats
* you need the capability 'NET_ADMIN' to  managed interfaces, priviledged mode is overkill and unsafe
* you need to use networkmode 'host' to manage a server instance
* you need CA certificates bundle to autenticate AWS API servers
