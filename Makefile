.PHONY: run build docker.run docker.up docker.down

run:
	go run ./main.go

build:	## Build backend Docker image
	docker build . \
		-t website \
		--no-cache \

docker.run:
	docker run -d \
	-p 8080:8080 \
	--name website website

docker.up:
	docker container start website

docker.down:
	docker container stop website