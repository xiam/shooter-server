DOCKER_IMAGE          ?= xiam/shooter-server
DOCKER_CONTAINER      ?= shooter-server

EXTERNAL_PORT         ?= 3223
CONTAINER_PORT        ?= 3223

DEPLOY_TARGET         ?= prod

POSTGRESQL_USER       ?=
POSTGRESQL_PASSWORD   ?=

all: vendor-sync build

run: build
	./shooter-server

build:
	go build -o shooter-server

vendor-sync:
	govendor sync

clean:
	rm -f shooter-server && \
	rm -f shooter_*

shooter_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o shooter_linux_amd64

docker-build: shooter_linux_amd64
	docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build
	(docker stop $(DOCKER_CONTAINER) | exit 0) && \
	(docker rm $(DOCKER_CONTAINER) | exit 0) && \
	docker run -d \
		--restart=always \
		--link postgresql \
		-e POSTGRESQL_USER="$(POSTGRESQL_USER)" \
		-e POSTGRESQL_PASSWORD="$(POSTGRESQL_PASSWORD)" \
		--name=$(DOCKER_CONTAINER) \
		-p 127.0.0.1:$(EXTERNAL_PORT):$(CONTAINER_PORT) \
		$(DOCKER_IMAGE)

deploy: clean shooter_linux_amd64
	ansible-playbook playbook.yml -e "host=$(DEPLOY_TARGET)" -e "dbuser=$(POSTGRESQL_USER)" -e "dbpass=$(POSTGRESQL_PASSWORD)"
