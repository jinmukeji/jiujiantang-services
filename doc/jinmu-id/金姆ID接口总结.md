## 喜马把脉ID接口总结
### 1. 喜马把脉ID所有接口
-  查看版本信息 GET /version 
-  提交客户端授权 POST ​/client​/auth 
- 手机号获取信息 POST ​/notification​/sms 
- 邮件通知 POST ​/notification​/email
- 登录 POST ​/signin
- 获取登录二维码 POST ​/wx​/signin​/qrcode
- 设置显示语言 POST ​/user​/{user_id}​/language
- 得到显示语言 GET ​/user​/{user_id}​/language
- 获取用户的个人档案 GET /user​/{user_id}​/profile
- 修改个人档案 PUT ​/user​/{user_id}​/profile
- 用户设置安全邮箱 POST ​/user​/{user_id}​/safe_email
- 用户解除设置安全邮箱 POST ​/user​/{user_id}​/safe_email​/delete
- 用户设置密码 POST ​/user​/{user_id}​/password 
- 用户修改密码 PUT ​/user​/{user_id}​/password
- 用户设置用户名 POST ​/user​/{user_id}​/username
- 用户首次设置密保问题 POST ​/user​/{user_id}​/set_secure_question
- 修改密保问题前验证原来的密保问题 POST ​/user​/{user_id}​/ validate_question_before_modify
- 修改密保问题 POST ​/user​/{user_id}​/modify_secure_question
- 根据密保问题修改密码前验证手机号码或者用户名是否存在 POST ​/user​/validate_username_or_phone 
- 根据密保问题修改密码前验证回答的密保问题是否正确 POST ​/user​/validate_question_before_modify_password 
- 根据密保问题修改密码 POST /user​/modify_password_via_question
- 获取绑定微信二维码 POST ​/wx​/user​/{user_id}​/connect​/qrcode
- 获取解绑微信二维码 DELETE ​/wx​/user​/{user_id}​/connect​/qrcode 
- 根据验证码重置密码 POST ​/user​/{user_id}​/reset_password
- 获取所有密保问题列表 GET ​/user​/{user_id}​/secure_question_list
- 选择区域 POST ​/user​/region
- 设置登录手机号码 POST ​/user​/{user_id}​/signin_phone
- 修改登录手机号码 PUT /user​/{user_id}​/signin_phone
- 验证登录手机号码 POST ​/validate_signin_phone
- 获取最新验证码 POST ​/_debug​/user​/latest_verification_code
- 获取资源列表 GET ​/resource
- 验证邮箱验证码是否正确 POST ​/user​/validate_email_verification_code
- 修改密保前获取已设置密保列表 GET ​/user​/{user_id}​/secure_question_to_modify
- 根据邮箱找回用户名 POST /user​/find_username_by_email
- 根据用户名或者手机号获取当前设置的密保问题 POST /user​/secure_question
- 修改安全邮箱 POST ​/user​/{user_id}​/modify_secure_email
- 获取用户的首选项配置 GET ​/user​/{user_id}​/preferences
- user 查询订阅 GET ​/user​/{user_id}​/subscription
- 使用验证码注册 POST ​/signup​/mvc
- 用户得到使用的device GET ​/user​/{user_id}​/device
- 注销用户登录 POST /user​/signout
### 2. 接口与业务流程
| 模块       | 功能点                   | 备注         |
| ---------- | ------------------------ | ------------ |
| 注册       | 手机号注册               |              |
|            | 微信注册                 |              |
| 登录       | 手机号验证码登录         |              |
|            | 手机号密码登录           |              |
|            | 用户名密码登录           |              |
|            | 微信登录                 |              |
|            | 忘记密码>重置            |              |
|            | 忘记用户名>找回          |              |
| 账号安全   | 手机号设置、修改         |              |
|            | 微信绑定、解绑           |              |
|            | 密码设置、修改           |              |
|            | 安全邮箱设置、修改、解绑 |              |
|            | 密保问题设置、修改       |              |
|            | 最近使用的登录设备       | 手机、电脑等 |
| 个人信息   | 基础资料                 |              |
|            | 首选项                   |              |
| 设备与服务 | 使用过的脉诊仪           |              |
|            | 正在使用的服务           |              |
| 通知提醒   | 设置是否接收消息         |              |
| 隐私政策   | 隐私政策网站跳转         |              |
| 帮助中心   | 联系客服                 |              |

