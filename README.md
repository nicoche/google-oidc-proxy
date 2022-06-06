# google-oidc-proxy

Forward proxy HTTP requests to resources protected by google OIDC (e.g. Identity Aware Proxy)

Inspired by https://github.com/awslabs/aws-sigv4-proxy / leverages https://pkg.go.dev/google.golang.org/api/idtoken

## Usage

```
$ GOOGLE_APPLICATION_CREDENTIALS=xxxx.json TARGET_AUDIENCE=yyy.apps.googleusercontent.com TARGET_HOST=zzz.com ADDRESS=localhost:8080 go run cmd/google-oidc-proxy/main.go &
$ curl localhost:8080/a/b/c
# http get proxified to zzz.com/a/b/c, with a valid google id token
```
