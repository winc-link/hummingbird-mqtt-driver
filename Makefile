.PHONY: build

build:
	docker buildx build --platform linux/amd64 -t 'registry.cn-shanghai.aliyuncs.com/winc-driver/mqtt-driver:2.7' -f docker/Dockerfile . --push


