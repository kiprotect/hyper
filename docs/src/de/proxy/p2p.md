# Peer-to-Peer (P2P) Proxy

Der Proxy unterstützt auch einen Peer-to-Peer (P2P) Modus, der es zwei EPS-Servern ermöglicht, eine sichere Verbindung über den Proxy herzustellen.

## Vermittlung von Verbindungen

Der Vermittlungsprozess für eine P2P-Verbindung läuft wie folgt ab:

* EPS-Server A möchte eine Verbindung zu EPS-Server B herstellen.
* A erkennt anhand des Verzeichniseintrags von B, dass er nur über den Proxy P erreichbar ist.
* A sendet eine `connectionRequest` Nachricht über das EPS-System an P und gibt dabei den gRPC-Serverkanal von B als Empfänger an.
* P erstellt ein Token und sendet über das EPS-System eine Nachricht an B, in der er die Anfrage von A weiterleitet und einen Proxy-Endpunkt angibt, mit dem er sich verbinden soll.
* B empfängt die Nachricht und leitet sie an den entsprechenden Kanal weiter, der sie verarbeitet und eine Verbindung zum Endpunkt von P herstellt, wobei der Token als Routing-Schlüssel gesendet wird.
* P empfängt die Verbindung von B und speichert sie.
* P sendet eine Bestätigung an A, die den Token und denselben Endpunkt enthält.
* A verbindet sich mit dem Endpunkt von P und sendet ebenfalls das Token.
* P nimmt die Verbindung von A an, ruft die passende Verbindung von B ab und vermittelt den Datenverkehr zwischen ihnen.

## Testen Sie

Um eine Testinfrastruktur einzurichten, führen Sie einfach (in verschiedenen Shells):

```bash
# run the service directory
SD_SETTINGS=settings/dev/roles/sd-1 sd run
# run the public proxy
PROXY_SETTINGS=settings/dev/roles/public-proxy-1 proxy run public

# run all EPS servers
EPS_SETTINGS=settings/dev/roles/hd-1 eps server run
EPS_SETTINGS=settings/dev/roles/hd-2 eps server run
EPS_SETTINGS=settings/dev/roles/public-proxy-eps-1 eps server run
```

Stellen Sie sicher, dass Sie `make sd-setup` ausführen, um das Dienstverzeichnis mit den erforderlichen Einträgen zu aktualisieren. Anschließend sollten Sie in der Lage sein, über den Proxy einen Ping vom HD-2 Server über den HD-1 JSON-RPC Server anzufordern:

```bash
curl --cert settings/dev/certs/hd-1.crt --key settings/dev/certs/hd-1.key --cacert settings/dev/certs/root.crt --resolve hd-1:5555:127.0.0.1 https://hd-1:5555/jsonrpc --header "Content-Type: application/json" --data '{"method": "hd-2._ping", "id": "1", "params": {}, "jsonrpc": "2.0"}' | jq .

```
