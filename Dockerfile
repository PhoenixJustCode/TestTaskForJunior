FROM golang:1.19-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN apt-get update && apt-get -y install postgresql-client
RUN chmod +x wait-for-postgres.sh

# build go app
RUN go mod download
RUN go build -o test-task-app ./cmd/main.go

CMD ["./test-task-app"]