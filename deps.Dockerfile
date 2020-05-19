FROM alpine:edge AS deps
RUN apk update
RUN apk upgrade
RUN apk add --update go gcc g++ openssh git
# Add the module files and download dependencies.
COPY ./go.mod /go/src/app/go.mod
COPY ./go.sum /go/src/app/go.sum
WORKDIR /go/src/app
RUN go mod download

ARG SSH_PRIVATE_KEY

# Set GOPRIVATE
ENV GOPRIVATE=github.com/safe-distance

# Copy the application source code.
WORKDIR /go/src/app
COPY ./auth ./auth
COPY ./common ./common

# # add ssh credentials on build
# RUN mkdir /root/.ssh/
# RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa
# # fix permissions
# RUN chmod 600 /root/.ssh/id_rsa

# # make sure your domain is accepted
# RUN touch /root/.ssh/known_hosts
# RUN chmod 0700 /root/.ssh
# RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

# Set git config so go get uses ssh instead of https
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

# # Download and install imports
RUN go get -v ./...