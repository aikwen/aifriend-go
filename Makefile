CONFIG ?= ./config.local.yaml
IMAGE_NAME ?= aifriend-go
REGISTRY ?= ccr.ccs.tencentyun.com/kwen

.PHONY: run
run:
	go run ./cmd/api --config-file=$(CONFIG)


.PHONY: docker-build
docker-build:
	@if [ -z "$(VERSION)" ]; then \
		echo "错误: 请通过 VERSION 参数指定版本号。示例: make docker-build VERSION=v1.0.0"; \
		exit 1; \
	fi
	@if [ -z "$$(docker images -q $(IMAGE_NAME):$(VERSION) 2> /dev/null)" ]; then \
		echo "本地未找到镜像 $(IMAGE_NAME):$(VERSION)，开始构建..."; \
		docker build -t $(IMAGE_NAME):$(VERSION) .; \
	else \
		echo "镜像 $(IMAGE_NAME):$(VERSION) 已存在，跳过构建。"; \
	fi
	@echo "开始为远程仓库打 Tag..."
	docker tag $(IMAGE_NAME):$(VERSION) $(REGISTRY)/$(IMAGE_NAME):$(VERSION)
	@echo "tag 完成! 现在可以直接推送: docker push $(REGISTRY)/$(IMAGE_NAME):$(VERSION)"


.PHONY: build-sync
build-sync:
	@echo "正在为 Linux (amd64) 编译数据同步脚本..."
	GOOS=linux GOARCH=amd64 go build -o sync_meili cmd/sync_meili/main.go
	@echo "编译完成! 请将当前目录下的 sync_meili 文件上传到云服务器。"