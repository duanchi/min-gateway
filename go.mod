module github.com/duanchi/min-gateway

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/duanchi/min v1.8.10
	github.com/gin-gonic/gin v1.5.0
	github.com/go-redis/redis/v8 v8.11.2
	github.com/gorilla/websocket v1.4.2
	github.com/satori/go.uuid v1.2.0
)

replace github.com/duanchi/min => ../min
