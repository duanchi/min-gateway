# min-gateway

## 1 初始化

----

min-gateway推荐基于docker部署，可通过docker镜像直接进行网关部署

> 在部署之前，请确认已有redis实例和MySQL实例在运行中

1. 将数据库初始化文件导入数据库并配置数据库连接字符串

2. 通过docker部署
   
   ```shell
   docker run -d --restart always \
   --name gateway -p 9080 \
   -e AUTHORIZATION_DSN="redis://127.0.0.1:6379/0" \
   -e CACHE_DSN="redis://127.0.0.1:6379/1" \
   -e DB_DSN="mysql://{DB_USER}:{DB_PASS}@{DB_HOST}:{DB_PORT}/{DB_NAME}" \
   -e GATEWAY_CONSOLE_API_ACCESS_TOKEN="{CUSTOM-ACCESS-TOKEN}" \
   docker.in-mes.com/min/min-gateway
   ```

3. 需要GUI配置，在任意可访问网关的位置部署GUI管理环境
   
   ```shell
   docker run -d \
   --name gateway-ui -p 9800 \
   docker.in-mes.com/min/min-gateway-ui
   ```

4. 进入`http://your.host:9080`，输入网关地址`http://gateway.host:9080`和`GATEWAY_CONSOLE_API_ACCESS_TOKEN`即可进行网关管理

## 2 网关配置

网关配置分为服务配置路由配置，用于管理网关API路由和路由所指向的具体服务实例

### 2.1 服务配置

服务可通过以下几种方式进行配置

#### 2.2 GUI配置服务

可在GUI管理界面直接配置，可配置多组服务实例和灰度实例

#### 2.3 内置初始化服务配置

可在docker运行时通过置入特定环境变量实现内置服务配置

```shell
-e GATEWAY_AUTO_SERVICE_intellimes-service-biz=
-e GATEWAY_AUTO_SERVICE_intellimes-service-iot-executor=http://172.20.0.9:9501
-e GATEWAY_AUTO_SERVICE_intellimes-service-iot=http://172.20.0.11:8365,http://172.20.0.12:8365
```

置入`GATEWAY_AUTO_SERVICE_`开头的环境变量，在其后添加服务名称

> 如果不添加变量值，则只添加服务，不添加服务实例
> 
> 可以添加一个或多个服务实例，用`,`分割

> 通过环境变量初始化的服务配置，不能使用GUI删除或修改实例，但可以添加普通实例或灰度实例

#### 2.4 注册中心服务发现

在网关中维护的服务，通过与注册中心绑定，可以实现实例自动上下线，自动添加灰度实例等功能，免于手动维护服务实例

> 通过注册中心同步的普通实例和灰度实例，可以进行删除操作，但下一个同步周期会再次全量同步注册中心的实例

#### 2.5 灰度发布和灰度测试

> 灰度发布涉及多种配置，暂未实现

可以通过在前端请求中置入`X-Gateway-Instance-Id` 即可通过灰度实例访问后端接口

在GUI界面中添加灰度实例时，会自动生成Instance-Id

使用注册中心服务发现时，在Discovery-Client中配置Metadata的`instance-id`项，即可自动同步至网关灰度实例

> Instance-Id需遵循`UUID`规范

### 3 路由配置

支持设置前缀、匹配、正则、fnmatch类型的API规则路由

路由配置中可以使用网关内置的认证和授权机制，若使用内置认证和授权，需按照`使用网关认证和授权`一节中的规范进行应用改造

路由支持配置单目标登录，勾选`授权路由`并设置基于Header和querystring的特征字符串进行不同目标登录匹配限定，从而实现如同时只能在一个浏览器、一部手机登录等限定功能。

可配置路由访问超时时间，超时后强制结束接口请求

路由也支持websocket的API转发

> 带有网关jwt认证的websocket请求需要将jwt认证字符串通过querystring形式添加在token字段中

## 4 使用网关认证和授权

### 4.1 认证授权配置

在路由配置中勾选`需要授权`，则此路由将默认在网关进行认证， 认证通过后将认证信息带入服务实例请求中。

