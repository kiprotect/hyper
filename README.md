**This software is still a work in progress and not ready for production use!**

# IRIS Endpoint Server (EPS)

This repository contains the code of the IRIS endpoint server (EPS), which manages the communication between different actors in the IRIS ecosystem. It provides a gRPC server & client to exchange messages between different actors, as well as a JSON-RPC API for location communication.

## Getting Started

To build the `eps` binary, simply run

```
make
```

For testing and development you'll also need TLS certificates, which you can generate with

```
make test-and-dev-certs
```

Please note that you need `openssl` on your system for this to work. This will generate all required certificates and put them in the `settings/dev/certs` and `settings/dev/test` folders. Please do not use these certificates in a production setting and do not check them into version control.

Please see below for additional dependencies you might need to install for various purposes (e.g. to recompile protobuf code).

## Running The Server

To run the development EPS server simply run (from the main directory)

```
EPS_SETTINGS=settings/dev eps server run
```

For this to work you need to ensure that your `GOPATH` is in your `PATH`. This will open the JSON RPC server and (depending on the settings) also a gRPC server.

## Running The Test Client

To run the development EPS test client simply run (from the main directory)

```
EPS_SETTINGS=settings/dev eps client test
```

This will currently try to connect to the development server via gRPC and send a message to it. In the future it will perform various connection and functional tests.

## Testing

To run the tests

```
make test # run normal tests
make test-races # test for race conditions
```

## Benchmarks

To run the benchmarks

```
make bench
```

## Copyright Headers

You can generate and update copyright headers as follows

```
make copyright-headers
```

This will add appropriate headers to all Golang files. You can edit the generation and affected file types directly in the script (in `.scripts`). You should run this before committing code. Please note that any additional comments that appear directly at the top of the file will be replaced by this.

## License

Currently this code is licensed under Affero GPL 3.0.

## Development Requirements

If you make modifications to the protocol buffers (`.proto` files) you need to recompile them using `protoc`. To install this on Debian/Ubuntu systems:

```
sudo apt install protobuf-compiler
```

To generate TLS certificates for testing and development you need to have `openssl` installed.