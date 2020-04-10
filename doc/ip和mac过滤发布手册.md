## 发布内容：根据不同的客户端实现ip和mac的过滤

### 前期准备

* 确认配置文件pkg/blocker的config_doc.yml已经更新完毕，并且config_doc.yml和相关代码已经合并

### 发布代码到仓库

* 使用Jenkins把代码发布到仓库，tag 是2.1.1

### 修改配置文件

* 修改`docker-compose.yml`文件，把原来的`2.1.0`换成`2.1.1`
* 修改`local.svc-biz-core.env`文件，编辑文件，在文件的末尾新加两行
```
X_BLOCKER_CONFIG_FILE=/blocker/config_doc.yml
X_BLOCKER_DB_FILE=/blocker/GeoLite2-Country.mmdb.gz
```
### 停机

* 登录2台正式环境的机器

* 运行`sudo docker-compose down`

### 发布

* 拉取镜像 `sudo docker-compose pull`
* 启动镜像 `sudo docker-compose up -d`
* 查看镜像运行情况 `sudo docker ps -a`

### 验证

* 测试可以使用app进行验证

  * 功能测试 使用app内部测量的模块

  * API接口测试 可以使用API接口测试代码跑测量接口


