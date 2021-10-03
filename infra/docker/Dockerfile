### Build Images
FROM golang:alpine AS build-go
WORKDIR /build
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg

FROM build-go AS build-core
RUN go build -o out/core ./cmd/core/*.go

FROM build-go AS build-provider-cloudfoundry
RUN go build -o out/core ./cmd/provider/cloudfoundry/*.go

FROM node:16-alpine AS build-frontend
WORKDIR /build
COPY ./frontend/package.json ./frontend/yarn.lock /build/
RUN yarn install --frozen-lockfile
COPY ./frontend/src /build/src
COPY ./frontend/public /build/public
RUN yarn build

### Final Images
FROM alpine AS core
WORKDIR /app
COPY --from=build-core /build/out/core /app/core
ENTRYPOINT ["/app/core"]

FROM alpine AS provider-cloudfoundry
WORKDIR /app
COPY --from=build-provider-cloudfoundry /build/out/core /app/provider-cloudfoundry
ENTRYPOINT ["/app/provider-cloudfoundry"]

FROM nginx:alpine AS frontend
WORKDIR /usr/share/nginx/html
COPY ./frontend/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build-frontend /build/build /usr/share/nginx/html