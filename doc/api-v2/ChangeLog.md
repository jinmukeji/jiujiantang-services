## ChangeLog

### 2.0.4修改内容

```
账户等级不用 数字表示，用字符串
修改组织信息中 删除created_date，added_user_num，max_user_limits，level
新增多个组织的owner 这个描述有问题，改成 user
删除组织，取消组织拥有者，取消组织管理者，将用户移除组织，解绑 mac，不返回其他信息
修改用户信息中，不允许register_type ，gender
测量添加record_type
添加健康趋势开关
```

### 2.0.5版本修改内容

```
添加example
```

### 2.0.6版本修改内容

```
添加删除响应成功后的提示
组织下用户提供user_id
添加owner查看组织信息接口：响应数据中服务到期时间、固定电话
```

### 2.0.7版本修改内容

```
Owner新增多个组织的User接口命名错误，改成是向组织添加多个拥有者
Owner查看拥有的组织接口中成功响应的返回的data中organization_id和profile是分开的，但是Owner查看组织拥有者接口中user_id是在profile里面的，Owner查看组织下的用户接口也一样，查看用户个人档案接口，修改个人档案接口，统一一下
添加Owner取消一个Admin的组织管理者身份接口响应成功后不需要返回profile，返回提示
添加Owner将用户移出组织接口响应成功后不需要返回profile，返回提示
添加注销登录接口响应成功给提示
添加意见反馈接口响应成功给提示
Owner将用户移出组织接口中只是移除单个用户改成单个或者多个
example中把interger改成数字
```

### 2.0.8版本修改内容

```
注册一个用户接口Response中user_id放入profile中
查看用户个人档案接口Response中user_id放入profile中
修改个人档案接口Response中user_id放入profile中
Owner将用户移出组织位置不对，修改放入Owner中
```

### 2.0.9版本修改内容

```
删除查看用户个人档案Response中多了一个user_id
修改个人档案接口Response中user_id放入profile中
```

### 2.0.10版本修改内容

```
企业信息中添加邮箱
创建组织，修改组织时添加line
admin 查看组织下用户添加失败的结果
'所在省份'字错误纠正
查看历史记录添加智能分析的部分
```

### 2.0.11版本修改内容

```
添加判断mac是否关联的api
版本信息中添加stoneDescription
```

### 2.0.12版本修改内容

```
把http中的body提取出来
修改ower创建组织的body
修改获取版本api的返回值
添加是否关联mac的api文档
```

### 2.0.13版本的修改内容

```
提取用户，组织，测量结果的profile
Responses换成Response,""换成''
```

### 2.0.14版本的修改内容

```
添加owner在组织中添加多个用户的api
删除owner在组织中添加多个owner的api
Admin查看组织用户的api，orginization_id改成integer
```

### 2.0.15版本的修改内容

```
把XX-profile改成profile
添加数据的格式，最小值和枚举。修改schema
修改注释，命名
在User中添加国家
```

### 2.0.16版本的修改内容

```
error中去掉data
data中去掉error
把response中的data提取出来
```

### 2.0.17版本的修改内容

```
修改Device,User部分的api
```