#### 2.1 用户注册
 * 使用手机号码注册
用户输入手机号码之后，输入验证码，调用此时调用接口**手机号获取信息**接口发送验证码，用户接受验证码之后填入验证码，调用**使用验证码注册**，传入用户输入的手机号码和验证码等信息 ，返回token和user_id。下面进入完善资料页面。
* 使用微信登录
暂无
#### 2.2 登录
##### 2.2.1 登录手机号码验证码登录
点击获取验证码，调用**邮件通知**
输入验证码之后，调用**登录**，此时登录方式为手机验证码登录
** 手机号密码登录
输入手机号码和密码之后，调用**登录**，此时传入的登录方式为手机密码登录
** 用户名密码登录
输入用户名和密码，点击登录，调用接口**登录**，此时登录方式为用户名密码登录
** 微信登录
暂无
##### 2.2.2 忘记密码>重置
* 使用登录手机重置
首先验证手机号码：输入手机号码，点击获取验证码，调用**手机号获取信息**，填写验证码之后调用**验证登录手机号**，该接口会查看当前手机号码是否存在和验证码是否正确。
然后是重置密码，密码输入完成后调用**根据验证码重置密码**
* 使用安全邮箱重置
首先验证邮箱，输入邮箱点击获取验证码，点击获取验证码，调用**邮件通知**， 填写验证码之后调用**验证邮箱验证码是否正确**，验证当前填入的验证码和邮箱是否对应。
然后是重置密码，密码输入完成后调用**根据验证码重置密码**
* 使用安全问题重置
首先输入用户名或者手机号，调用接口**根据密保问题修改密码前验证手机号码或者用户名是否存在**，返回验证是否正确。到达下一个页面，验证密保问题，获取用户已经设置的密保问题列表根据用户名或者手机号码调用接口**根据用户名或者手机号获取当前设置的密保问题**，返回用户已经设置的密保问题。用户回答完毕之后，点击确定，调用**根据密保问题修改密码前验证回答的密保问题是否正确**，如果用户回答错误，会返回错误问题的Key。若回答正确，则进入重置密码。输入密码之后，调用**根据密保问题重置密码**
* 未进行任何绑定
##### 2.2.2 忘记用户名>找回
* 使用安全邮箱找回
填写邮箱，调用**邮件通知**获取验证码，填写验证码之后，点击下一步之后调用**根据邮箱找回用户名**，该接口会验证邮箱是否存在，验证码是否正确，并且返回用户的用户名。
#### 2.3 帐号安全
* 帐号安全页面
进入该页面，获取用户的帐号安全信息，调用接口**获取用户的个人档案**
##### 2.3.1 手机号设置、修改
* 设置手机号码，调用**设置登录手机号码**
* 修改手机号
首先输入当前手机号码，获取验证码，调用**手机号获取信息**，输入验证码之后点击下一步，输入新的手机号，点击获取验证码，调用**手机号获取信息**，确定后调用**修改登录手机号** ，这个接口会判断手机号是否与原来的手机号相同
##### 2.3.2 微信绑定，解绑
##### 2.3.3 密码设置、修改
* 设置密码
输入密码与确认密码：调用**用户设置密码**
* 修改密码
输入旧的密码，输入两次新密码，点击确定后修改密码，调用接口**用户修改密码**，该接口会验证新旧密码是否相同，并且会设置用户密码为输入的新密码
##### 2.3.4 安全邮箱设置、修改，解绑
* 设置安全邮箱
输入安全邮箱 ，点击获取验证码，调用**邮件通知**，返回serial_number
输入验证码之后，点击确定，调用**设置安全邮箱**
* 安全邮箱修改
此时邮箱已经设置
输入旧的安全邮箱，调用**邮件通知**发送验证码，输入验证码完毕之后调用**验证邮箱验证码是否正确**，验证成功后会返回标识符verification_number，接下来输入新邮箱并且调用**邮件通知**发送验证码，验证码填写完毕之后也需要把verification_number传过去调用接口**修改安全邮箱**
* 安全邮箱解绑
输入旧的安全邮箱，调用**邮件通知**发送验证码，输入验证码完毕之后调用**用户解除设置安全邮箱**，该接口实现验证验证码是否正确并且实现解绑邮箱功能
* 最近使用的登录设备
##### 2.3.5 密保问题设置
* 设置密保问题
有三个密保问题，首先获取密保问题列表
调用接口**获取所有密保问题列表**
答案填完后，点击确定，调用**用户首次设置密保问题**

