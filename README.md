# Heron-Gateway

## 初始化

使用预先生成的docker镜像初始化微服务网关

如果需要认证和授权, 需要先启动一个redis作为令牌存储

```shell
docker run -d --name gateway -p 9080 -p 9800 -e AUTHORIZATION_DSN="redis://127.0.0.1:6379/0" docker.heurd.com/heron-go/gateway
```

## 网关配置

启动完成后通过 `http://ip:9800`进入网关配置页面进行服务和路由配置

## 认证和授权

在网关配置界面选择`需要授权`后, 网关将会使用`JWT Bearer`的方式进行token的生成、验证与刷新

token默认有效期为7200秒, 也可以设置`JWT SIGNATURE`进行`jwttoken`签名

对应环境变量为

```
AUTHORIZATION_SIGNATURE_KEY = ''
AUTHORIZATION_TTL = 7200
```



> 网关进行token验证和授权的过程对微服务是透明的, 不需微服务额外实现对token的验证和授权, 但这并不影响微服务授权信息(主要是用户信息)在token中的存储与传递

## 数据传递

网关进行认证和授权后会生成全局唯一的token返回至客户端, 并且在网关连接的redis中会存储token及对应存储的数据(这部分数据可以是从微服务传递至网关), 下次当客户端携带正确的token访问时, 网关将把这部分数据传递至微服务

可传递的数据有

#### `X-Gateway-Request-Id`

用与存储当前客户端请求的唯一请求标识, 通常用作请求追踪



#### `X-Gateway-Data`

使用`base64`编码存储，解码后通常为一个json字符串, 默认结构如下

```json
{
  "url": "request-url",
  "data": {
    "token": "jwt token",
    "user": "user ID",
    "more": "any more data"
  }
}
```



> 上述数据传递方式为网关向微服务HTTP请求的header部分,
>
> **`X-Gateway-Data`内容长度应不超过512字节**

> 若微服务使用`heron`框架, 可以初始化框架内的 `middleware.GatewayMiddleware ` 中间件从而自动解析`X-Gateway-Data`内容至 `gateway.Data`类型的结构体中, 并通过 `ctx.Get('GATEWAY_DATA')`获取



#### `X-Gateway-Authorization-Action`

此数据传输为反向传输, 即微服务反馈至网关, 通过微服务HTTP应答的Header反馈至网关

此数据反馈给网关认证和授权相关结果，如用户登录时, 请求微服务后微服务若验证用户登录成功, 则在应答Header中加入此header数据, 

可选值如下:

`CREATE`: 生成token, 反馈至网关生成jwt token, 同时反馈给网关需要存储的数据, 以`X-Gateway-Authorization-Data`反馈

`REFRESH`: 刷新一个token

`REMOVE`: 移除token



#### `X-Gateway-Authorization-Data`

此数据传输为反向传输, 即微服务反馈至网关, 通过微服务HTTP应答的Header反馈至网关

> 数据长度不应超过`256`字节

## 内置API

内置API是网关自身进行操作的restful-api, 路由可以通过 `GATEWAY_NATIVE_API_PREFIX`进行配置, 同时可以定义 `GATEWAY_NATIVE_API_ACCESS_TOKEN`保证访问的安全性

#### `PUT /routes`

刷新网关路由配置缓存

#### `PUT /services`

刷新网关微服务配置缓存

> 为了保证网关的性能, 网关配置是在网关启动时将配置预加载至内存缓存中, 对网关配置更改时, 特别是直接修改了配置文件`/heron/data/*.json`时, 需要调用上面的api对网关的缓存进行刷新才能够刷新配置
>
> **使用网关提供的配置管理界面时, 不需要额外调用上述内置API**