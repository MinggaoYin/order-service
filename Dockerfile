# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.13.1

# Force the go compiler to use modules
ENV GO111MODULE=on

WORKDIR /go/src/order-service

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

# This is the ‘magic’ step that will download all the dependencies that are specified in
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the go mod download
# command will _ only_ be re-run when the go.mod or go.sum file change
# (or when we add another docker instruction this line)
RUN go mod download

# Copy the local package files to the container's workspace.
COPY . /go/src/order-service

# Build the order-service command inside the container.
RUN go install order-service

# Run the order-service command by default when the container starts.
ENTRYPOINT /go/bin/order-service

# Document that the service listens on port 8080.
EXPOSE 8080