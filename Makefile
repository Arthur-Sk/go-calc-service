up:
	docker-compose up -d
down:
	docker-compose down
recreate:
	docker-compose up -d --build
ps:
	docker-compose ps
sh:
	docker-compose exec go-grpc sh