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

[![image](https://user-images.githubusercontent.com/62207008/217578966-c1daa0b2-5040-4906-abe8-aa7a2f276956.png)](https://ip.minpeter.xyz)

## how to use?

```sh
$ curl ip.minpeter.xyz -L
```

or <https://ip.minpeter.xyz>

## deployment

with docker

```
$ docker build -t iplogger .
$ docker run -dp 10000:10000 iplogger
```

or pre-built image

```
$ docker run -dp 10000:10000 ghcr.io/minpeter/iplogger:latest
```

with golang

```
$ go mod tidy
$ go run .
```

now running on <http://localhost:10000>

## ✨ result post (Korean) ✨

프로젝트에 대한 간단한 [설명 글](docs/result.md)  
A brief [description](docs/result.md) of the project

이 프로젝트를 진행하면서 작성한 [블로그 글](https://minpeter.github.io/uncategorized/%EB%B0%B1%EC%97%94%EB%93%9C%EC%97%90%EC%84%9C-client%EC%9D%98-ip%EB%A5%BC-%EB%A1%9C%EA%B9%85%ED%95%98%EB%8A%94-%EB%B0%A9%EB%B2%95)  
[Blog post](https://minpeter.github.io/uncategorized/%EB%B0%B1%EC%97%94%EB%93%9C%EC%97%90%EC%84%9C-client%EC%9D%98-ip%EB%A5%BC-%EB%A1%9C%EA%B9%85%ED%95%98%EB%8A%94-%EB%B0%A9%EB%B2%95) written while working on this project
