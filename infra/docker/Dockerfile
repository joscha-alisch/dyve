FROM golang:alpine AS build
WORKDIR /build
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg

FROM build AS build-core
RUN go build -o out/core ./cmd/core/*.go

FROM build AS build-provider-cloudfoundry
RUN go build -o out/core ./cmd/provider/cloudfoundry/*.go

FROM alpine AS core
COPY --from=build-core /build/out/core /app

FROM alpine AS provider-cloudfoundry
COPY --from=build-provider-cloudfoundry /build/out/core /app