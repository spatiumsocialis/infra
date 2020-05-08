FROM alpine:edge AS dep
RUN apk update
RUN apk upgrade
RUN apk add --update go gcc g++ openssh git
# Add the module files and download dependencies.
COPY ./go.mod /go/src/app/go.mod
COPY ./go.sum /go/src/app/go.sum
WORKDIR /go/src/app
RUN go mod download
# Add the shared packages.
COPY ./auth /go/src/app/auth
COPY ./common /go/src/app/common