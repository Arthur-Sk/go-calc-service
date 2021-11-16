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
generate:
	docker-compose exec go-grpc sh -c "cd src && protoc grpc-service/greet/greetpb/greet.proto --go_out=plugins=grpc:."