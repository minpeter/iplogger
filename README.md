# ipLogger

👀 traefik reverse proxy real ip test

## purpose

A study on how services located behind multiple reverse proxies log real client IPs.

여러 리버스 프록시 뒤에 위치한 서비스가 실제 클라이언트 IP를 기록하는 방법에 대한 연구

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

.. to be continued
