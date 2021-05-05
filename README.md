# Comparison of Go web frameworks

This repository includes some realstic examples of using popular go web frameworks:
- [echo](https://github.com/labstack/echo)
- [fiber](https://github.com/gofiber/fiber)
- [gin](https://github.com/gin-gonic/gin)
- [httprouter](https://github.com/julienschmidt/httprouter)

## Usage
Example command for /v1/bananas:
```bash
ht localhost:8080/v1/bananas
```

Example command for /v1/banana/:id:
```bash
ht localhost:8080/v1/banana/1234
```

Example command for /v1/admin with basic auth header:
```bash
ht -a foo:bar POST localhost:8080/v1/admin value=hello
```
