import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserFindUsernameByEmail(ApiTestCase):
    u'''根据邮箱找回用户名'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')
    serial_number = Util().get_email_find_username_serial_number()
    verification_code = Util().email_verification_code()

    def find_username_by_email(self, verification_code, serial_number, email):
        data = {
            "verification_code": verification_code,
            "serial_number": serial_number,
            "email": email
        }
        url = self.host + "user/find_username_by_email"
        LOG.info("请求url:%s" % url)
        req = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("请求参数：%s" % data)
        return req.json()

    def test_find_username_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------根据邮箱找回用户名：start!---------")
        result = self.find_username_by_email(self.verification_code, self.serial_number, data.email)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_find_username_verification_code_is_null(self):
        u"""verification_code为空"""
        LOG.info("------根据邮箱找回用户名：start!---------")
        result = self.find_username_by_email('', self.serial_number,data.email)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrExpiredVcRecord)
        LOG.info("------pass!---------")

    def test_find_username_verification_code_is_error(self):
        u"""verification_code错误"""
        LOG.info("------根据邮箱找回用户名：start!---------")
        result = self.find_username_by_email('123456', self.serial_number, data.email)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrExpiredVcRecord)
        LOG.info("------pass!---------")

    def test_find_username_serial_number_is_null(self):
        u"""serial_number为空"""
        LOG.info("------根据邮箱找回用户名：start!---------")
        result = self.find_username_by_email(self.verification_code, '', data.email)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrExpiredVcRecord)
        LOG.info("------pass!---------")

    def test_find_username_email_is_null(self):
        u"""email为空"""
        LOG.info("------根据邮箱找回用户名：start!---------")
        result = self.find_username_by_email(self.verification_code, self.serial_number, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrExpiredVcRecord)
        LOG.info("------pass!---------")

    def test_find_username_email_is_error(self):
        u"""email错误"""
        LOG.info("------根据邮箱找回用户名：start!---------")
        result = self.find_username_by_email(self.verification_code, self.serial_number, '2270262765')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrExpiredVcRecord)
        LOG.info("------pass!---------")

    def test_find_username_email_is_not_eixt(self):
        u"""email不存在"""
        LOG.info("------根据邮箱找回用户名：start!---------")
        result = self.find_username_by_email(self.verification_code, self.serial_number, data.email_not_exist)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrExpiredVcRecord)
        LOG.info("------pass!---------")
