Message can carry a `replyTo` field that can specify a queue that will handle a reply. A reply queue and consumer are created on the RPC client before calling the RPC server. The server processes the request and puts a response in the reply queue.

The server has QoS settings on its channel so the exchange will go to the next available server (if it exists) when the target is processing a message.

A `CorrelationId` follows the transaction so the RPC client knows which server response traffic it owns. This should only become relevant when a server dies while the request is in flight.

```shell
go run rpc_server.go

# in another shell
go run rpc_client.go 30

# You can start multiple servers and they will share requests
```
