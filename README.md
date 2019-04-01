# Execute Go On the Server, With WebAssembly

This project is a server and a small CLI that lets you write code locally, upload it, execute it on the server, and return the result to the CLI.

# Details

The [CLI](./cli) simply reads the specifies the file and does a `POST`s some JSON to the server. Then, the server does all the cool stuff:

- Creates a temp dir
- Writes some WASM-ey stuff into the temp dir:
  - Javascript to execute the WASM binary
  - A shall script to run a node process that executes the WASM binary
  - The uploaded Go code
- Executes the WASM inside a node process
- Fetches the output and sends back to the client in the HTTP response

>This executes Go code on the server in the same sandbox as if it were executed on the browser. So, it's as secure as if the same code were compiled into WebAssembly and executed in the browser.

# Run It Yourself

You need some stuff to run this:

- The ability to run Makefiles
- Go 1.11 or higher
- A recent version of Node.js

First, build the CLI:

```console
$ make build-cli
```

Next, run the server:

```console
$ make run-server
```

Then, in a new terminal window, from the same directory, run the CLI:

```console
$ ./runcli -file=sample/sample.go -server=http://localhost:8080
```

It succeeded if you see something like this:

```console
2019/04/01 13:45:25 curl -X 'POST' -d '{"Code":"cGFja2FnZSBtYWluCgppbXBvcnQgImZtdCIKCmZ1bmMgbWFpbigpIHsKCWZtdC5QcmludGxuKCJJdCBXb3JrZWQhIikKfQo="}' -H 'Content-Type: application/json' 'http://localhost:8080'
2019/04/01 13:45:26 It Worked!
```

# On The Shoulders Of Giants

This project wouldn't be possible without the work of all the folks who contributed WASM support across the ecosystem, including:

- All the contributors to Go, who added WASM support to the compiler and the de-facto standard Javascript template
- Mozilla, for leading the WASM spec work
- Plenty of other things

Finally, this project wouldn't be possible without the section in the WebAssembly Wiki called "Executing WebAssembly with Node.js". Check it out here: https://github.com/golang/go/wiki/WebAssembly#executing-webassembly-with-nodejs
