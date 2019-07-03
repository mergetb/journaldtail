.PHONY: journaldtail

# Set HUBUSER to build an image that you can push to a registry
HUBUSER ?= local

JOURNALDSRC = pkg/journald/read.go
STORAGESRC = pkg/storage/memory.go

cmd/journaldtail: cmd/main.go $(JOURNALDSRC) $(STORAGESRC)
	go build -o cmd/journaldtail cmd/main.go

# builds inside the Docker container, then make a runtime image
docker:
	docker build -t $(HUBUSER)/journaldtail .

.PHONY: clean
clean:
	rm -f cmd/journaldtail
