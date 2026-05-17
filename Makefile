buildup:
	docker compose -f docker-compose.yaml --env-file .env up --build

up:
	docker compose -f docker-compose.yaml --env-file .env up 