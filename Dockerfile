# Step 1: Modules caching
FROM golang:latest as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM --platform=$BUILDPLATFORM golang:latest AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
ARG TARGETPLATFORM
ENV CGO_ENABLED=0
WORKDIR /app
RUN go build -o /bin/app .

# GOPATH for scratch images is /
FROM scratch
COPY --from=builder /bin/app /app
COPY ./templates /templates
EXPOSE 10000
CMD ["/app"]
