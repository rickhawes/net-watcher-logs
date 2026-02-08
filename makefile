SHELL=/bin/zsh

# Initialize the project
init:
	go mod download

# Run a local version 
run:
	go run . 

# Run using docker compose
run-docker:
	docker compose -f compose.yml up

# Make docker images
#  requires docker to be running. (eg. orb)
#
build:
	./autotag.sh
	docker build --tag rickhawes/net-watcher-logs:amd-latest --platform linux/amd64 .
	docker build --tag rickhawes/net-watcher-logs:latest .

# Push docker images
#
push:
	docker login
	docker push rickhawes/net-watcher-logs:amd-latest
	docker push rickhawes/net-watcher-logs:latest

.PHONY: init gather_requirements build push