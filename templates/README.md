# Templates

## 命名规则

命名空间：

- **API** 类型的微服务使用命名空间 `com.jinmuhealth.api`

- **SRV** 类型的微服务使用命名空间 `com.himalife.srv`

- **WEB** 类型的微服务使用命名空间 `com.jinmuhealth.web`

Docker Image 命名：

- 所有系统应用的微服务使用 `jm-app` 作为 Docker Image 的命名空间名。
- 使用 `<Docker Namespace>/<Service Name>` 格式。例如，微服务 `svc-biz-core` 的 Docker Image 的名称为 `jm-app/svc-biz-core`

## 命名速查表

| Service Name   | Nacmespace          | Type | Docker Image Name         | Remarks |
| -------------- | ------------------- | ---- | ------------------------- | ------- |
| rest-api-l-v2  | com.jinmuhealth.api | API  | jm-app/web-rest-api-l-v2  |         |
| rest-api-v2    | com.jinmuhealth.api | API  | jm-app/web-rest-api-v2    |         |
| svc-biz-core   | com.himalife.srv | SRV  | jm-app/svc-biz-core       |         |
| rest-websocket | com.jinmuhealth.web | WEB  | jm-app/web-rest-websocket |         |
| rest-wechat    | com.jinmuhealth.web | WEB  | jm-app/web-rest-wechat    |         |



### Type of Services

- API
- WEB
- SRV

**API** - Served by the **micro api**, an API service sits at the edge of your infrastructure, most likely serving public facing traffic and your mobile or web apps. You can either build it with HTTP handlers and run the micro api in reverse proxy mode or by default handle a specific RPC API request response format which can be found [here](https://github.com/micro/micro/blob/master/api/proto/api.proto).

**WEB** - Served by the **micro web**, a Web service focuses on serving html content and dashboards. The micro web reverse proxies HTTP and WebSockets. 

**SRV** - These are backend RPC based services. They’re primarily focused on providing the core functionality for your system and are most likely not be public facing. You can still access them via the micro api or web using the /rpc endpoint if you like but it’s more likely API, Web and other SRV services use the go-micro client to call them directly.

![Arch](arch.png)



Namespacing

The micro api and web will compose a service name of the namespace and first path of a request path e.g. request to api **/customer** becomes **go.micro.api.customer**.

The default namespaces are:

- **API** - go.micro.api
- **Web** - go.micro.web
- **SRV** - go.micro.srv

You should set these to your domain e.g *com.example.{api, web, srv}*. The micro api and micro web can be configured at runtime to route to your namespace.



- svc: RPC service
- api: An API Gateway that serves HTTP and routes requests to appropriate micro services. It acts as a single entry point and can either be used as a reverse proxy or translate HTTP requests to RPC.
- web: a web dashboard and reverse proxy for micro web applications. It behaves much the like the API reverse proxy but also includes support for web sockets.
- sidecar: provides all the features of go-micro as a HTTP service.



## Sync vs Async

#### Synchronous

![RequestResponse](request-response.png)



#### Asynchronous

![PubSub](pub-sub.png)



#### Example: Audit

![Audit](audit.png)



## 安装 Micro 框架

> https://micro.mu/docs/install-guide.html

#### 更新 micro

参考: https://github.com/micro/micro

```sh
go get github.com/micro/micro
# Or update
# go get -u github.com/micro/micro
docker pull microhq/micro
```

#### 安装 Consul

```sh
brew install consul
```

#### 安装 Protobuf

```sh
# install protobuf
brew install protobuf

# install protoc-gen-go
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}

# install protoc-gen-micro
go get -u github.com/micro/protoc-gen-micro
```
