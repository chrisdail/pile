FROM ubuntu:18.04

RUN apt-get update && apt-get install -y curl jq

ENV GO_VERSION 1.13.4
RUN curl "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz" -o go.tar.gz \
  && tar -C /usr/local -xzf go.tar.gz \
  && rm go.tar.gz \
  && ln -s /usr/local/go/bin/go /usr/bin/go

ADD entrypoint.sh /entrypoint.sh
ADD build.sh /build.sh
ENTRYPOINT ["/entrypoint.sh"]
