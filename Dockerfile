FROM alpine:latest

COPY activity_api_go/server .

LABEL MAINTAINER="Benonymity"
LABEL Github="https://github.com/benonymity/spotify-buddylist"
LABEL version="1"
LABEL description="A Docker container to provide a web version of Spotify's friend activity"

CMD ./server
