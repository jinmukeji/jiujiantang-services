import requests

from common.check_result import ApiTestCase
from common.data import data
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserSigninPhone(ApiTestCase):
    u'''设置登录手机号'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')
    serial_number = Util().get_modify_phone_number_serial_number()
    mvc = Util().phone_verification_code()
    verification_code = Util().phone_verification_code()

    def signin_phone(self, phone, nation_code, mvc, serial_number, old_phone, old_nation_code, verification_code):
        data = {
            "phone": phone,
            "nation_code": nation_code,
            "mvc": mvc,
            "serial_number": serial_number,
            "old_phone": old_phone,
            "old_nation_code": old_nation_code,
            "verification_code": verification_code
        }
        url = self.host + 'user/' + str(Util().get_user_id()) + '/signin_phone'
        LOG.info("请求url:%s" % url)
        res = requests.put(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_signin_phone_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------设置登录手机号：start!---------")
        result = self.signin_phone(data.phone, data.nation_code, self.mvc, self.serial_number, data.phone_not_exist,data.nation_code,
                                   self.verification_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signin_phone_phone_is_null(self):
        u"""手机号为空"""
        LOG.info("------设置登录手机号：start!---------")
        result = self.signin_phone('', data.nation_code, '', self.serial_number,data.phone_not_exist, data.nation_code, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signin_phone_phone_not_exit(self):
        u"""手机号格式不正确"""
        LOG.info("------设置登录手机号：start!---------")
        result = self.signin_phone(data.phone, data.nation_code, '', self.serial_number, data.phone_not_exist, data.nation_code, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signin_phone_nation_code_is_null(self):
        u"""区号为空"""
        LOG.info("------设置登录手机号：start!---------")
        result = self.signin_phone(data.phone, '', '', self.serial_number, data.phone_not_exist, data.nation_code, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signin_phone_phone_not_correspond_nation_code(self):
        u"""手机号和区号不对应"""
        LOG.info("------设置登录手机号：start!---------")
        result = self.signin_phone(data.phone, '+1', '', self.serial_number, data.phone_not_exist, data.nation_code, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signin_phone_mvc_is_null(self):
        u"""mvc为空"""
        LOG.info("------设置登录手机号：start!---------")
        result = self.signin_phone(data.phone, data.nation_code, '', self.serial_number, data.phone_not_exist, data.nation_code, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signin_phone_serial_number_is_null(self):
        u"""serial_number为空"""
        LOG.info("------设置登录手机号：start!---------")
        result = self.signin_phone(data.phone, data.nation_code, '', '', data.phone_not_exist, data.nation_code, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signin_phone_old_phone_is_null(self):
        u"""手机号为空"""
        LOG.info("------设置登录手机号：start!---------")
        result = self.signin_phone(data.phone, data.nation_code, '', self.serial_number, '', data.nation_code, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signin_phone_old_phone_not_exit(self):
        u"""手机号格式不正确"""
        LOG.info("------设置登录手机号：start!---------")
        result = self.signin_phone('1322105', data.nation_code, '', self.serial_number, '1306569', data.nation_code, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signin_phone_old_nation_code_is_null(self):
        u"""区号为空"""
        LOG.info("------设置登录手机号：start!---------")
        result = self.signin_phone(data.phone, data.nation_code, '', self.serial_number, '13065697748', '', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signin_phone_verification_code_is_null(self):
        u"""verification_code为空"""
        LOG.info("------设置登录手机号：start!---------")
        result = self.signin_phone(data.phone, data.nation_code, '', self.serial_number, data.phone_not_exist, data.nation_code, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signin_phone_old_phone_not_correspond_nation_code(self):
        u"""旧手机号和区号不对应"""
        LOG.info("------设置登录手机号：start!---------")
        result = self.signin_phone(data.phone, '+1', '', self.serial_number, data.phone_not_exist, '+8', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signin_phone_same(self):
        u"""新旧手机号一致"""
        LOG.info("------设置登录手机号：start!---------")
        result = self.signin_phone(data.phone_not_exist, data.nation_code, '', self.serial_number, data.phone_not_exist,data.nation_code, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
