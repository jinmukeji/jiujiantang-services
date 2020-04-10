# REST API V2 代理微服务

### 1 启动微服务

**启动运行：**

```sh
go run main.go \
  --server_address=:8080
```

**测试：**

```sh
http GET http://localhost:8080/version \
  Cache-Control:no-cache
```



### 2 使用 `V2-api` 作为基准URL启动微服务

**启动运行：**

```sh
go run main.go \
  --server_address=:8080 \
  --x_api_base=v2-api
```

**测试：**

```sh
http GET http://localhost:8080/v2-api/version \
  Cache-Control:no-cache
```

