from golang:1.13
label mantainer="Kevin Aleman"

run groupadd -r kestar && useradd -m --no-log-init -r -g kestar kestar && mkdir /app && chown -R kestar: /app

workdir /app

copy go.mod go.sum ./
run go mod download

copy . .

run go build -o main .

cmd ["./main"]