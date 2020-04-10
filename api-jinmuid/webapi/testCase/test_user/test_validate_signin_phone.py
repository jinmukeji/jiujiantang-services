import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class ValidateSigninPhone(ApiTestCase):
    u'''验证登录手机号码'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')
    serial_number = Util().get_reset_password_serial_number()
    mvc = Util().phone_verification_code()

    def validate_signin_phone(self, phone, nation_code, mvc, serial_number):
        data = {
            "phone": phone,
            "nation_code": nation_code,
            "mvc": mvc,
            "serial_number": serial_number
        }
        url = self.host + 'validate_signin_phone'
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_signin_phone_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------验证登录手机号码：start!---------")
        result = self.validate_signin_phone(data.phone, data.nation_code, self.mvc, self.serial_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signin_phone_phone_is_null(self):
        u"""手机号为空"""
        LOG.info("------验证登录手机号码：start!---------")
        result = self.validate_signin_phone('', data.nation_code, self.mvc, self.serial_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInValidMVC)
        LOG.info("------pass!---------")

    def test_signin_phone_phone_not_exit(self):
        u"""手机号格式不正确"""
        LOG.info("------验证登录手机号码：start!---------")
        result = self.validate_signin_phone('13221058', data.nation_code, self.mvc, self.serial_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInValidMVC)
        LOG.info("------pass!---------")

    def test_signin_phone_nation_code_is_null(self):
        u"""区号为空"""
        LOG.info("------验证登录手机号码：start!---------")
        result = self.validate_signin_phone(data.phone, '', self.mvc, self.serial_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInValidMVC)
        LOG.info("------pass!---------")

    def test_signin_phone_mvc_is_null(self):
        u"""mvc为空"""
        LOG.info("------验证登录手机号码：start!---------")
        result = self.validate_signin_phone(data.phone, data.nation_code, '', self.serial_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInValidMVC)
        LOG.info("------pass!---------")

    def test_signin_phone_serial_number_is_null(self):
        u"""serial_number为空"""
        LOG.info("------验证登录手机号码：start!---------")
        result = self.validate_signin_phone(data.phone, data.nation_code, self.mvc, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInValidMVC)
        LOG.info("------pass!---------")

    def test_signin_phone_phone_not_correspond_nation_code(self):
        u"""手机号和区号不对应"""
        LOG.info("------验证登录手机号码：start!---------")
        result = self.validate_signin_phone(data.phone, '+1', self.mvc, self.serial_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInValidMVC)
        LOG.info("------pass!---------")
