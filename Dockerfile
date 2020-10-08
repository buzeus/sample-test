FROM golang AS builder
USER root
# COPY main.go /tmp/
COPY ./ /go/src/
# To make it run with Alpine - https://stackoverflow.com/questions/36279253/go-compiled-binary-wont-run-in-an-alpine-docker-container-on-ubuntu-host
# Final Image size ~ 20 MB
RUN cd /go/src/ && \
go build -tags netgo -a main.go

# To run with ubuntu (final Image size ~ 87 MB)
# RUN cd /tmp && \
# go build main.go


FROM alpine
USER root
COPY --from=builder /go/src/main /tmp/
EXPOSE 8080
RUN cd /tmp && \
mkdir /.cache && \
chmod -R 777 /tmp /.cache
ENTRYPOINT ["/tmp/main"]
USER 1001