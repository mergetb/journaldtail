
# NOTE: Example only. The Systemd c libraries used to build and run
#       Journaldtail must match the host version. So you may need to
#       make build your own image.

FROM ubuntu:18.04 as builder

RUN apt-get update \
            && apt-get install -yq libsystemd-dev make wget git build-essential \
            && rm -rf /var/lib/apt/lists/*

ENV GOVERSION 1.12.1

RUN cd /opt && wget https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz && \
    tar zxf go${GOVERSION}.linux-amd64.tar.gz && rm go${GOVERSION}.linux-amd64.tar.gz && \
    ln -s /opt/go/bin/go /usr/bin/

WORKDIR /root/src/github.com/hikhvar/journaldtail/
COPY . /root/src/github.com/hikhvar/journaldtail/

RUN make clean
RUN make

# I use this image with Docker Swarm, and during development run:
#
# docker run -ti -v /etc/machine-id:/etc/machine-id \
#       -v /var/run/systemd/journal/:/var/run/systemd/journal/ \
#       -v /var/log/journal:/run/log/journal \
#       --network loki_default \
#       -e LOKI_URL=http://loki:3100/api/prom/push \
#           svendowideit/journaldtail

FROM ubuntu:18.04

COPY --from=builder /root/src/github.com/hikhvar/journaldtail/cmd/journaldtail /usr/bin
ENTRYPOINT ["/usr/bin/journaldtail"]
