## 现存WEB API清单

* Api-v2中存在的web API
  * `POST` `/client/auth` 授权接口
  * `POST` `/users/signin` 登录接口
  * `PUT` `/owner/measurements/{record_id:int}/remark` 修改备注接口
  * `POST` `/owner/measurements/{record_id:int}/analyze` 智能分析接口
  * `POST` `/res/getUrl` 获取res中的资源文件接口
  * `GET` `/owner/measurements/{record_id:int}/token` 导出分析报告token接口
  * `GET` `/owner/measurements/token/{token:string}/analyze` 通过token获取分析报告接口

* Api-l-v2中存在的web API
  * `POST` `/client/auth` 授权接口
  * `POST` `/account/signin` 登录接口
  * `POST` `/measurements/{record_id:int}/analyze` 智能分析接口
  * `GET` `/measurements/{record_id:int}/analyze` 获取分析报告接口

* wechat中存在的web API
  * `GET` `/wx/oauth` 微信OAuth登录接口
  * `GET` `/wx/oauth/callback` 微信OAuth登录回调接口
  * `GET` `/wx/api/measurements` 微信测量接口
  * `PUT` `/wx/api/measurements/{record_id:int}/remark` 微信修改备注接口
  * `GET` `/wx/api/measurements/{record_id:int}/analyze` 微信获取分析报告接口
  * `POST` `/wx/api/payment` 微信伪支付接口
  * `POST` `/wx/api/jssdk/config` 微信获取JS配置接口
  * `GET` `/wx/api/measurements/{record_id:int}/token` 微信导出分析报告token接口
  * `GET` `/wx/api/measurements/token/{token:string}/analyze` 微信通过token获取分析报告接口
  * `POST` `/wx/api/res/getUrl` 微信获取res中的资源文件接口

* 其中获取web static 页面的接口
  * Api-v2中的`POST` `/res/getUrl` 获取res中的资源文件接口
  * wechat中的`POST` `/wx/api/res/getUrl` 微信获取res中的资源文件接口