##### 2.3.6 密保问题修改
* 验证密保问题
调用接口**修改密保前获取已设置密保列表**
答案填完后，点击确定，调用**修改密保问题前验证原来的密保问题**
调用接口**获取所有密保问题列表**，选择新密保问题
答案填完后，点击确定，调用**修改密保问题**

#### 2.4 个人信息
##### 2.4.1 基础资料
* 获取
进入个人信息页面，调用**获取用户的个人档案**
* 设置
选择性别
设置生日
设置身高
设置体重：调用**修改个人档案**
*设置常用语言，调用**设置显示语言**
* 修改
基础资料修改完成后确认，调用接口**修改个人档案**
##### 2.4.2 首选项
* 设置语言
语言选择结束之后调用**设置显示语言**，该账号能够登录的喜马把脉所有的应用程序在下次打开时会自动更换成该语言
得到用户显示的语言调用**得到显示语言**
* 地区设置
地区设置后不允许修改
#### 2.5 设备与服务
##### 2.5.1 使用过的脉诊仪
进入使用过的脉诊仪界面，请求接口**用户得到使用的device**
##### 2.5.2 正在使用的服务
#### 2.6 通知提醒
新建的默认选择
### 3. 异常处理方案
#### 3.1 注册
* 使用手机号码注册
判断手机号码格式是否正确，不正确则提醒用户。
判断手机号码是否已经被其他人注册了，已经被注册了则提醒用户。
* 微信注册
#### 3.2 登录
##### 3.1.1 手机验证码登录
调用**登录**，此时请求的请求方式为手机号码验证码，如果当前手机号码没与被任何人设置，**登录**接口则返回错误码，在用户获取验证码之后提醒用户该手机号码没有被绑定。
##### 3.1.2 手机号码密码登录
调用**登录**，此时请求的请求方式为手机号码密码登录，如果当前手机号码没与被任何人设置，在用户输入手机号码和密码结束之后，则返回错误码，在用户获取验证码之后提醒用户该手机号码没有被绑定
##### 3.1.3 用户名密码登录
调用**登录**，此时请求的请求方式为用户名密码登录，判断用户名是否存在，不存在的话**登录**异常提示：`用户名不存在`。如果用户名存在，才进行正常的用户名与密码的校验流程 。
##### 3.1.4 微信登录
##### 3.1.5 忘记密码>重置
* 登录手机找回
输入手机号，**手机号获取信息**验证当前手机号码是否被设置，是否有人注册了这个手机号。有的话则进行提示。
当前手机号码已经注册过喜马把脉ID，才发送验证码。
**验证登录手机号**判断验证码验证是否过期，是否已经被使用，成功后则填写新密码，密码判断密码为8-20个字符，同时包含数字和字母，并且不能与旧密码相同。否则的话会进行提示。
重置密码之后跳转到登录，重新登录。
* 安全邮箱找回
输入安全邮箱，**邮件通知**验证当前手安全邮箱是否被设置，是否有人注册了这个手机号。有的话则进行提示。
当前安全邮箱已经注册过喜马把脉ID，才发送验证码。
**验证邮箱验证码是否正确**验证码验证是否过期，是否已经被使用，成功后则填写新密码，密码判断密码为8-20个字符，同时包含数字和字母，并且不能与旧密码相同。否则的话会进行提示。
重置密码之后跳转到登录，重新登录。
* 密保问题找回
首先验证手机号码或者用户名，**根据密保问题修改密码前验证手机号码或者用户名是否存在**判断当两个都填写的时候直接返回错误，不再进行下面的获取用户的密保问题步骤。用户名或者手机号码只能填写其中的一项，对用户名或者手机号码进行验证是否存在，存在且唯一的话则**根据用户名或者手机号获取当前设置的密保问题**返回对应的密保问题。
返回密保问题之后，用户填写每个密保问题对应的答案，**根据密保问题修改密码前验证回答的密保问题是否正确**这里对答案不做格式判断，只验证是否正确。验证成功后则可以输入密码，否则返回会在错误答案的位置对应提示答案错误。密保问题回答正确后则可以输入新密码，**根据密保问题重置密码**对新密码仍然做格式判断。
##### 3.1.6 忘记用户名>找回
* 根据安全邮箱找回
输入邮箱后判断邮箱是否格式正确，点击获取验证码之后判断是否绑定，不正确均返回错误码。验证码填写完毕后请求接口，返回对应的用户名。
#### 3.3 帐号安全
##### 3.3.1 手机号设置、修改
* 调用**设置登录手机号码**，验证手机号码是否已经注册了和格式是否正确，不正确均返回错误码。
* 修改手机号
向原来手机号码发送验证码的时候验证手机号码是否是当前用户的，验证码验证成功之后，输入新的手机号，确定后调用**修改登录手机号** ，这个接口会判断手机号是否与原来的手机号相同，判断新手机号码是否被其他用户设置了，错误的话都会返回对应提示。
##### 3.3.2 密码设置、修改
* 设置密码
输入密码与确认密码：调用**用户设置密码**，验证密码的长度，字符等信息，不符合的则返回提示内容。
* 修改密码
输入旧的密码，输入两次新密码，点击确定后修改密码，调用接口**用户修改密码**，该接口会验证新旧密码是否相同，并且会设置用户密码为输入的新密码
##### 3.3.3 安全邮箱设置、修改，解绑
* 设置安全邮箱
验证邮箱格式
调用**邮件通知**，根据不同场景下的邮箱通知，判断在除了设置安全邮箱与修改安全邮箱之外的场景邮箱必须是被已经设置的，**设置安全邮箱**判断邮箱不能被其他人设置了
* 安全邮箱修改
用户输入邮箱之后，**邮件通知**判断旧邮箱是否与用户的邮箱一致，填写新邮箱后**修改安全邮箱**判断新旧邮箱不能相同，否则会返回对应提示。
* 安全邮箱解绑
输入旧的安全邮箱，**邮件通知**判断旧邮箱是否与用户的邮箱一致。
##### 3.3.4 密保问题设置
* 设置密保问题
调用接口**获取所有密保问题列表**，返回密保问题的列表，返回对应的问题的key和描述。
用户填写答案，**用户首次设置密保问题**判断答案是否满足在15字符内，不包含敏感字符等信息，否则的话都会给出提示。请求服务器时传入问题的key，接口会判断key是否在指定范围内，不在的话会给出提示。
##### 3.3.5 密保问题修改
新答案填完后，点击确定，调用**修改密保问题**，这个时候该接口会对传进来的旧密保问题进行验证之后再进行修改。
#### 3.4 个人信息
##### 3.4.1 基础资料
* 设置
选择性别，一旦设置后不可修改
设置生日
设置身高
设置体重：调用**修改个人档案**，会判断所有的参数是否合法，是否在指定的范围
调用接口**修改个人档案**，会判断所有的参数是否合法，是否在指定的范围
