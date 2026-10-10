run:
	go run ./cmd/tollgate

test:
	go test -race ./...

bench:
	go test -bench=. -benchmem ./internal/limiter

cluster:
	docker compose up --build
