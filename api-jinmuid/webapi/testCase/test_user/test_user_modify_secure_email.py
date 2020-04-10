import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserModifySecureEmail(ApiTestCase):
    u'''修改安全邮箱'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')
    new_serial_number = Util().get_email_modify_secure_email_serial_number()
    new_verification_code = Util().email_verification_code_new()
    old_verification_number = Util().get_verification_number_email_modify_email()

    def modify_secure_email(self, new_email, new_verification_code, new_serial_number, old_email,
                            old_verification_number):
        data = {"new_email": new_email,
                "new_verification_code": new_verification_code,
                "new_serial_number": new_serial_number,
                "old_email": old_email,
                "old_verification_number": old_verification_number
                }
        url = self.host + "user/" + str(Util().get_user_id()) + "/modify_secure_email"
        LOG.info("请求url:%s" % url)
        LOG.info("请求url:%s" % data)
        req = requests.post(url=url, json=data, headers=Util().get_token())
        return req.json()

    def test_modify_secure_email_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------修改安全邮箱：start!---------")
        result = self.modify_secure_email(data.new_email, self.new_verification_code, self.new_serial_number,
                                          data.email, self.old_verification_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_modify_secure_email_new_email_is_null(self):
        u"""新邮箱为空"""
        LOG.info("------修改安全邮箱：start!---------")
        result = self.modify_secure_email('', self.new_verification_code, self.new_serial_number,
                                          data.email, self.old_verification_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNotExistNewSecureEmail)
        LOG.info("------pass!---------")

    def test_modify_secure_email_new_verification_code_is_null(self):
        u"""新验证码为空"""
        LOG.info("------修改安全邮箱：start!---------")
        result = self.modify_secure_email(data.new_email, '', self.new_serial_number,
                                          data.email, self.old_verification_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNotExistNewSecureEmail)
        LOG.info("------pass!---------")

    def test_modify_secure_email_new_serial_number_is_null(self):
        u"""serial_number为空"""
        LOG.info("------修改安全邮箱：start!---------")
        result = self.modify_secure_email(data.new_email, self.new_verification_code, '',
                                          data.email, self.old_verification_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNotExistNewSecureEmail)
        LOG.info("------pass!---------")

    def test_modify_secure_email_old_email_is_null(self):
        u"""old_email为空"""
        LOG.info("------修改安全邮箱：start!---------")
        result = self.modify_secure_email(data.new_email, self.new_verification_code, self.new_serial_number,
                                          '', self.old_verification_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrSecureEmailNotSet)
        LOG.info("------pass!---------")

    def test_modify_secure_email_old_verification_number_is_null(self):
        u"""old_verification_number为空"""
        LOG.info("------修改安全邮箱：start!---------")
        result = self.modify_secure_email(data.new_email, self.new_verification_code, self.new_serial_number,
                                          data.email, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrSecureEmailNotSet)
        LOG.info("------pass!---------")

    def test_modify_secure_email_old_email_format(self):
        u"""old_email格式不正确"""
        LOG.info("------修改安全邮箱：start!---------")
        result = self.modify_secure_email(data.new_email, self.new_verification_code, self.new_serial_number,
                                          'cuimin', self.old_verification_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNotExistNewSecureEmail)
        LOG.info("------pass!---------")

    def test_modify_secure_email_old_email_not_exist(self):
        u"""old_email不存在"""
        LOG.info("------修改安全邮箱：start!---------")
        result = self.modify_secure_email('601224464@qq.com', self.new_verification_code, self.new_serial_number,
                                          '2270262765@qq.com', self.old_verification_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNotExistNewSecureEmail)
        LOG.info("------pass!---------")

    def test_modify_secure_email_new_email_format(self):
        u"""new_email格式不正确"""
        LOG.info("------修改安全邮箱：start!---------")
        result = self.modify_secure_email('60122446', self.new_verification_code, self.new_serial_number,
                                          data.email, self.old_verification_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNotExistNewSecureEmail)

    def test_modify_secure_email_new_email_is_exist(self):
        u"""new_email已存在 """
        LOG.info("------修改安全邮箱：start!---------")
        result = self.modify_secure_email('nichanglan@jinmuhealth.com', self.new_verification_code,
                                          self.new_serial_number,
                                          data.email, self.old_verification_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNotExistNewSecureEmail)

    def test_modify_secure_email_same(self):
        u"""新旧密码一致"""
        LOG.info("------修改安全邮箱：start!---------")
        result = self.modify_secure_email( data.email, self.new_verification_code,
                                          self.new_serial_number,
                                          data.email, self.old_verification_number)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNotExistNewSecureEmail)
