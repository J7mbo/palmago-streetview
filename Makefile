include docker/docker.env
export $(shell sed 's/=.*//' docker/docker.env)


SHELL := /bin/bash

#
# Injector gen binary location
#
gen_bin := ${GOPATH}/src/github.com/j7mbo/goij/bin/gen

#
# Container names
#
app_container_name := palmago-app
elasticsearch_container_name := palmago-elasticsearch
kibana_container_name := palmago-kibana
redis_container_name := palmago-redis

#
# Container statuses
#
elasticstack_status = $(shell docker inspect -f {{.State.Running}} $(elasticsearch_container_name) 2> /dev/null)
redis_status = $(shell docker inspect -f {{.State.Running}} $(redis_container_name) 2> /dev/null)

status: ## Shows the application containers status
	@docker ps -s

#
# General, Injector and Utility commands
#
gen: ## Use Goij to generate an src/Registry file for dir src/
	@chmod +x $(gen_bin)
	@ $(gen_bin) -o src/Registry.go -dir src/

genc: ## Use Goij to generate an src/ConfigRegistry file for dir config/
	@chmod +x $(gen_bin)
	@ $(gen_bin) -o src/ConfigRegistry.go -dir config/
	@sed -i.bak s/GetRegistry/GetConfigRegistry/g src/ConfigRegistry.go
	@rm src/ConfigRegistry.go.bak

genv: ## Use Goij to generate an src/VendorRegistry file for dir vendor/
	@chmod +x $(gen_bin)
	@ $(gen_bin) -o src/VendorRegistry.go -dir vendor/
	@sed -i.bak s/GetRegistry/GetVendorRegistry/g src/VendorRegistry.go
	@rm src/VendorRegistry.go.bak

proto: ## Generate the protobuf files for use in the codebase from the definitions in: /api/proto/v1
	@protoc --go_out=plugins=grpc:. ./api/proto/v1/*.proto

redis-cli: ## Access the redis-cli interface directly
	@COMPOSE_IGNORE_ORPHANS=True docker-compose -f docker/redis.yml exec ${redis_container_name} redis-cli -p ${REDIS_PORT}

#
# Build
#
build: build-app ## Build the application so it's ready to be run

build-app:
	@docker-compose -f docker/app.yml build

#
# Run - -appd is daemonised (background), whilst -app is in the foreground so you can see STDOUT / STDERR
#
run: run-elasticstack run-redis run-app ## Runs whole stack with the application in the foreground to see err output

rund: run-elasticstack run-redis run-appd ## Runs the whole stack with the application in the background

run-appd: create-network
	@COMPOSE_IGNORE_ORPHANS=True docker-compose -f docker/app.yml up -d

run-app: create-network
	@COMPOSE_IGNORE_ORPHANS=True docker-compose -f docker/app.yml up

run-elasticstack:
	@COMPOSE_IGNORE_ORPHANS=True docker-compose -f docker/elastic-stack.yml up -d

# --build is required here so the redis port usage is not cached and is used differently each time
run-redis:
	@COMPOSE_IGNORE_ORPHANS=True REDIS_PORT=${REDIS_PORT} docker-compose -f docker/redis.yml up -d --build

#
# Test a GRPC call. Note you need Grpcc installed for this.
#
test-grpc-call:
	@grpcc --proto ./api/proto/v1/service.proto --address=${GRPC_HOST}:${GRPC_PORT} -i

#
# Kill
#
kill: kill-app kill-elasticstack kill-redis ## Kill the docker containers

kill-app: .kill-${app_container_name}

kill-redis: .kill-${redis_container_name}

kill-elasticstack: .kill-$(elasticsearch_container_name) .kill-$(kibana_container_name)

#
# Destroy
#
destroy: destroy-app destroy-elasticstack destroy-redis ## Kill and remove the docker containers

destroy-app: .destroy-${app_container_name}

destroy-redis: .destroy-${redis_container_name}

destroy-elasticstack: .destroy-$(elasticsearch_container_name) .destroy-$(kibana_container_name)

create-network:
	@docker network create palmago-net 2> /dev/null || true

help:
	@echo 'Usage: make [target] ...'
	@echo
	@echo 'targets:'
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

#
# Pattern rules
#

.kill-%:
	@docker kill $* 2> /dev/null || true

.destroy-%:
	@docker kill $* 2> /dev/null || true
	@docker rm $* 2> /dev/null || true

.DEFAULT_GOAL := help
