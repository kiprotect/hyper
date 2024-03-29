# Welcome!

The **Hyper** provides several server and client components that manage and secure the communication in the Hyper ecosystem. Notably, the Hyper provides two core components:

* A **message broker / mesh router** services that transmits requests between different actors in the system and ensures mutual authorization and authentication.
* A **distributed service directory** that stores cryptographically signed information about actors in the system, and is used by the message broker for authentication, service discovery and connection establishment.
* An **overlay network** that enables you to connect arbitrary TCP & UDP services over end-to-end encrypted channels (work in progress).

In addition it provides a **TLS passthrough proxy service** that enables direct, end-to-end encrypted communication between client endpoints and health departments.

## Getting Started

To get started please read the [Getting Started Guide]({{'getting-started'|href}}). After that you can check out the [detailed Hyper documentation]({{'hyper.index'|href}}). If you want to run the proxy or service directory you can check out the respective [proxy documentation]({{'proxy.index'|href}}) as well as [service directory documentation]({{'sd.index'|href}}).

If you encounter a problem please [open an issue on Github](https://github.com/iris-connect/hyper) where our community can help you.