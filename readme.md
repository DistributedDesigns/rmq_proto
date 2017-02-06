Rabbit MQ Prototype
====

How is rabbit formed?

## Demo
A  `quote_manager` accepts requests for new quotes and broadcasts the new value when the quote resolves. The `quote_cache` listens to quote broadcasts and logs the value. The `quote_requester` passes requests to the `quote_manager` and logs the value for the quote when it's broadcast. The broadcast does not have to originate from the request that `quote_requester` pushed.

### Run it
0. Start a local rabbit server. This requires a running instance of Docker.
  - `make run`
  - You can monitor the rabbit server at [localhost:8080](http://localhost:8080/)
0. `go run quote_manager.go`
0. `go run quote_cache.go`
0. `go run quote_requester.go AAPL GOOG SLOW`
  - `SLOW` will always take 20sec for a response from `quote_manager`
  - You can fake a simultaneous request for `SLOW` by publishing a message directly to the `quote_broadcast` queue with the [RMQ management interface](http://localhost:8080/#/exchanges/%2F/quote_broadcast)
    - Routing key: `SLOW`
    - Payload: `someone_else,SLOW,123.45`

### Logs
We want new quotes for `AAPL`, `GOOG` and `SLOW`.

`SLOW` is manually updated at 19:52:03.
```
$go run quote_requester.go AAPL GOOG SLOW     
2017/02/05 19:51:58  [.] Getting updates for [AAPL GOOG SLOW]
2017/02/05 19:51:58  [-] Waiting for updates to GOOG
2017/02/05 19:51:58  [-] Waiting for updates to SLOW
2017/02/05 19:51:58  [↑] Requesting quote for GOOG
2017/02/05 19:51:58  [↑] Requesting quote for SLOW
2017/02/05 19:51:58  [-] Waiting for updates to AAPL
2017/02/05 19:51:58  [↑] Requesting quote for AAPL
2017/02/05 19:51:58  [↓] Received: jappleseed,GOOG,653.33
2017/02/05 19:51:58  [x] Got update for GOOG
2017/02/05 19:52:00  [↓] Received: jappleseed,AAPL,342.61
2017/02/05 19:52:00  [x] Got update for AAPL
2017/02/05 19:52:03  [↙] Intercepted: someone_else,SLOW,123.45 <-- manual update
2017/02/05 19:52:03  [x] Got update for SLOW
<< requested updated is uncaught >>
```
```
$ go run quote_cache.go
2017/02/05 19:51:58  [-] Waiting for broadcasts
2017/02/05 19:51:58  [↓] Update: jappleseed,GOOG,653.33
2017/02/05 19:52:00  [↓] Update: jappleseed,AAPL,342.61
2017/02/05 19:52:03  [↓] Update: someone_else,SLOW,123.45 <-- manual update
2017/02/05 19:52:18  [↓] Update: jappleseed,SLOW,977.21   <-- requested update
```
```
$ go run quote_manager.go                 
2017/02/05 19:51:54  [-] Waiting for new pending quotes
2017/02/05 19:51:54  [-] Monitoring quote_req queue
2017/02/05 19:51:58  [↓] Received a quote request: jappleseed,GOOG
2017/02/05 19:51:58  [.] New pending quote request
2017/02/05 19:51:58  [-] Waiting for 0 sec
2017/02/05 19:51:58  [.] Got a response: jappleseed,GOOG,653.33
2017/02/05 19:51:58  [↑] Broadcast update for GOOG
2017/02/05 19:51:58  [↓] Received a quote request: jappleseed,SLOW
2017/02/05 19:51:58  [.] New pending quote request
2017/02/05 19:51:58  [-] Waiting for 20 sec
2017/02/05 19:51:58  [↓] Received a quote request: jappleseed,AAPL
2017/02/05 19:51:58  [.] New pending quote request
2017/02/05 19:51:58  [-] Waiting for 2 sec
2017/02/05 19:52:00  [.] Got a response: jappleseed,AAPL,342.61
2017/02/05 19:52:00  [↑] Broadcast update for AAPL
2017/02/05 19:52:18  [.] Got a response: jappleseed,SLOW,977.21
2017/02/05 19:52:18  [↑] Broadcast update for SLOW        <-- requested update
```
