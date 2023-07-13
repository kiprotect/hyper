# Willkommen!

**Hyper** bietet mehrere Server- und Client-Komponenten, die die Kommunikation im Hyper-Ökosystem verwalten und sichern. Hyper bietet vor allem zwei Kernkomponenten:

* Ein **Message-Broker / Mesh-Router-Dienst**, der Anfragen zwischen verschiedenen Akteuren im System weiterleitet und die gegenseitige Autorisierung und Authentifizierung sicherstellt.
* Ein **verteiltes Dienstverzeichnis**, das kryptografisch signierte Informationen über Akteure im System speichert und vom Message Broker für die Authentifizierung, die Dienstsuche und den Verbindungsaufbau verwendet wird.
* Ein **Overlay-Netzwerk**, mit dem Sie beliebige TCP- und UDP-Dienste über Ende-zu-Ende-verschlüsselte Kanäle verbinden können (in Arbeit).

Zusätzlich bietet es einen **TLS-Passthrough-Proxy-Dienst, der** eine direkte, Ende-zu-Ende-verschlüsselte Kommunikation zwischen Client-Endpunkten und Gesundheitsämtern ermöglicht.

## Erste Schritte

Für die ersten Schritte lesen Sie bitte das [Handbuch Erste Schritte]({{'getting-started'|href}}). Danach können Sie sich die [ausführliche Hyper-Dokumentation]({{'hyper.index'|href}}) ansehen. Wenn Sie den Proxy oder das Serviceverzeichnis ausführen möchten, können Sie sich die entsprechende [Proxy-Dokumentation]({{'proxy.index'|href}}) sowie die [Serviceverzeichnis-Dokumentation]({{'sd.index'|href}}) ansehen.

Wenn Sie auf ein Problem stoßen, [öffnen Sie](https://github.com/iris-connect/hyper) bitte [ein Issue auf Github](https://github.com/iris-connect/hyper), wo unsere Community Ihnen helfen kann.
