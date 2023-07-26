# QUIC-Kanal

**Bitte beachten Sie, dass sich der QUIC-Kanal noch in der Entwicklung befindet.**

Der QUIC-Kanal ermöglicht die Übertragung beliebiger TCP-Streams zwischen zwei Hosts.

## Konfiguration

Die beiden Beispieleinstellungen für QUIC befinden sich in den Einstellungsordnern `quic-1` und `quic-2` (`settings/dev/roles/...`).
Der `quic-1` Server ist so konfiguriert, dass er einen einzelnen lokalen TCP-Port (4444) an einen Port des entfernten Servers (5555) weiterleitet.

Die entsprechenden Einträge im Dienstverzeichnis können wie folgt geladen werden (denken Sie daran, das Dienstverzeichnis zuerst über `SD_SETTINGS=settings/dev/roles/sd-1 sd run` zu starten):

```bash
make sd-setup SD=quic
```

Dann können Sie die beiden QUIC-Server einfach wie folgt starten:

```bash
# quic-1
HYPER_SETTINGS=settings/dev/roles/quic-1 hyper server run
# quic-2 (in a different terminal)
HYPER_SETTINGS=settings/dev/roles/quic-2 hyper server run
```

Um einen lokalen TCP-Server zu simulieren, können Sie z.B. `ncat` wie folgt verwenden:

```bash
ncat -l 4444 --keep-open --exec "/bin/cat"
```

Jetzt sollten Sie in der Lage sein, sich mit dem lokalen Port `5555` zu verbinden und alle Daten vom ncat-Server über die beiden QUIC-Server wiedergeben zu lassen:

```bash
> telnet localhost 5555
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
Hi
Hi
```

Herzlichen Glückwunsch! Sie haben soeben den einfachsten möglichen QUIC-Kanal zwischen zwei Hosts eingerichtet. Sie können in den Einstellungen des `quic-1` Servers Channel-Einträge hinzufügen, um weitere Ports zuzuordnen.
