.PHONY: run build docker.run docker.up docker.down

run:
	go run ./main.go

build:	## Build backend Docker image
	docker build . \
		-t website \
		--no-cache \
		--platform linux/amd64 \

docker.run:
	docker run -d \
	-p 80:80 \
	--name website website

docker.up:
	docker container start website

docker.down:
	docker container stop website