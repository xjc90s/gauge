# Experiment running Gauge as a webassembly in a browser

## How

### 1. Build

`GOOS=js GOARCH=wasm go build -o gauge.wasm ./gauge_js.go `

### 2. Serve

Use any local http server. For example, one can use [`goexec`](https://github.com/shurcooL/goexec) and run this command:

`'http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))'`