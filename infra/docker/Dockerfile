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
RUN go build -o out/cmd ./cmd/core/*.go

FROM build-go AS build-provider-cloudfoundry
RUN go build -o out/cmd ./cmd/provider/cloudfoundry/*.go

FROM build-go AS build-provider-demo
RUN go build -o out/cmd ./cmd/provider/demo/*.go

FROM build-go AS build-provider-github
RUN go build -o out/cmd ./cmd/provider/github/*.go

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
COPY infra/docker/defaults/core.yaml /app/config.yaml
COPY --from=build-core /build/out/cmd /app/core
ENTRYPOINT ["/app/core"]

FROM alpine AS provider-cloudfoundry
WORKDIR /app
COPY --from=build-provider-cloudfoundry /build/out/cmd /app/provider-cloudfoundry
ENTRYPOINT ["/app/provider-cloudfoundry"]

FROM alpine AS provider-demo
WORKDIR /app
COPY --from=build-provider-demo /build/out/cmd /app/provider-demo
ENTRYPOINT ["/app/provider-demo"]

FROM alpine AS provider-github
WORKDIR /app
COPY --from=build-provider-github /build/out/cmd /app/provider-github
ENTRYPOINT ["/app/provider-github"]

FROM nginx:alpine AS frontend
WORKDIR /usr/share/nginx/html
COPY ./frontend/nginx.conf.template /etc/nginx/templates/default.conf.template
COPY --from=build-frontend /build/build /usr/share/nginx/html
