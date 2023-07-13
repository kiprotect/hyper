# Integration

Die Integration in die IRIS-Infrastruktur ist einfach (hoffen wir). Zunächst müssen Sie den `hyper` Server zusammen mit den Einstellungen und Zertifikaten, die wir Ihnen zur Verfügung gestellt haben, einrichten. Dazu laden Sie einfach die neueste Version von `hyper` von unserem Server herunter, entpacken das Archiv mit den Einstellungen, das wir Ihnen zur Verfügung gestellt haben, und führen

```bash
HYPER_SETTINGS=path/to/settings hyper server run
```

Dies sollte einen lokalen JSON-RPC-Server auf Port `5555` öffnen, mit dem Sie sich über TLS verbinden können (dazu müssen Sie das Root-CA-Zertifikat zu Ihrer Zertifikatskette hinzufügen). Dieser Server ist Ihr Gateway zu allen IRIS-Diensten. Schauen Sie einfach nach den Diensten, die ein bestimmter Betreiber anbietet, und senden Sie eine Anfrage, die den Namen des Betreibers und die Dienstmethode enthält, die Sie aufrufen möchten. Um z. B. mit dem Dienst "locations" zu interagieren, der vom Betreiber "ls-1" bereitgestellt wird, würden Sie einfach eine JSON-RPC-Nachricht wie diese senden:

```json
{
	"method": "ls-1.add",
	"id": "1",
	"params": {
		"name": "Ginos",
		"id": "af5ca4da5caa"
	},
	"jsonrpc": "2.0"
}
```

Das Gateway kümmert sich darum, diese Nachricht an den richtigen Dienst weiterzuleiten und eine Antwort an Sie zurückzuschicken.

Wenn Sie Anfragen von anderen Diensten im IRIS-Ökosystem akzeptieren möchten, können Sie die `jsonrpc_client` verwenden. Dabei geben Sie einfach einen API-Endpunkt an, an den eingehende Anfragen mit der gleichen Syntax wie oben zugestellt werden sollen.

Das war's!

## Asynchrone Anrufe

Die Anrufe, die wir oben gesehen haben, waren alle synchron, d.h. ein Anruf führte zu einer direkten Antwort. Manchmal sind jedoch asynchrone Anrufe erforderlich, z.B. weil die Beantwortung der Anrufe Zeit in Anspruch nimmt. Wenn Sie einen asynchronen Aufruf an einen anderen Dienst tätigen, erhalten Sie zunächst eine Bestätigung zurück. Sobald der von Ihnen aufgerufene Dienst eine Antwort bereit hat, sendet er diese über das Netzwerk `hyper` an Sie zurück, wobei er dieselbe `id` verwendet, die Sie angegeben haben (wodurch Sie die Antwort Ihrer Anfrage zuordnen können). Ebenso können Sie auf Aufrufe von anderen Diensten asynchron reagieren, indem Sie die Antwort einfach an Ihren lokalen JSON-RPC-Server mit dem Methodennamen `respond` (ohne den Namen des Dienstes) weiterleiten. Vergessen Sie nicht, dieselbe `id` anzugeben, die Sie mit der ursprünglichen Anfrage erhalten haben, da diese die "Rücksendeadresse" der Anfrage enthält.

## Integration Beispiel

Um eine konkrete Vorstellung von der Integration mit der IRIS-Infrastruktur unter Verwendung des Hyper-Servers zu erhalten, haben wir ein einfaches Demo-Setup erstellt, das alle Komponenten veranschaulicht. Die Demo besteht aus drei Komponenten:

* Ein `hyper` Server, der eine `health department` simuliert, namens `hd-1`
* Eine `hyper` Server-Simulation eines Betreibers, der einen "Standort"-Dienst anbietet, namens `ls-1`
* Der tatsächliche vom Betreiber angebotene Ortungsdienst `hyper-ls` `ls-1`

## Aufstehen und loslegen

Lesen Sie bitte zuerst in der README nach, wie Sie alle notwendigen TLS-Zertifikate erstellen und die Software bauen. Starten Sie dann die einzelnen Dienste auf verschiedenen Terminals:

```bash
# run the `hyper` server of the "locations" operator ls-1
HYPER_SETTINGS=settings/dev/roles/ls-1 hyper --level debug server run
# run the `hyper` server of the health department hd-1
HYPER_SETTINGS=settings/dev/roles/hd-1 hyper --level debug server run
# run the "locations" service
hyper-ls
```

## Testen

Jetzt sollte Ihr System einsatzbereit sein. Der Demo-Dienst "locations" bietet eine einfache, authentifizierungsfreie JSON-RPC-Schnittstelle mit zwei Methoden: `add`, die einen Ort zur Datenbank hinzufügt, und `lookup`, die einen Ort anhand seiner `name` nachschlägt. Zum Beispiel, um dem Dienst einen Ort hinzuzufügen:

```bash
curl --cacert settings/dev/certs/root.crt --resolve hd-1:5555:127.0.0.1 https://hd-1:5555/jsonrpc --header "Content-Type: application/json" --data '{"method": "ls-1.add", "id": "1", "params": {"name": "Ginos", "id": "af5ca4da5caa"}, "jsonrpc": "2.0"}' 2>/dev/null | jq 
```

Dies sollte eine einfache JSON-Antwort zurückgeben:

```json
{
  "jsonrpc": "2.0",
  "result": {
    "_": "ok"
  },
  "id": "1"
}
```

Die Anfrage ging zunächst an den `hyper` Server des Gesundheitsamtes, wurde zunächst über gRPC an den `hyper` Server von `ls-1` weitergeleitet und dann an die JSON-RPC API des dort laufenden lokalen `hyper-ls` Dienstes übergeben. Das Ergebnis wurde dann über die gesamte Kette zurückgereicht.

Sie können auch eine Suche nach dem Ort durchführen, den Sie gerade hinzugefügt haben:

```bash
curl --cacert settings/dev/certs/root.crt --resolve hd-1:5555:127.0.0.1 https://hd-1:5555/jsonrpc --header "Content-Type: application/json" --data '{"method": "ls-1.lookup", "id": "1", "params": {"name": "Ginos"}, "jsonrpc": "2.0"}' 2>/dev/null | jq .
```

die Folgendes zurückgeben sollte

```json
{
  "jsonrpc": "2.0",
  "result": {
    "id": "af5ca4da5caa",
    "name": "Ginos"
  },
  "id": "1"
}
```

Daher ist die Interaktion mit dem entfernten "locations"-Dienst genauso wie der Aufruf eines lokalen JSON-RPC-Dienstes, außer dass Sie den Namen des Operators, der den Dienst ausführt, `ls-1.lookup`, angeben, anstatt einfach `lookup` aufzurufen.

