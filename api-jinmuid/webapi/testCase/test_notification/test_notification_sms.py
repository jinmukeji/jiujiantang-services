import time

import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class NotificationSMS(ApiTestCase):
    u'''手机号获取信息'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def notification_sms(self, sms_Notification_type, phone, language, nation_code):
        url = self.host + "notification/sms"
        LOG.info("请求url:%s" % url)
        data = {"sms_Notification_type": sms_Notification_type,
                "phone": phone,
                "language": language,
                "nation_code": nation_code  # 国家区号
                }
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_notification_sms_signin(self):
        u"""登录短信验证"""
        LOG.info("------手机号获取信息用例：start!---------")
        url = self.host + "notification/sms"
        LOG.info("请求url:%s" % url)
        data = {"sms_Notification_type": "sign_in",
                "phone":"13221058643",
                "language": "zh-Hans",
                "nation_code": "+86"  # 国家区号
                }
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("获取测试结果：%s" % res.json())
        self.assertOkResult(res.json())
        LOG.info("------pass!---------")

    def test_notification_sms_signup(self):
        u"""注册短信验证"""
        LOG.info("------手机号获取信息用例：start!---------")
        result = self.notification_sms('sign_up', data.phone, data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_notification_sms_reset_password(self):
        u"""重置密码短信验证"""
        LOG.info("------手机号获取信息用例：start!---------")
        result = self.notification_sms('reset_password',  data.phone, data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_notification_sms_modify_phone_number(self):
        u"""修改手机号码"""
        LOG.info("------手机号获取信息用例：start!---------")
        result = self.notification_sms('modify_phone',  data.phone, data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_notification_sms_set_phone_number(self):
        u"""设置手机号码"""
        LOG.info("------手机号获取信息用例：start!---------")
        result = self.notification_sms('set_phone',  data.phone, data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_sms_error_type(self):
        u'''sms_Notification_type为空'''
        LOG.info("------手机号获取信息用例：start!---------")
        result = self.notification_sms('',  data.phone, data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_sms_error_phone(self):
        u'''phone为空'''
        LOG.info("------手机号获取信息用例：start!---------")
        result = self.notification_sms('sign_in', '', data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_sms_error_phone_format(self):
        u'''phone格式不正确'''
        LOG.info("------手机号获取信息用例：start!---------")
        result = self.notification_sms('sign_in', '13221058', data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_sms_error_language(self):
        u'''语言为空'''
        LOG.info("------手机号获取信息用例：start!---------")
        result = self.notification_sms('sign_in', data.phone, '', data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_sms_error_nation_code(self):
        u'''国家区号为空'''
        LOG.info("------手机号获取信息用例：start!---------")
        result = self.notification_sms('sign_up', data.phone, data.language, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signup_phone_exit(self):
        u'''注册短信验证phone已存在'''
        LOG.info("------手机号获取信息用例：start!---------")
        result = self.notification_sms('sign_up', data.phone_is_exist, data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signin_phone_no_exit(self):
        u'''登录短信phone不存在'''
        LOG.info("------手机号获取信息用例：start!---------")
        result = self.notification_sms('sign_in', data.phone_not_exist, data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_reset_password_phone_no_exit(self):
        u'''重置phone不存在'''
        LOG.info("------手机号获取信息用例：start!---------")
        result = self.notification_sms('reset_password', data.phone_not_exist, data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result,const.ErrInvalidSigninPhone)
        LOG.info("------pass!---------")

    def test_sms_in_1_minutes(self):
        u'''1分钟内再次发送短信，'''
        LOG.info("------手机号获取信息用例：start!---------")
        result = self.notification_sms('sign_up', data.phone, data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
        time.sleep(30)
        LOG.info("------手机号获取信息用例：start!---------")
        resp = self.notification_sms('sign_up', data.phone, data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % resp)
        self.assertErrorResult(result,const.ErrInvalidRequestCountIn1Minute)
        LOG.info("------pass!---------")

    def test_sms_in_2_minutes(self):
        u'''1分钟后再次发送短信，两条短信返回 的值一致'''
        LOG.info("------手机号获取信息用例：start!---------")
        result = self.notification_sms('sign_up', data.phone, data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
        time.sleep(60)
        LOG.info("------手机号获取信息用例：start!---------")
        resp = self.notification_sms('sign_up', data.phone, data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % resp)
        self.assertOkResult(resp)
        self.assertEqual(result,resp)
        LOG.info("------pass!---------")


    def test_sms_more_2_minutes(self):
        u'''2分钟后再次发送短信，两条短信返回 的值一致'''
        LOG.info("------手机号获取信息用例：start!---------")
        result = self.notification_sms('sign_up', data.phone, data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
        time.sleep(120)
        LOG.info("------手机号获取信息用例：start!---------")
        resp = self.notification_sms('sign_up', data.phone, data.language, data.nation_code)
        LOG.info("获取测试结果：%s" % resp)
        self.assertOkResult(resp)
        self.assertNotEqual(result,resp)
        LOG.info("------pass!---------")
