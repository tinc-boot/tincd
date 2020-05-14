# TINCD

This is a supporting library for running tincd daemon from Go with tinc-web-boot supporting protocol.

## States

![image](https://user-images.githubusercontent.com/6597086/81945803-1cdc8780-9631-11ea-9a4e-64c772e3af8a.png)

Each state can be interrupted by error or canceled context.

## Greeting protocol

Over JSON-RPC 2.0 / HTTP on VPN IP on `CommunicationPort`  

![Untitled (7)](https://user-images.githubusercontent.com/6597086/81946280-c28ff680-9631-11ea-90c5-8a284af0c9ba.png)