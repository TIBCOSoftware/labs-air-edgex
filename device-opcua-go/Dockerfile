#
# Copyright (c) 2018, 2019 Intel
#
# SPDX-License-Identifier: Apache-2.0
#

ARG ALPINE=golang:1.15-alpine
FROM ${ALPINE} AS builder
ARG ALPINE_PKG_BASE="build-base git openssh-client"
ARG ALPINE_PKG_EXTRA=""

# Replicate the APK repository override.
RUN sed -e 's/dl-cdn[.]alpinelinux.org/mirrors.ustc.edu.cn/g' -i~ /etc/apk/repositories
# Install our build time packages.
RUN apk add --no-cache ${ALPINE_PKG_BASE} ${ALPINE_PKG_EXTRA}

WORKDIR /device-opcua-go

COPY . .

# To run tests in the build container:
#   docker build --build-arg 'MAKE=build test' .
# This is handy of you do your Docker business on a Mac
ARG MAKE=build
RUN make $MAKE

# Next image - Copy built Go binary into new workspace
FROM scratch

# expose command data port
ENV APP_PORT=49997
EXPOSE $APP_PORT

COPY --from=builder /device-opcua-go/cmd /

ENTRYPOINT ["/device-opcua-go","--confdir=/res"]
