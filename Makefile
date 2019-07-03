.PHONY: journaldtail

# Set HUBUSER to build an image that you can push to a registry
HUBUSER ?= local

journaldtail:
	go build -o cmd/journaldtail cmd/main.go

# builds inside the Docker container, then make a runtime image
docker:
	docker build -t $(HUBUSER)/journaldtail .

