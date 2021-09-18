FROM golang:1.16.7 as build

COPY go build.sh /go/src/github.com/SongCastle/ggnb
RUN cd /go/src/github.com/SongCastle && \
    mv ggnb/build.sh . && chmod 755 ./build.sh && ./build.sh

FROM golang:1.16.7

WORKDIR /go/src/github.com/SongCastle/ggnb

# https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/go-image.html
ADD https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie /usr/bin/aws-lambda-rie
COPY entry.sh build.sh /
COPY --from=build /bin/ggnb /bin/ggnb
RUN chmod 755 /usr/bin/aws-lambda-rie /entry.sh /build.sh && \
    mv /build.sh ../

CMD [ "/bin/ggnb" ]
ENTRYPOINT [ "/entry.sh" ]
