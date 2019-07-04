# Set HUBUSER to build an image that you can push to a registry
HUBUSER ?= local
distro ?= stretch

JOURNALDSRC = pkg/journald/read.go
STORAGESRC = pkg/storage/memory.go

cmd/journaldtail.$(distro): cmd/main.go $(JOURNALDSRC) $(STORAGESRC)
	go build -o $@ cmd/main.go

jessie:
	docker build $(DOCKER_ARGS) -f jessie.dock -t $(HUBUSER)/journaldtail-jessie .
	docker create -ti --name jdt $(HUBUSER)/journaldtail-jessie
	docker cp jdt:/usr/bin/journaldtail cmd/journaldtail.jessie
	docker rm -fv jdt

stretch:
	docker build $(DOCKER_ARGS) -f stretch.dock -t $(HUBUSER)/journaldtail-stretch .
	docker create -ti --name jdt $(HUBUSER)/journaldtail-stretch
	docker cp jdt:/usr/bin/journaldtail cmd/journaldtail.stretch
	docker rm -fv jdt

buster:
	docker build $(DOCKER_ARGS) -f buster.dock -t $(HUBUSER)/journaldtail-buster .
	docker create -ti --name jdt $(HUBUSER)/journaldtail-buster
	docker cp jdt:/usr/bin/journaldtail cmd/journaldtail.buster
	docker rm -fv jdt

.PHONY: clean
clean:
	rm -f cmd/journaldtail.$(distro)
