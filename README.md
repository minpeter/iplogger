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

## screenshot
1. Create baseurl using host header  
![image](https://user-images.githubusercontent.com/62207008/213464093-97febea2-edb1-4580-a351-d57675538159.png)
2. Your connection ip address - <https://ip.minpeter.cf/test>
