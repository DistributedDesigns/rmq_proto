Messages passed into the `logs` exchange will consumed by every instance of `receive_logs.go`.

```shell
# open in as may different shells as you like
go run receive_logs.go

# all receivers will echo this message
go run emit_log.go Hello friends!
```
