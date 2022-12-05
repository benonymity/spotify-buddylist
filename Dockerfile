FROM alpine:latest AS build
RUN apk update
RUN apk upgrade
RUN apk add --update go=1.8.3-r0 gcc=6.3.0-r4 g++=6.3.0-r4
WORKDIR /app
COPY api/server.go .
ENV GOPATH /app
RUN CGO_ENABLED=1 GOOS=linux go build server.go

FROM alpine:latest

COPY --from=build /app/server /
COPY frontend/dist /dist/

LABEL MAINTAINER="benonymity"
LABEL Github="https://github.com/benonymity/spotify-buddylist"
LABEL version="1.0"
LABEL description="A Docker container to provide a web version of Spotify's friend activity"

CMD ["bin/sh", "-c", "/server"]
