############################
# STEP 1 build executable binary
############################

ARG GO_VERSION=1.20.6

FROM golang:${GO_VERSION}-alpine AS builder_base

# Install git.  - TODO Not sure if this is needed if we using Vendor
# Git is required for fetching the dependencies.
RUN apk update \
    && apk add --no-cache git \
    && apk add --no-cache make

RUN mkdir -p /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

FROM builder_base AS builder

# Copy sources
COPY . /app

ENV CGO_ENABLED=0

# Build the binary.  Make the world...
RUN go build

############################
# STEP 2 build a small image
############################
FROM alpine:latest

# Copy our executable.
COPY --from=builder /app/ittybittyfeaturechecker /ittybittyfeaturechecker
COPY features.json /

EXPOSE 8081

# Run the hello binary.
ENTRYPOINT ["/ittybittyfeaturechecker"]
 