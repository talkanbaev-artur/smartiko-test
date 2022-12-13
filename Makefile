IMAGE_NAME = ehdw/smartiko-test
VERSION = 0.1

include .env
export

run: 
	go run src/cmd/*.go

build: mod_tidy create_docker tag_latest

run_d: mod_tidy build run_docker

create_docker:
	docker build --tag $(IMAGE_NAME):$(VERSION) .

tag_latest:
	docker image tag $(IMAGE_NAME):$(VERSION) $(IMAGE_NAME):latest

mod_tidy:
	go mod tidy

run_docker:
	docker run --name msu-orders-loader --env-file .env.docker --network=host -d $(IMAGE_NAME):$(VERSION)
