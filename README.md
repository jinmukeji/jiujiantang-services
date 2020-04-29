# 喜马把脉产品 - API服务

[TOC]

## 0 说明

本项目是 **gf-api** 微服务。

**fqdn:** `com.himalife.srv.svc-biz-core`

## 1 配置开发环境

### 1.1 GO 环境

- 安装 GO

```shell
brew install go
```

- 设置 GOPATH

  在 .bashrc 或 .zshrc 中设置

```shell
 export GOROOT=/usr/local/opt/go/libexec
 export GOPATH=$HOME/go
 export PATH=$GOPATH/bin:$PATH
 export GOBIN=$GOROOT/bin
```



### 1.2 Repository初始化

1. Clone 本项目
   ```sh
   cd $HOME
   git clone https://github.com/jinmukeji/gf-api.git
   ```

   > 开发人员应当 Clone 自己的 fork 版本
   
   注意不得 clone 项目至 GOPATH

### 1.3 安装依赖项

#### 1.3.1 安装 Consul

> macOS 下使用 Homebrew
>
> Linux 下使用 LinuxBrew

```sh
brew install consul
```

或者：

```sh
docker run consul
```
#### 1.3.3 配置go代理服务器(翻墙不好的情况下使用)

macOS or Linux 下使用 `export GOPROXY=https://goproxy.cn` or `echo "GOPROXY=https://goproxy.cn" >> ~/.profile && source ~/.profiles`

Windows 下使用 `$env:GOPROXY = "https://goproxy.cn"`

#### 1.3.4 Go module依赖管理

* 项目在$GOPATH路径外面

* 查看mod下所有依赖
```sh
go list -m all
```

* 更新依赖
```sh
go get -u all
```

## 2 开发调试

### 2.2 为本地测试设置配置文件

#### 2.1.1 环境变量配置文件

环境变量配置文件位于 `build` 文件夹下面。
- `template.svc-biz-core.env` 是一个配置文件模板。开发人员需要复制一份模板文件，并且命名为 `local.svc-biz-core.env`。
- `template.svc-sms-gw.env` 是一个配置文件模板。开发人员需要复制一份模板文件，并且命名为 `local.svc-sms-gw.env`。
- `template.svc-sem-gw.env` 是一个配置文件模板。开发人员需要复制一份模板文件，并且命名为 `local.svc-sem-gw.env`。
完成以上步骤，才能正常在 Docker Compose 或者 脚本启动方式下使用。

```sh
cd build
cp template.svc-biz-core.env local.svc-biz-core.env
cp template.svc-sms-gw.env local.svc-sms-gw.env
cp template.svc-sem-gw.env local.svc-sem-gw.env
```

**注意事项：**

- 首次创建 `local.svc-biz-core.env` 文件后，开发人员可以根据自己需要修改其中的环境变量内容；
- `local.svc-biz-core.env` 文件为每个开发人员自己本地使用的文件，**不应当**提交到 Git 项目之中；
- `local.svc-biz-core.env` 文件配置了敏感信息，例如数据库访问密码等，该文件应当妥善保管，防止泄露敏感信息。

### 2.2 启动Go微服务(microservice)

1. 首先，启动 Consul:

  ```sh
  consul agent -dev
  ```

2. 然后，启动当前 Service:

  ```sh
  cd service
  # 运行服务
  go run main.go \
     --x_db_address "localhost:3306" \
     --x_db_username "root" \
     --x_db_password "p@ssw0rd" \
     --x_db_database "jinmu" \
     --x_db_enable_Log \
     --x_db_max_connections 1
  ```

  其中启动参数含义可以运行 `go run main.go --help` 查看。

  或者，从环境配置文件里面加载环境变量，然后运行：

  ```sh
  env $(cat ../build/local.svc-biz-core.env | grep -v ^# | xargs) go run main.go
  ```

  配置文件位于 `build/local.svc-biz-core.env`

### 2.3 通过Micro Plugins启动 gRPC 微服务

有时需要为使用gRPC的客户端启动微服务(gprc不兼容micro的组件)。

```sh
cd service
# 运行服务
go run main.go grpc_plugins.go \
   --x_db_address "localhost:3306" \
   --x_db_username "root" \
   --x_db_password "p@ssw0rd" \
   --x_db_database "jinmu" \
   --x_db_enable_Log \
   --x_db_max_connections 1
```

