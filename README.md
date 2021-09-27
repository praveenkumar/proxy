# proxy
Run http/https proxy server.

```
 $ ./proxy -h
Usage of ./pro:
  -ca-cert-path string
        CA cert path for TLS
  -ca-key-path string
        CA key path for TLS
  -host string
        host (domain/IP) for alternate dns in certificate
  -port int
        Port to serve non tls proxy (default 8080)
  -ssl-port int
        SSL port to serve tls proxy (default 8443)
```

If a user already have CA cert/key then use `--ca-cert-path`
and `--ca-key-path` otherwise it will generate and provide
the details about created cert/key.

Example
------

```
$ ./proxy
2021/09/27 14:25:42 CA cert file /tmp/proxy-ca290921658/proxy-cert.pem
2021/09/27 14:25:42 CA key file /tmp/proxy-ca290921658/proxy-key.pem
```

```
// HTTP proxy
$ curl -x http://localhost:8080 -Ik gandi.net
HTTP/1.1 301 Moved Permanently
Connection: keep-alive
Content-Length: 0
Date: Mon, 27 Sep 2021 08:56:58 GMT
Location: http://www.gandi.net/
Server: Varnish

// HTTPS proxy
$ curl --proxy-cacert /tmp/proxy-ca290921658/proxy-cert.pem -x https://localhost:8443 -Ik https://gandi.net
HTTP/1.0 200 OK

HTTP/1.1 301 Moved Permanently
Date: Mon, 27 Sep 2021 08:57:42 GMT
Server: Varnish
Location: https://www.gandi.net/
Content-Length: 0
Connection: keep-alive
```