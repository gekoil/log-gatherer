# log-gatherer
Go app that gathers logs from docker containers

# Demo
- Run the server
```shell script
go run ./cmd/log-gatherer/main.go
```

- Launch containers to follow
For example:
```shell script
docker build ./containers/. -t foo:latest
docker run -d --name foo
```

- Call the APIs
```http request
PUT /attach HTTP/1.1
Host: localhost:8080
Content-Type: application/json
cache-control: no-cache

{
  "stdOut": true,
  "stdErr": true,
  "filter": {
    "name": "foo"
  }
}
```
