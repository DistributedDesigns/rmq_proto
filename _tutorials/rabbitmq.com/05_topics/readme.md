Topics allow subscriptions based on pattern matching. Typically, patterns are `<topic1>.<topic2>.<topic3>` and so on. `.` delimits topics. The pattern must be <= 255 characters.

There are special characters:
- `*` matches exactly one topic
- `#` matches zero or more topics

```shell
# Topics will have the pattern <app>.<severity>
# Receive all logs, like a fanout
go run receive_log_topic.go "#"

# Receive everything from "kern" and all critical logs
go run receive_log_topic.go "kern.*" "*.critical"

# Send some message!
# these will be captured by the second logger
go run emit_log_topic.go "kern.info" "What is going on?"
go run emit_log_topic.go "sh.critical" "Holy shit!"

# Only one message is displayed even though it matches
# both patterns for the second logger.
go run emit_log_topic.go "kern.critical" "Piiiiissss"

# Only hits the universal logger
go run emit_log_topic.go "cron.fatal" "Shit is crazy"
go run emit_log_topic.go "." "It's so hot in here."

```

[link](https://www.youtube.com/watch?v=vUz9xCTOPRw)
