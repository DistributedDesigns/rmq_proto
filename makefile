IMG_NAME = rmq-proto
AMQP_PORT = 44430
MGMT_PORT = 8080

.PHONY: run start stop clean tail

run:
	@docker run \
		--hostname ${IMG_NAME} \
		--name ${IMG_NAME} \
		--publish ${AMQP_PORT}:5672 \
		--publish ${MGMT_PORT}:15672 \
		--detach \
		rabbitmq:3.6.6-management

	@echo "${IMG_NAME} is running on `docker inspect --format '{{ .NetworkSettings.IPAddress }}' ${IMG_NAME}`"

	@echo "AMPQ: ${AMQP_PORT}"

	@echo "Management GUI: ${MGMT_PORT}"

start:
	docker start ${IMG_NAME}

stop:
	docker stop ${IMG_NAME}

clean: stop
	docker rm ${IMG_NAME}

tail:
	@docker logs -f ${IMG_NAME}
