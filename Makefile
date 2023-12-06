run:
	docker compose -f ./deployment/docker-compose.yml up --build -d

stop:
	docker compose -f ./deployment/docker-compose.yml down && \
	docker volume rm todo_postgres