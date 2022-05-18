.PHONY:
go-gen-deps:
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/tinylib/msgp@v1.1.6