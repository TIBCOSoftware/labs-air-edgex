FROM golang:1.15-alpine AS builder

# add git for go modules
RUN apk update;\
    apk add --no-cache make git gcc libc-dev libsodium-dev zeromq-dev;

WORKDIR /app-service-zmq-to-mqtt

COPY go.mod .

RUN go mod download

COPY . .

RUN apk info -a zeromq-dev

RUN make build

# Next image - Copy built Go binary into new workspace
FROM alpine


RUN apk --no-cache add zeromq
COPY --from=builder /app-service-zmq-to-mqtt/cmd/res /res
COPY --from=builder /app-service-zmq-to-mqtt/cmd/app-service-zmq-to-mqtt /app-service-zmq-to-mqtt

CMD [ "/app-service-zmq-to-mqtt" ,"--registry","--confdir=/res"]