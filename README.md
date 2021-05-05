# Comparison of Go web frameworks

This repository includes some realstic examples of using popular go web frameworks:
- [echo](https://github.com/labstack/echo)
- [fiber](https://github.com/gofiber/fiber)
- [gin](https://github.com/gin-gonic/gin)
- [httprouter](https://github.com/julienschmidt/httprouter)

## Starting the examples
Starting the 4 different examples can be done like follows:

echo (port `8080`):
```bash
go run ./cmd/echo
```

fiber (port `8081`):
```bash
go run ./cmd/fiber
```

gin (port `8082`):
```bash
go run ./cmd/gin
```

httprouter (port `8083`):
```bash
go run ./cmd/httprouter
```

## Statsd

Start the statsd lister with docker-compose:
```bash
docker-compose up -d
```

and follow the logs:
```bash
docker-compose logs -f
```

## Accessing the examples
These commands use [httpie](https://httpie.io/). Though I prefer the
[go version](https://github.com/nojima/httpie-go), as it starts super quick,
unlike the Python one.

Example command for `/v1/bananas`:
```bash
http localhost:8080/v1/bananas
```

Example command for `/v1/banana/:id`:
```bash
http localhost:8080/v1/banana/1234
```

Example command for `/v1/admin` with basic auth:
```bash
http -a foo:bar POST localhost:8080/v1/admin value=hello
```
