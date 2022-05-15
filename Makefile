env-init:
	cp -n ./.env.example ./.env || true

run-sender:
	go run ./cmd/sender/main.go