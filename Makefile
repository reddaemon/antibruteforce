precommit:
	gofmt -w -s -d .
	go vet .
	golangci-lint run --enable-all
	go mod tidy
	go mod verify
unit-test:
	go test -race -cover ./internal/service/api/bucket/...
integr-test:
	docker-compose -f docker/docker-compose/docker-compose.test.yml up --build -d ;\
		docker-compose -f docker/docker-compose/docker-compose.test.yml run integration_tests go test ./integr_test;\
		test_status_code=$$? ;\
		docker-compose -f docker/docker-compose/docker-compose.test.yml down ;\
		exit $$test_status_code ;\

gen-proto:
	protoc -I. protofiles/antibruteforce.proto --go_out=plugins=grpc:.
run:
	go run -race main.go api

up:
	docker-compose -f docker/docker-compose/docker-compose.yml up
down:
	docker-compose -f docker/docker-compose/docker-compose.yml down

restart: down up