FROM golang:1.16.7

WORKDIR /go/src/github.com/SongCastle/ggnb
COPY build.sh ../
RUN chmod 744 ../build.sh
