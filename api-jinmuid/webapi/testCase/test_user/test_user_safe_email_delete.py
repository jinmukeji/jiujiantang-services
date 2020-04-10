import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserSafeEmailDelete(ApiTestCase):
    u'''用户解除设置安全邮箱'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def safe_email_delete(self, verification_code, serial_number, email):
        data = {"verification_code": verification_code,
                "serial_number": serial_number,
                "email": email}
        url = self.host + 'user/' + str(Util().get_user_id()) + '/safe_email/delete'
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        return res.json()

    serial_number = Util().get_email_unset_secure_email_serial_number()
    verification_code = Util().email_verification_code()

    def test_safe_email_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------用户解除设置安全邮箱：start!---------")
        result = self.safe_email_delete(self.verification_code, self.serial_number, data.email)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_safe_email_mvc_null(self):
        u"""mvc为空"""
        LOG.info("------用户解除设置安全邮箱：start!---------")
        result = self.safe_email_delete('', self.serial_number, data.email)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.InvalidVcRecord)
        LOG.info("------pass!---------")

    def test_safe_email_error_mvc(self):
        u"""mvc验证码错误"""
        LOG.info("------用户解除设置安全邮箱：start!---------")
        result = self.safe_email_delete('123', self.serial_number, data.email)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.InvalidVcRecord)
        LOG.info("------pass!---------")

    def test_safe_email_error_serial_number(self):
        u"""serial_number验证码错误"""
        LOG.info("------用户解除设置安全邮箱：start!---------")
        result = self.safe_email_delete(self.verification_code, '', data.email)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.InvalidVcRecord)
        LOG.info("------pass!---------")

    def test_safe_email_error_email(self):
        u"""email验证码错误"""
        LOG.info("------用户解除设置安全邮箱：start!---------")
        result = self.safe_email_delete(self.verification_code, self.serial_number, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidEmailAddress)
        LOG.info("------pass!---------")

    def test_safe_email_error_email_format(self):
        u"""email格式错误"""
        LOG.info("------用户解除设置安全邮箱：start!---------")
        result = self.safe_email_delete(self.verification_code, self.serial_number, 'cuimin')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.InvalidVcRecord)
        LOG.info("------pass!---------")

    def test_safe_email_is_not_exit(self):
        u"""email已被解缄"""
        LOG.info("------用户解除设置安全邮箱：start!---------")
        result = self.safe_email_delete(self.verification_code, self.serial_number,data.email)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrSecureEmailNotSet)
        LOG.info("------pass!---------")
