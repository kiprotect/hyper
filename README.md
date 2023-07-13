# Hyper

This repository contains the code of the Hyper mesh router & overlay network system, which provides a programmable network & communication layer for distributed systems. It provides a gRPC server & client to exchange messages between different actors (`hyper`), a service directory (`sd`) that provides a public registry of signed information about actors in the system, and a TLS passthrough proxy (`proxy`) that enables actors to directly receive data from each other and third-parties via TLS.

Please also check out the [documentation](https://kiprotect.github.io/hyper/docs/) for more detailed information. Please check the `docs` subfolder for instructions on how to build/view the documentation locally.

## Getting Started

Please ensure your Golang version is recent enough (>=1.13) before you attempt to build the software. 

To build the `hyper`, `sd` and `proxy` binaries, simply run

```bash
make
```

For testing and development you'll also need TLS certificates, which you can generate via

```bash
make certs
```

Please note that you need `openssl` on your system for this to work. This will generate all required certificates and put them in the `settings/dev/certs` and `settings/dev/test` folders. Please do not use these certificates in a production setting and do not check them into version control. **Warning:** Running  this command again will delete an re-create all certificates from scratch.

Please see below for additional dependencies you might need to install for various purposes (e.g. to recompile protobuf code).

To build the example services (e.g. the "locations" services `hyper-ls`) simply run

```bash
make examples
```

## Defining Settings

The `hyper` binary will look for settings in a list of colon (`:`) separated directories as defined by the `HYPER_SETTINGS` environment variable (or, if it is undefined in the `settings` subdirectory of the current directory). The development settings include an environment-based variable `HYPER_OP` that allows you to use different certificates for testing. You should define these variables before running the development server:

```bash
export HYPER_SETTINGS=`readlink -f settings/dev`
export HYPER_OP=hd-1 # run server as the 'hd-1' operator
```

You can also source these things from the local `.dev-setup` script, which includes everything you need to get started:

```bash
source .dev-setup # load all development environment variables
```

There are also role-specific development/test settings in the `settings/dev/roles` directory. Those can be used to set up multiple Hyper servers and test the communication between them. Please have a a look at the [integration guidelines](docs/integration.md) for more information about this.

**Important: The settings parser includes support for variable replacement and many other things. But with great power comes great responsibility and attack surface, so make sure you only feed trusted YAML input to it, as it is not designed to handle untrusted or potentially malicious settings.**

## Running The Service Directory

All Hyper servers rely on the service directory (SD) to discover each other and learn about permissions, certificates and other important settings. For development, you can either use a JSON-based service directory, or run the service directory API like this:

```bash
SD_SETTINGS=settings/dev/roles/sd-1 sd run
```

To initialize the service directory you can upload the JSON-based directory:

```bash
# for development
make sd-setup
# for testing
make sd-test-setup
```

This should give you a fully functional API-based service directory with certificate and service information.

## Running The Hyper Server

To run the development Hyper server simply run (from the main directory)

```bash
HYPER_SETTINGS=settings/dev/roles/hd-1 hyper server run
```

This will run the Hyper server for the role `hd-1` (simulating a health department in the system). For this to work you need to ensure that your `GOPATH` is in your `PATH`. This will open the JSON RPC server and (depending on the settings) also a gRPC server.

## Running The Proxy Servers

To run the public and private proxy servers simply run (from the main directory)

```bash
# private proxy server
PROXY_SETTINGS=settings/dev/roles/private-proxy-1 proxy run private
# public proxy server
PROXY_SETTINGS=settings/dev/roles/private-proxy-1 proxy run public
```

## Testing

To run the tests

```bash
make test # run normal tests
make test-races # test for race conditions
```

## Benchmarks

To run the benchmarks

```bash
make bench
```

## Debugging

If you're stuck debugging a problem please have a look at the [debugging guidelines](docs/debugging.md), which contain a few pointers that might help you to pinpoint problems in the system.

## Copyright Headers

You can generate and update copyright headers as follows

```bash
make copyright
```

This will add appropriate headers to all Golang files. You can edit the generation and affected file types directly in the script (in `.scripts`). You should run this before committing code. Please note that any additional comments that appear directly at the top of the file will be replaced by this.

## License

Currently this code is licensed under Affero GPL 3.0.

## Development Requirements

If you make modifications to the protocol buffers (`.proto` files) you need to recompile them using `protoc`. To install this on Debian/Ubuntu systems:

```bash
sudo apt install protobuf-compiler
```

To generate TLS certificates for testing and development you need to have `openssl` installed.

## Deployment

You can easily deploy the server as a service using `systemd` or Docker. Specific documentation coming up soon.

# Feedback

If you have any questions [just contact us](mailto:hyper@kiprotect.com).

# Participation

We are happy about your contribution to the project! In order to ensure compliance with the licensing conditions and the future development of the project, we require a signed contributor license agreement (CLA) for all contributions in accordance with the [Harmony standard](http://selector.harmonyagreements.org). Please sign the corresponding document for [natural persons](.clas/hyper-individual.pdf) or for [organizations](.clas/hyper-entity.pdf) and send it to [us](mailto:hyper@kiprotect.com).

## Supporting organizations

- The software on which Hyper is based, originally named EPS (Endpoint System), was generously suppored by the Björn Steiger Stiftung SbR - https://www.steiger-stiftung.de. The source code of the system is available at https://github.com/iris-connect/eps.
