# veil-evasion-api
An HTTP API and client for Veil-Evasion
![](https://raw.githubusercontent.com/tomsteele/veil-evasion-api/master/docs/client.png)

## Install
The easiest way to install is to use the docker image. This will start Veil-Evasion's rpc server for you.

#### Using Docker
```
$ docker pull tomsteele/veil-evasion-api
$ docker run --rm -e ADMIN_USER=someuser -e ADMIN_PASS=somesecret -p 80:80 tomsteele/veil-evasion
```

#### Building
If you don't wish to use docker, you can also build it yourself. These instructions will also serve as development instructions. You will need to install Go and node.

Somewhere start Veil-Evasion's rpc server:
```
$ ./Veil-Evasion.py --rpc
```

The web client needs to be built with node:
```
$ cd client
$ npm i
$ npm run-script build
$ mv dist ../public
```

This project uses godep. To compile and run the server, do the following (change environment variables as needed:
```
$ go get github.com/tools/godep
$ export VEIL_LISTENER=localhost:4242
$ export VEIL_OUTPUT_DIR=/usr/share/veil-output
$ export SERVER_LISTENER=0.0.0.0:8000
$ export ADMIN_USER=admin
$ export ADMIN_PASS=secret
$ godep go run main.go // to run
$ godep go build // to build
```
## Using in your Go program
The handlers for communicating with the API are modular. You can embed them in your own Go server. Example:
```
func main() {
	c, _ := jsonrpc.Dial("tcp", "localhost:4242")
	h := handlers.New(&handlers.H{C: c})
	r := mux.NewRouter()
	r.HandleFunc("/api/version", h.Version).Methods("GET")
	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":3000")
}
```

## Security
Before deploying there are a few things you should consider
 - The server uses basic auth with a shared username and password for authentication.
 - Payloads do not require authentication to assist with delivery.
 - CSRF is prevented by enforcing content-type of application/json
 - No validation is performed on options passed to Veil-Evasion. It's possible that there is some command injection in this process. I haven't looked. It's probably best to assume that if you give someone access to the API that they could do this.
 - NO TLS/SSL. Encryption is good, and you should use it. I would suggest deploying this behind a reverse proxy such as nginx or caddy.
