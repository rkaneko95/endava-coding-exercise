lint:
	golangci-lint run -v

fmt:
	go fmt ./...

test:
	go clean -testcache && go test -v -cover ./...

run_redis:
	docker pull redis && docker run -d -p 6379:6379 --name redis-container redis