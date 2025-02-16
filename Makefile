server:
	go run ./cmd/server/main.go

test:
	go test -c ./core/state/tests/... -o core.test
	./core.test -test.v
