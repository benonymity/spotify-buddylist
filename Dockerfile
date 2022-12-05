FROM alpine:latest AS build
RUN apk update
RUN apk upgrade
RUN apk add --update go	gcc g++
WORKDIR /app
COPY api/ .
RUN CGO_ENABLED=1 GOOS=linux go build server.go

FROM alpine:latest

COPY --from=build /app/server /
COPY frontend/dist /dist/

LABEL MAINTAINER="benonymity"
LABEL Github="https://github.com/benonymity/spotify-buddylist"
LABEL version="1.0"
LABEL description="A Docker container to provide a web version of Spotify's friend activity"

CMD ["bin/sh", "-c", "/server"]