> 网关进行token验证和授权的过程对微服务是透明的, 不需微服务额外实现对token的验证和授权, 但这并不影响微服务授权信息(主要是用户信息)在token中的存储与传递

> 若服务实例需要对认证信息进行修改、删除等操作，需按照下面的`认证数据传递`、`认证状态控制`进行处理和应答。

### 4.2 默认授权和自定义授权

#### 默认授权

网关将会使用`JWT Bearer`的方式进行token的生成、验证与刷新

token默认有效期为7200秒, 也可以设置`JWT SIGNATURE`进行`jwttoken`签名

对应环境变量为

```
AUTHORIZATION_SIGNATURE_KEY = ''
AUTHORIZATION_TTL = 7200 // -1 为用不过期
```

#### 自定义授权

自定义授权网关将不参与授权token的生成过程，服务需按照`认证状态控制`中的授权规范返回相关字段，网关进行透明存储，并进行后续的认证信息验证

### 4.3 认证数据传递

当网关转发用户请求时，后端的认证(或用户登录)服务可正常根据账号密码等验证因子进行身份鉴权，当鉴权通过后，若使用默认授权，可通过认证信息应答通道像网关返回认证结果和此次认证携带的附属信息，由网关生成jwt token，并由客户端和网关之间实现认证鉴权过程。

> 客户端和网关之间的身份认证鉴权过程，后端服务无需参与

若使用自定义授权，与默认授权不同的是，客户端与网关之间鉴权的token信息需要后端的认证服务自行生成，并通过应答认证数据通道像网关返回。

> 同样，使用自定义授权时，客户端和网关之间的身份认证鉴权过程，后端服务无需参与

### 4.4 认证状态控制

#### 认证信息携带通道

认证信息携带通道是当客户端发送API请求时，通过网关的鉴权后，网关会将存储的鉴权token对应的用户认证信息通过代理请求携带至后端服务请求中，协助后端服务识别当前用户和辅助实现其他业务逻辑。

目前认证信息携带通道使用的是向服务端请求时添加请求头信息`X-Gateway-Data`实现



**`X-Gateway-Data`**

使用`base64`编码存储，解码后通常为一个json字符串, 默认结构如下

```json
{
  "url": "request-url",
  "data": {
    "token": "jwt token",
    "user": "user ID",
    "more": {
      "foo": "bar"
    }
  }
}
```

> 上述数据传递方式为网关向微服务HTTP请求的header部分,
> 
> **`X-Gateway-Data`内容长度应不超过512字节**

> 若微服务使用`min`框架, 可以初始化框架内的 `middleware.GatewayMiddleware` 中间件从而自动解析`X-Gateway-Data`内容至 `gateway.Data`类型的结构体中, 并通过 `ctx.Get('GATEWAY_DATA')`获取

#### 认证信息应答通道

认证信息应答通道用于用户登录或鉴权、用户注销等场景时，向网关应答用户信息，并通知网关进行相应token的生成、删除、更新等操作，目前使用在服务端返回应答时，在应答头中添加`X-Gateway-Authorization-Action`，`X-Gateway-Authorization-Data`，`X-Gateway-Authorization-Singleton`等信息，用于控制网关认证状态。



**`X-Gateway-Authorization-Action`**

此数据传输为反向传输, 即微服务反馈至网关, 通过微服务HTTP应答的Header反馈至网关

此数据反馈给网关认证和授权相关结果，如用户登录时, 请求微服务后微服务若验证用户登录成功, 则在应答Header中加入此header数据,

可选值如下:

`CREATE`: 生成token, 反馈至网关生成jwt token, 同时反馈给网关需要存储的数据, 以`X-Gateway-Authorization-Data`反馈

`REFRESH`: 刷新一个token

`REMOVE`: 移除token



**`X-Gateway-Authorization-Data`**

此数据传输为反向传输, 即微服务反馈至网关, 通过微服务HTTP应答的Header反馈至网关

> 数据长度不应超过`256`字节



**`Gateway-Authorization-Singleton`**



## 5 接口集成

### 5.1 接口集成配置

## 6 内置API

## 7 GUI配置界面

## 附: Road Map

- 调用日志
- 状态监控
- 告警信息
- 接口集成GUI