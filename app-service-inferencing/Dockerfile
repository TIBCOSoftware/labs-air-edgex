FROM golang:1.15-alpine AS builder

# add git for go modules
RUN apk update;\
    apk add --no-cache make git gcc libc-dev libsodium-dev zeromq-dev;

WORKDIR /app-service-inferencing

COPY go.mod .

RUN go mod download

COPY . .

RUN apk info -a zeromq

RUN make build

# Next image - Copy built Go binary into new workspace
FROM alpine


RUN apk --no-cache add zeromq

RUN apk info -a zeromq

COPY --from=builder /app-service-inferencing/cmd/res /res
COPY --from=builder /app-service-inferencing/cmd/app-service-inferencing /app-service-inferencing

CMD [ "/app-service-inferencing" ,"--registry","--confdir=/res"]