或者，从环境配置文件里面加载环境变量，然后运行：

```sh
env $(cat ../build/local.svc-biz-core.env | grep -v ^# | xargs) go run main.go grpc_plugins.go
```

配置文件位于 `build/local.svc-biz-core.env`

### 2.4 运行Go客户端

```sh
cd service/client
go run main.go
```

### 2.5 使用 Micro API 网关进行HTTP方式的调试

首先，启动 Micro API网关，默认监听`8080`端口:

```sh
micro api --handler=rpc --namespace=com.himalife.srv
```

其中 `com.himalife.srv` 是RPC服务的命名空间

> 参考信息: https://micro.mu/blog/2016/04/18/micro-architecture.html



然后，通过 HTTP 请求调用RPC方法。以下例子调用了`gf-api`服务中的 `Jinmuhealth/Echo` RPC方法:

```sh
# 或者，使用 REST API 路由映射，使用 JSON 形式提交请求
http --json POST \
	http://localhost:8080/svc-biz-core/Jinmuhealth/Echo \
	content='Hello, Sky'

# 直接访问 /rpc 路径调用方法，使用 Form 表单形式提交请求
http --form POST \
	http://127.0.0.1:8080/rpc \
	service='com.himalife.srv.svc-biz-core' \
	method='Jinmuhealth.Echo' \
	request='{"content": "Hello, Sky"}'

# 直接访问 /rpc 路径调用方法访问短信网关
http  \
http://127.0.0.1:8080/rpc \
service='com.himalife.srv.svc-sms-gw' \
method='SMS.GetVersion' 

# 直接访问 /rpc 路径调用方法访问邮件网关
http  \
http://127.0.0.1:8080/rpc \
service='com.himalife.srv.svc-sem-gw' \
method='SEM.GetVersion' 


```

> 使用`http` 命令需要安装 [httpie](https://httpie.org/)

### 2.6 使用 Micro Web UI 进行调试

启动 Micro Web UI，默认监听`8082`端口:

```sh
micro web
```

浏览器访问: http://localhost:8082

> 参考: https://github.com/micro/micro/tree/master/web

### 2.7 查看本项目的Go文档

启动 godoc:

```sh
godoc -http=:6060 -src -timestamps -index
```

浏览器访问： http://localhost:6060/pkg/github.com/jinmukeji/jiujiantang-services/

## 3 使用 Docker

### 3.1 Build Docker Image

**Docker Image Name:** `svc-biz-core`

运行以下命令构建 Docker Image：

```sh
cd build
./go-build-all.sh
./docker-build-all.sh
```

### 3.2 运行 Docker Compose

启动各个 Docker Container：

```sh
cd build
docker-compose up
```

运行后，以下服务映射到Host：

- **Micro Web:** http://localhost:8082
- **Micro API Gateway:** http://localhost:8080
- **svc-biz-core 服务:** http://localhost:9090

退出运行后，执行以下命令清除已经停止的 Docker Container:

```sh
docker-compose rm -f
```

## 开发流程

- 本项目使用 [GitHub flow](https://guides.github.com/introduction/flow/) 流程进行源代码管理。**(注意不是 git-flow)**
- 本项目分为两个分支：
  - **master:** 主分支，且为默认分支。其它项目引用本项目时，均使用主分支。主分支体现最新稳定代码。
  - **develop:** 开发分支。本项目研发迭代过程使用develop分支。
- 无特殊必要，本项目不开 **feature** 或 **release** 分支

- Go 语言包使用 go modules 进行管理

## 开发环境

> 以 macOS 为例

```sh
# 更新 brew
brew update && brew upgrade

# gRPC - go
#		https://grpc.io/docs/quickstart/go/
go get -u google.golang.org/grpc

# Micro 框架
# 	https://micro.mu/docs/install-guide.html
go get -u github.com/micro/go-micro/v2
```

定期更新相关工具

```sh
# 更新工具

brew update && brew upgrade

go get -u google.golang.org/grpc

go get -u github.com/micro/go-micro/v2
```

### 定期更新 Go 包

```sh
go get -u all
```



## 4 开发人员必读资料

开发人员**必须阅读**以下内容资料：

- https://github.com/jinmukeji/arch-guide/issues/11 
