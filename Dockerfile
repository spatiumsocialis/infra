ARG EXECUTABLE
FROM dependencies AS builder
ARG SERVICE
ARG EXECUTABLE
ARG SSH_PRIVATE_KEY

# Set GOPRIVATE
ENV GOPRIVATE=github.com/safe-distance

# Copy the application source code.
WORKDIR /go/src/app
COPY ./${SERVICE} ./${SERVICE}

# add ssh credentials on build
RUN mkdir /root/.ssh/
RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa
# fix permissions
RUN chmod 600 /root/.ssh/id_rsa

# make sure your domain is accepted
RUN touch /root/.ssh/known_hosts
RUN chmod 0700 /root/.ssh
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

# Set git config so go get uses ssh instead of https
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

# Download and install imports
RUN go get -v  -insecure ./...

# Build the application.
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /go/bin/${EXECUTABLE} /go/src/app/${SERVICE}/cmd/${EXECUTABLE}

ENTRYPOINT /go/bin/${EXECUTABLE}

FROM alpine:latest
ARG EXECUTABLE
ENV EXECUTABLE=${EXECUTABLE}
# Copy executable to /bin/
COPY --from=builder /go/bin/${EXECUTABLE} /bin/${EXECUTABLE}

# Set env defaults
ENV DB_PROVIDER sqlite3
ENV DB_CONNECTION_STRING :memory:
ENV PORT 8080

ENTRYPOINT /bin/${EXECUTABLE}