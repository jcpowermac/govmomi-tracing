FROM quay.io/fedora/fedora:latest AS builder

#ARG BRANCH=master
#ARG REPO=https://github.com/openshift/installer
#ARG DELVE_BRANCH=master
#ARG DELVE_REPO=https://github.com/go-delve/delve
ENV WORKDIR_PATH="/src"
#ENV DELVE_WORKDIR_PATH="/go/src/github.com/go-delve/delve/"

RUN mkdir -p $WORKDIR_PATH

WORKDIR $WORKDIR_PATH
COPY . .
RUN dnf install golang git zip -y && \
    go build main.go

FROM quay.io/fedora/fedora:latest AS run
COPY --from=builder /src/main /bin/govmomi-trace

CMD ["/bin/govmomi-trace"]
