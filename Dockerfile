# Build Stage
FROM lacion/alpine-golang-buildimage:1.12.4 AS build-stage

LABEL app="build-bugzillagithubsync"
LABEL REPO="https://github.com/alexander-demichev/bugzillagithubsync"

ENV PROJPATH=/go/src/github.com/alexander-demichev/bugzillagithubsync

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/alexander-demichev/bugzillagithubsync
WORKDIR /go/src/github.com/alexander-demichev/bugzillagithubsync

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/alexander-demichev/bugzillagithubsync"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/bugzillagithubsync/bin

WORKDIR /opt/bugzillagithubsync/bin

COPY --from=build-stage /go/src/github.com/alexander-demichev/bugzillagithubsync/bin/bugzillagithubsync /opt/bugzillagithubsync/bin/
RUN chmod +x /opt/bugzillagithubsync/bin/bugzillagithubsync

# Create appuser
RUN adduser -D -g '' bugzillagithubsync
USER bugzillagithubsync

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/bugzillagithubsync/bin/bugzillagithubsync"]
