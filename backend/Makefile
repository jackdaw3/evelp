APP=evelp
.PHONY: clean test run build

clean:
	go clean

test:
	go test ./... -gcflags=-l --count=1 -cover

run:
	go run main.go

build: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${APP}
