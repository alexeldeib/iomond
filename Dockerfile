# Build the manager binary
FROM golang:1.13 as builder

WORKDIR /workspace

RUN mkdir tools \
    && cd tools \
    && go mod init tools \
    && go get golang.org/x/tools/cmd/goimports

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . .

RUN goimports -w .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o iomond main.go

# Use distroless as minimal base image to package the iomond binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# FROM gcr.io/distroless/base-debian10:nonroot
FROM ubuntu:bionic
WORKDIR /

RUN apt update && apt install -y curl jq fio sysstat nano

RUN echo "[global]\n\
name=fio-rand-RW\n\
ioengine=libaio\n\
rw=randrw\n\
rwmixread=100\n\
bs=4K\n\
direct=1\n\
time_based=1\n\
runtime=60\n\
iodepth=64\n\
size=4G\n\
numjobs=2\n\
\n\
[sda]\n\
filename=/dev/sda\n\
" > job

COPY --from=builder /workspace/iomond .

ENTRYPOINT ["/iomond"]