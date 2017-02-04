Multiple workers consume tasks from a queue. Tasks have variable length corresponding to the number of `...` after the message.

Queues are _durable_ and will be retained between `docker restart` cycles. Removing an image and re-running it will flush everything.

```shell
go run worker.go

# in another shell
go run worker.go

# in yet another shell
go run new_task.go one.  # 1 sec task
go run new_task.go two.. # 2 sec task...
go run new_task.go three...
go run new_task.go four....
go run new_task.go five.....
```
