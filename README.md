# iplogger

[![Go Report Card](https://goreportcard.com/badge/github.com/minpeter/iplogger)](https://goreportcard.com/report/github.com/minpeter/iplogger)
[![Go Reference](https://pkg.go.dev/badge/github.com/minpeter/iplogger.svg)](https://pkg.go.dev/github.com/minpeter/iplogger)
[![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/minpeter/iplogger)](https://hub.docker.com/r/minpeter/iplogger)
[![Docker Pulls](https://img.shields.io/docker/pulls/minpeter/iplogger)](https://hub.docker.com/r/minpeter/iplogger)  
👀 Project: What is my IP?

## purpose

A study on how services located behind multiple reverse proxies log real client IPs.

여러 리버스 프록시 뒤에 위치한 서비스가 실제 클라이언트 IP를 기록하는 방법에 대한 연구

## screenshot

[![image](https://user-images.githubusercontent.com/62207008/217578966-c1daa0b2-5040-4906-abe8-aa7a2f276956.png)](https://ip.minpeter.uk)

## how to use?

```sh
curl ip.minpeter.uk -L
```

or <https://ip.minpeter.uk>

## deployment

with docker

```sh
docker build -t iplogger .
docker run -dp 10000:10000 iplogger
```

or pre-built image

```sh
docker run -dp 10000:10000 ghcr.io/minpeter/iplogger:latest
```

with golang

```sh
go mod tidy
go run .
```

now running on <http://localhost:10000>

## ✨ result post ✨

A brief [description](docs/result.md) of the project (Korean only)  
[Blog post](https://minpeter.uk/blog/how-loggin-real-ip) written while working on this project (Korean only)

프로젝트에 대한 간단한 [설명 글](docs/result.md)  
이 프로젝트를 진행하면서 작성한 [블로그 글](https://minpeter.uk/blog/how-loggin-real-ip)
