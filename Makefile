up:
	docker-compose up -d
down:
	docker-compose down
restart:
	docker-compose down && docker-compose up -d
recreate:
	docker-compose up -d --build
ps:
	docker-compose ps
logs:
	docker-compose logs
sh:
	docker-compose exec calc-service sh
generate:
	docker-compose exec calc-service sh -c "cd src/grpc-service && protoc greet/greetpb/greet.proto --go_out=plugins=grpc:."
	docker-compose exec calc-service sh -c "cd src/grpc-service && protoc calculator/calcpb/calc.proto --go_out=plugins=grpc:."