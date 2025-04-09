FROM golang:1.24 AS builder

LABEL maintainer="marius@stud.ntnu.no"
LABEL stage=builder

# Set up execution environment in container's GOPATH
WORKDIR /go

# Copy go mod and sum to pre-download dependencies before copying rest of files
COPY ./go.mod /go/go.mod
COPY ./go.sum /go/go.sum
RUN go mod download

# Copy relevant folders into container
COPY ./clients /go/clients
COPY ./config /go/config
COPY ./database /go/database
COPY ./handlers /go/handlers
COPY ./services /go/services
COPY ./utils /go/utils
COPY ./main.go /go/main.go


# Compile binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o server

# Indicate port on which server listens
EXPOSE 8080

# Instantiate binary
CMD ["./server"]
