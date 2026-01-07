.PHONY: build

build:
	docker buildx build --platform linux/amd64 -t 'registry.cn-shanghai.aliyuncs.com/winc-driver/mqtt-driver:3.0' -f docker/Dockerfile . --push


