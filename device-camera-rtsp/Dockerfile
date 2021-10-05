FROM gocv/opencv:4.5.0 AS builder


ENV GOPATH /go
WORKDIR $GOPATH

RUN go get -u -d gocv.io/x/gocv


WORKDIR /device-camera-rtsp

COPY . .

RUN ls

# To run tests in the build container:
#   docker build --build-arg 'MAKE=build test' .
# This is handy of you do your Docker business on a Mac
ARG MAKE=build
RUN make $MAKE


ENV APP_PORT=49520
EXPOSE $APP_PORT

WORKDIR /device-camera-rtsp/cmd

RUN ls



# COPY --from=builder /device-camera-rtsp/cmd /

ENTRYPOINT ["./device-camera-rtsp","--confdir=./res"]
