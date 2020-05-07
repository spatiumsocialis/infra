FROM golang:1.13

ARG SSH_PRIVATE_KEY

# Set default port
ENV PORT 8080
# Set default DB vars
ENV DB_PROVIDER sqlite3
ENV DB_CONNECTION_STRING :memory:

# Create dir where service account json will go
RUN mkdir /root/google

# add ssh credentials on build
RUN mkdir /root/.ssh/
RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa
# fix permissions
RUN chmod 600 /root/.ssh/id_rsa

# make sure your domain is accepted
RUN touch /root/.ssh/known_hosts
RUN chmod 0700 /root/.ssh
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts


WORKDIR /go/src/app
COPY . .

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

RUN go get -d -v -insecure ./...
RUN go install -v ./...

RUN chmod 777 entrypoint.sh

ENTRYPOINT [ "/go/src/app/entrypoint.sh" ]