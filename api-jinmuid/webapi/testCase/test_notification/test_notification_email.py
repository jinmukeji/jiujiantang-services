import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class NotificationEmail(ApiTestCase):
    u'''邮箱获取信息'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')
    user_id=Util().get_user_id()

    def notification_email(self, type, email, language,user_id,send_to_new_if_modify):
        url = self.host + "notification/email"
        LOG.info("请求url:%s" % url)
        data = {"type": type,
                "email": email,
                "language": language,
                "user_id":user_id,
                "send_to_new_if_modify":send_to_new_if_modify
                }
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_notification_email_bind(self):
        u"""设置安全邮箱"""
        LOG.info("------邮箱获取信息用例：start!---------")
        result = self.notification_email('set_secure_email', data.email, data.language,self.user_id, False)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_notification_email_unbind(self):
        u"""解除绑定邮箱验证"""
        LOG.info("------邮箱获取信息用例：start!---------")
        result = self.notification_email('unset_secure_email', data.email, data.language,self.user_id,False)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_notification_email_modify_secure_email_old(self):
        u"""修改安全邮箱验证"""
        LOG.info("------邮箱获取信息用例：start!---------")
        result = self.notification_email('modify_secure_email', data.new_email, data.language, self.user_id,
                                         False)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_notification_email_modify_secure_email_new(self):
        u"""修改安全邮箱验证"""
        LOG.info("------邮箱获取信息用例：start!---------")
        result = self.notification_email('modify_secure_email', data.new_email, data.language,self.user_id,True)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_notification_email_reset_password(self):
        u"""重置密码邮箱验证"""
        LOG.info("------邮箱获取信息用例：start!---------")
        result = self.notification_email('reset_password', data.email, data.language,self.user_id, False)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_notification_email_find_username(self):
        u"""找回用户名邮箱验证"""
        LOG.info("------邮箱获取信息用例：start!---------")
        result = self.notification_email('find_username', data.email, data.language,self.user_id, False)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNoneExistSecureEmail)
        LOG.info("------pass!---------")

    def test_email_error_type(self):
        u'''type为空'''
        LOG.info("------邮箱获取信息用例：start!---------")
        result = self.notification_email('', data.email, data.language,self.user_id, False)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidEmailNotificationAction)
        LOG.info("------pass!---------")

    def test_email_error_email(self):
        u'''email为空'''
        LOG.info("------邮箱获取信息用例：start!---------")
        result = self.notification_email('set', '', data.language,self.user_id, False)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidEmailNotificationAction)
        LOG.info("------pass!---------")

    def test_email_error_language(self):
        u'''language为空'''
        LOG.info("------邮箱获取信息用例：start!---------")
        result = self.notification_email('set', data.email, data.language,self.user_id, False)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidEmailNotificationAction)
        LOG.info("------pass!---------")

    def test_email_error_email_format(self):
        u'''email格式不正确'''
        LOG.info("------邮箱获取信息用例：start!---------")
        result = self.notification_email('set', '123213', data.language,self.user_id, False)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidEmailNotificationAction)
        LOG.info("------pass!---------")

    def test_bind_email_exit(self):
        u'''绑定邮箱验证email已存在'''
        LOG.info("------邮箱获取信息用例：start!---------")
        result = self.notification_email('set',data.email, data.language,self.user_id, False)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result,const.ErrInvalidEmailNotificationAction)
        LOG.info("------pass!---------")

    def test_unbind_email_no_exit(self):
        u'''解除绑定邮箱email不存在'''
        LOG.info("------邮箱获取信息用例：start!---------")
        result = self.notification_email('unset',data.email_not_exist,  data.language,self.user_id, False)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidEmailNotificationAction)
        LOG.info("------pass!---------")

    def test_reset_password_email_no_exit(self):
        u'''重置邮箱email不存在'''
        LOG.info("------邮箱获取信息用例：start!---------")
        result = self.notification_email('reset_password',data.email_not_exist,  data.language,self.user_id, False)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNoneExistSecureEmail)
        LOG.info("------pass!---------")

    def test_find_username_email_no_exit(self):
        u'''找回用户名邮箱email不存在'''
        LOG.info("------邮箱获取信息用例：start!---------")
        result = self.notification_email('find_username', data.email_not_exist,  data.language,self.user_id, False)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNoneExistSecureEmail)
        LOG.info("------pass!---------")
