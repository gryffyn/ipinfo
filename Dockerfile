# Build
FROM golang:1.21-alpine AS build
WORKDIR /go/src/github.com/gryffyn/ipinfo
COPY . .

# Must build without cgo because libc is unavailable in runtime image
ENV GO111MODULE=on CGO_ENABLED=0
RUN make

# Run
FROM scratch
EXPOSE 8080

COPY --from=build /go/bin/ipinfo /opt/ipinfo/
COPY data /opt/ipinfo/data

WORKDIR /opt/ipinfo
ENTRYPOINT ["/opt/ipinfo/ipinfo"]
