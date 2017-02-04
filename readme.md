Rabbit MQ Prototype
====

How is rabbit formed?

## Setup
Start a local rabbit server. This requires a running instance of Docker.

```shell
make run
```

This will echo the ports used for AMQP and the management GUI. Access them from `localhost:$PORT`. Use the default login for the GUI: `guest` / `guest`.
