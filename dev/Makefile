.PHONY:
redis:
	docker run --rm --name NAME -p P:6379 -p P:8001 redis/redis-stack:latest

rediscli:
	docker exec -it NAME redis-cli
