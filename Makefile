#/bin/bash
NAME = $(shell basename $(shell pwd -P))
VERSION ?= $(BUILDER)-$(shell git rev-parse --short HEAD)
TARGET = main.go
GOPRIVATE="sonymimic1"
OUTPUT = $(NAME)
DOCKER_IMAGE = sonymimic/checkrtpapi
TAG = 1.0
clean-cache:
	go clean --modcache

tidy:
	GOPRIVATE=$(GOPRIVATE) go mod tidy

build:
	make tidy
	GOPRIVATE=$(GOPRIVATE) go build -o $(OUTPUT) ./cmd/$(TARGET)


# 檢測操作系統類型
UNAME_S := $(shell uname -s)

# 根據不同的操作系統設置獲取 IP 的命令
ifeq ($(UNAME_S),Linux)
    # Ubuntu (Linux) 系統
    IP := $(shell hostname -I | awk '{print $$1}')
else ifeq ($(UNAME_S),Darwin)
    # macOS 系統
    IP := $(shell ipconfig getifaddr en0 || ipconfig getifaddr en1)
else
    $(error Unsupported operating system)
endif

# 顯示 IP 的目標
show-ip:
	@echo $(IP)


docker-build:
	docker build -t $(DOCKER_IMAGE):$(TAG) .

docker-push:
	docker push $(DOCKER_IMAGE):$(TAG)

service-up:
	rm -rf ./shared-config/redis-cluster/redis
	sh ./shared-config/redis-cluster/redis-config.sh $(IP)
	docker-compose up -d
	docker exec -it redis-1 redis-cli --cluster create $(IP):6371 $(IP):6372 $(IP):6373 $(IP):6374 $(IP):6375 $(IP):6376 --cluster-replicas 1 --cluster-yes
	sleep 2
	docker exec -it redis-1 redis-cli -h redis-1 -p 6371 -c cluster nodes
	docker exec -it redis-1 redis-cli -h redis-1 -p 6371 -c set AT01-BET 100
	docker exec -it redis-1 redis-cli -h redis-1 -p 6371 -c set AT01-WIN 97 
	docker exec -it redis-1 redis-cli -h redis-1 -p 6371 -c set AB3-BET 100 
	docker exec -it redis-1 redis-cli -h redis-1 -p 6371 -c set AB3-WIN 96

service-down:
	docker-compose down
	docker volume prune -f
	docker network prune -f


clean:
	rm -f $(OUTPUT)
	rm -f $(OUTPUT).mac
	rm -f $(OUTPUT).linux
	rm -f *.txt
	rm -f *.log
	rm -f *.test
	rm -rf ./log/
	rm -f *.pid
 
.PHONY: build mac linux clean clean-cache up go-mod-name service-up service-down show-ip tidy build docker-build docker-push