server:
	TURSO_URL="libsql://settlers-dev-victoroliveirab.turso.io" \
	TURSO_TOKEN="eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJhIjoicnciLCJpYXQiOjE3NDAyMjM1NDQsImlkIjoiMjAzN2FjODEtZjIwNS00YmY3LTlmODEtODc2ZmYyY2IxNTg2In0.oktzAxlYtoTodTQojVX4cMmnQjZ9x1IENnP2fwozmVgSi7tuLBTwQwI1u9Q2hto_4OOPYdBAtam9MNh10puoAA" \
 	go run ./cmd/server/main.go
 
test:
	go test -c ./core/tests/settlement_test.go ./core/tests/constants.go -o core.test
	./core.test -test.v

test-coverage:
	go test -tags=test -c -cover -coverpkg=./... -o core.test ./core
	./core.test -test.v -test.coverprofile=coverage.out
	go tool cover -html=coverage.out
