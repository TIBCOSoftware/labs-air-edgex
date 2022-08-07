#
# Copyright (c) 2019 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# This file will work as is for local development. No need to use Dockerfile.build

#build stage
ARG ALPINE=golang:1.17-alpine3.15
FROM ${ALPINE} AS builder

ARG ALPINE_PKG_BASE="make git openssh-client gcc libc-dev zeromq-dev libsodium-dev"
ARG ALPINE_PKG_EXTRA=""

LABEL license='SPDX-License-Identifier: Apache-2.0' \
  copyright='Copyright (c) 2019: Intel'
RUN sed -e 's/dl-cdn[.]alpinelinux.org/nl.alpinelinux.org/g' -i~ /etc/apk/repositories
RUN apk add --no-cache ${ALPINE_PKG_BASE} ${ALPINE_PKG_EXTRA}

WORKDIR /device-sound-simulator

COPY . .

RUN go mod tidy
RUN go mod download

ARG MAKE=build
RUN make $MAKE


#final stage
FROM alpine:3.12

LABEL license='SPDX-License-Identifier: Apache-2.0' \
  copyright='Copyright (c) 2019: Intel'
LABEL Name=device-generic-rest Version=${VERSION}


RUN apk add --update --no-cache zeromq dumb-init

COPY --from=builder /device-sound-simulator/cmd /

EXPOSE 48985

ENTRYPOINT ["/device-sound-simulator"]
CMD ["--cp=consul://edgex-core-consul:8500", "--confdir=/res", "--registry"]
