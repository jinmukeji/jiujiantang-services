import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class ValidateEmailVerificationCode(ApiTestCase):
    u'''验证邮箱验证码是否正确'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')
    # reset_serial_number = Util().get_email_reset_password_serial_number()
    modify_serial_number = Util().get_email_modify_secure_email_serial_number_old()
    verification_code = Util().email_verification_code()

    def validate_email_verification_code(self, email, verification_code, serial_number, verification_type):
        data = {
            "email": email,
            "verification_code": verification_code,
            "serial_number": serial_number,
            "verification_type": verification_type
        }
        url = self.host + "user/validate_email_verification_code"
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        req = requests.post(url=url, json=data, headers=Util().get_authorization())
        return req.json()

    def test_validate_email_correct_parameters_reset_password(self):
        u"""正确参数reset_password"""
        LOG.info("------验证邮箱验证码是否正确：start!---------")
        result = self.validate_email_verification_code(data.email, self.verification_code,
                                                       self.reset_serial_number, 'reset_password')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_validate_email_correct_parameters_modify_secure_email(self):
        u"""正确参数modify_secure_email"""
        LOG.info("------验证邮箱验证码是否正确：start!---------")
        result = self.validate_email_verification_code(data.email, self.verification_code,
                                                       self.modify_serial_number,
                                                       'modify_secure_email')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_validate_email_is_null(self):
        u"""email为空"""
        LOG.info("------验证邮箱验证码是否正确：start!---------")
        result = self.validate_email_verification_code('', self.verification_code, self.modify_serial_number,
                                                       'modify_secure_email')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrSecureEmailUsedByOthers)
        LOG.info("------pass!---------")

    def test_validate_email_not_exit(self):
        u"""email不存在"""
        LOG.info("------验证邮箱验证码是否正确：start!---------")
        result = self.validate_email_verification_code(data.email_not_exist, self.verification_code,
                                                       self.modify_serial_number,
                                                       'modify_secure_email')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrSecureEmailUsedByOthers)
        LOG.info("------pass!---------")

    def test_validate_email_format(self):
        u"""email不存在"""
        LOG.info("------验证邮箱验证码是否正确：start!---------")
        result = self.validate_email_verification_code('267026276', self.verification_code, self.modify_serial_number,
                                                       'modify_secure_email')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrSecureEmailUsedByOthers)
        LOG.info("------pass!---------")

    def test_validate_email_serial_number_is_null(self):
        u"""serial_number为空"""
        LOG.info("------验证邮箱验证码是否正确：start!---------")
        result = self.validate_email_verification_code(data.email, self.verification_code, '',
                                                       'modify_secure_email')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_validate_email_verification_type_is_null(self):
        u"""verification_type为空"""
        LOG.info("------验证邮箱验证码是否正确：start!---------")
        result = self.validate_email_verification_code(data.email, self.verification_code,
                                                       self.modify_serial_number,
                                                       '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
