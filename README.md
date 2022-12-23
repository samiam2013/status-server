# status-server

simple json api to check system resources status via HTTP

set BASIC_PASS and BASIC_USER environment variables to configure auth (use .env file in root or `export`)

reverse proxy from nginx or similar on /status with https (or basic auth isn't secure)


then request /status to get json response

```bash
➜  status-server git:(main) ✗ curl localhost:8888/status -u "user:password" -v
```
```
*   Trying 127.0.0.1:8888...
* Connected to localhost (127.0.0.1) port 8888 (#0)
* Server auth using Basic with user 'user'
> GET /status HTTP/1.1
> Host: localhost:8888
> Authorization: Basic *************************
> User-Agent: curl/7.81.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Www-Authenticate: Basic realm="Restricted"
< Date: Fri, 23 Dec 2022 14:05:59 GMT
< Content-Length: 198
< Content-Type: text/plain; charset=utf-8
```
```json
{
    "time":1671804359,
    "up_time":"112h50m27s",
    "memory":{
        "used_perc":3.8350156024157673,
        "used_gib":2.587025408},
    "cpu":{
        "user":0.020842017507294707,
        "system":0.06252605252188412,
        "idle":99.91663192997082}
}
```
