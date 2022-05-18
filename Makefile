.PHONY:
go-gen-deps:
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/tinylib/msgp@v1.1.6

.PHONY:
docker-build:
	DOCKER_BUILDKIT=0 docker build --no-cache -f build/server/Dockerfile -t pow-server .
	DOCKER_BUILDKIT=0 docker build --no-cache -f build/client/Dockerfile -t pow-client .