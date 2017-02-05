Receivers can subscribe to the same `direct` exchange but filter based on routing keys. Each receiver must explicitly `.QueueBind()` to the `routingKey` that it cares about.

If you subscribe to a `direct` exchange but don't specify a routing key then you get nothing :(

```shell
# This receiver will subscribe to only errors
go run receive_log_direct.go error

# In another shell, make a receiver that subscribes to all messages
go run receive_log_direct.go info warning error

# Send some message!
go run emit_log_direct.go error "I display on both receivers"
go run emit_log_direct.go info "I only display on one receiver"
```
