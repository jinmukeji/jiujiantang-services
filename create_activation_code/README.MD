# 喜马激活码生辰工具

## 运行环境

* macOS Version 10.15.4 

## 运行脚本

订阅类型说明：
0：定制化 1：试用版 2：黄喜马 3：白喜马 4：钻石喜马 5：礼品版

```sh
go run create_activation_code.go -username=数据库用户名 -password='数据库密码' -address=数据库地址 -database=数据库名 -subscription_type=订阅类型 -contract_year=年限 -max_user_limits=人数 -key=密钥 -number=生成数目
```
