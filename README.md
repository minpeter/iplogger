# ipLogger
[![Go Report Card](https://goreportcard.com/badge/github.com/minpeter/iplogger)](https://goreportcard.com/report/github.com/minpeter/iplogger)
[![Go Reference](https://pkg.go.dev/badge/github.com/minpeter/iplogger.svg)](https://pkg.go.dev/github.com/minpeter/iplogger)  
👀 Project What is my IP?

## purpose

A study on how services located behind multiple reverse proxies log real client IPs.

여러 리버스 프록시 뒤에 위치한 서비스가 실제 클라이언트 IP를 기록하는 방법에 대한 연구

## screenshot

[![image](https://user-images.githubusercontent.com/62207008/217578966-c1daa0b2-5040-4906-abe8-aa7a2f276956.png)](https://ip.minpeter.cf)

## how to use?

```sh
$ curl ip.minpeter.cf -L
```
or <https://ip.minpeter.cf>

## deployment

with docker
```
$ docker build -t iplogger .
$ docker run -dp 10000:10000 iplogger
```

with golang
```
$ go mod tidy
$ go run .
```

now running on <http://localhost:10000>

## result

```go
forward := r.Header.Get("X-Forwarded-For")
  var ip string
  var err error
  if forward != "" {
    // With proxy
    ip = strings.Split(forward, ",")[0]
  } else {
    // Without proxy
    ip, _, err = net.SplitHostPort(r.RemoteAddr)
    if err != nil {
      http.Error(w, "Error parsing remote address ["+r.RemoteAddr+"]", http.StatusInternalServerError)
      return
    }
  }
```

As a de facto standard, when passing through a proxy, it is added one by one to the `X-Forwarded-For` header.
Therefore, if the corresponding header exists, `X-Forwarded-For` 0th IP is regarded as the client ip, and if not, RemoteAddr is used as the client ip.

사실상 표준 규격으로 프록시를 지날때 `X-Forwarded-For` 해더에 하나씩 추가되게 된다.  
따라서 해당 해더가 존재하는 경우 `X-Forwarded-For` 0번쨰 아이피를 client ip로 간주하고 없는 경우 RemoteAddr를 client ip로 사용한다.  

## but!!

> "절대로 유저에게서 오는 입력은 신뢰하지 마라"

X-Forwarded-For 헤더는 클라이언트가 임의로 조작할 수 있다.  

다음과 같이 간단하게 유조할 수 있다.  

```bash
curl -H "X-Forwarded-For: 1.1.1.1" http://localhost:10000
```

이걸 어떻게 고쳐야될까?  
좋은 go web framework인 echo에서 구현을 참고해보았다.  

[echo guide - ip address](https://echo.labstack.com/guide/ip-address/)  
[echo source - ip address](https://github.com/labstack/echo/blob/v4.10.0/ip.go)  

echo 구현 방식을 보면 프록시가 존재하는 경우에는 TrustOption을 통해 XFF IP 신뢰 여부를 결정하는데, 기본적으로 설정할 수 있는 옵션은 다음과 같다.  

루프백 주소를 신뢰할 것인가?, 링크로컬 주소를 신뢰할 것인가?, 프라이빗 주소를 신뢰할 것인가?  
그리고 사용자 지정 추가 주소를 설정할 수 있다.  

이게 완성된 구현이 존재하는 마당에 다시 만드는게 바퀴의 재발명 같을 수 있지만, 이렇게 구현해보는 것도 좋은 공부가 될 것 같아서 구현해보았다.  

구현된 패키지는 다음과 같이 사용할 수 있다.  

```go
package main

import (
	"log"
	"net/http"

	"github.com/minpeter/iplogger/pkg/ip"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientIP := ip.GetIP(r)

		w.Write([]byte(clientIP))
	})

	http.ListenAndServe(":8080", nil)
	log.Println("Server started on port 8080")
}
```

기본적으론 XFF 헤더를 사용하고 루프백, 링크로컬, 프라이빗 주소를 신뢰하고, 클라우드플레어 프록시 주소를 신뢰한다.  
추가적으로 사용자가 TrustOption을 통해 신뢰할 프록시 주소를 추가할 수 있다.  

```go
ip.TrustOption.TrustIPRanges = []*net.IPNet{
	{IP: net.IPv4(10, 10, 10, 0), Mask: net.IPv4Mask(255, 255, 255, 0)},
}
clientIP := ip.GetIP(r)
```

이렇게 설정해주면 10.10.10.0/24 대역의 프록시에서 오는 XFF 헤더를 신뢰한다.

또 각각 루프백, 링크로컬, 프라이빗, 클라우드플레어를 신뢰하는지 여부를 설정할 수 있다.

```go
ip.TrustOption.TrustLoopback = false
ip.TrustOption.TrustLinkLocal = false
ip.TrustOption.TrustPrivate = false
ip.TrustOption.TrustCloudflare = false
```

처음에는 막막했지만 막상 구현하고 나니 간단해서 할 말이 없다.  
나중에 쓸 일이 있는 패키지인지는 모르겠지만 IP 로깅을 하려고 할 때 참고하면 좋을 것 같다.  
