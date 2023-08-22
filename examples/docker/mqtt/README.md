# MQTT Example

This example Docker compose setup shows how to use the QUIC channel to forward traffic to a remote service via the Hyper network. It defines the following services:

- **quic-1**: A Hyper service running the QUIC channel that is configured to forward local ports to a remote server via the **quic-2** proxy.
- **quic-2**: A Hyper service running the QUIC channel that is configured as a remote proxy, forwarding connections to the **mqtt-1** container.
- **sd-1**: The Hyper service directory for this setup.
- **mqtt-1**: A RabbitMQ-based MQTT broker.s
- **hd-1**: An admin container that initializes the service directory.

To run this setup, first run `make certs` in the main Hyper directory to generate all required TLS certificates. Then, simply run

```bash
docker compose up
```

This should create all containers and run them. You should then be able to connect to the RabbitMQ admin API via `curl` through the forwarded local port:

```bash
curl http://localhost:6666/rabbitmqadmin
# should return a JSON error response
```

That's it! You have successfully established connectivity to a remote service through the Hyper network's QUIC channel.