#FROM golang:alpine AS builder
FROM golang:1.12 AS builder
# Add all the source code (except what's ignored
# under `.dockerignore`) to the build context.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group
COPY ./*.go ./
ADD ./config.yml ./
RUN go get -u github.com/Azure/azure-event-hubs-go/...
RUN go get -u github.com/Azure/azure-amqp-common-go/...
RUN go get -u github.com/Azure/go-autorest/...
RUN go get -u gopkg.in/yaml.v2/...
RUN pwd
RUN ls -lh
ENV GOPATH=/go
ENV GOBIN=/go/bin/
ENV CGO_ENABLED=0
RUN set -ex 
RUN go install -v
RUN ls /go/bin
FROM golang:alpine AS final
COPY --from=builder /go/bin/go /usr/bin/appdynamics-azure-event-hub
COPY --from=builder /go/config.yml ./
RUN chmod 777 /usr/bin/appdynamics-azure-event-hub
# Set the binary as the entrypoint of the container
ENTRYPOINT [ "/usr/bin/./appdynamics-azure-event-hub" ]
