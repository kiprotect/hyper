# QUIC Channel

**Please note the QUIC channel is still a work in progress.**

The QUIC channel allows transmission of arbitrary TCP streams between two hosts.

## Configuration

The two example settings for QUIC are in the `quic-1` and `quic-2` settings folders (`settings/dev/roles/...`).
The `quic-1` server is configured to forward a single local TCP port (4444) to a port of the remote server (5555).

The corresponding service directory entries can be loaded as follows (remember to start the service directory first via `SD_SETTINGS=settings/dev/roles/sd-1 sd run`):

```bash
make sd-setup SD=quic
```

Then, you can simply start the two QUIC servers as follows:

```bash
# quic-1
HYPER_SETTINGS=settings/dev/roles/quic-1 hyper server run
# quic-2 (in a different terminal)
HYPER_SETTINGS=settings/dev/roles/quic-2 hyper server run
```

To simulate a local TCP server, you can e.g. use `ncat` as follows:

```bash
ncat -l 4444 --keep-open --exec "/bin/cat"
```

Now, you should be able to connect to local port `5555` and have all data echoed by the ncat server through the two QUIC servers:

```bash
> telnet localhost 5555
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
Hi
Hi
```

Congrats! You just set up the simplest possible QUIC channel between two hosts. You can add channel entries to the settings of the `quic-1` server to map additional ports.