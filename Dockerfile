FROM alpine:latest

COPY api/server /
COPY frontend/dist /dist/

LABEL MAINTAINER="benonymity"
LABEL Github="https://github.com/benonymity/spotify-buddylist"
LABEL version="1.0"
LABEL description="A Docker container to provide a web version of Spotify's friend activity"

CMD ["bin/sh", "-c", "/server"]
