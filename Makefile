.PHONY: full-run

full-run: test run

run:
	docker compose -f ./deployment/docker-compose.yml up --build -d

stop:
	docker compose -f ./deployment/docker-compose.yml down && \
	docker volume rm todo_postgres

test:
	go test -race ./internal/...