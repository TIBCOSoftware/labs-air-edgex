FROM golang:1.15-alpine AS builder

# add git for go modules
RUN apk update;\
    apk add --no-cache make git gcc libc-dev libsodium-dev zeromq-dev;

WORKDIR /app-service-metadata

COPY go.mod .

RUN go mod download

COPY . .

RUN apk info -a zeromq-dev

RUN make build

# Next image - Copy built Go binary into new workspace
FROM alpine


RUN apk --no-cache add zeromq
COPY --from=builder /app-service-metadata/cmd/res /res
COPY --from=builder /app-service-metadata/cmd/app-service-metadata /app-service-metadata

CMD [ "/app-service-metadata" ,"--registry","--confdir=/res"]