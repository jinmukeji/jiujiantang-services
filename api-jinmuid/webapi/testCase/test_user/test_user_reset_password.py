import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class ResetPassword(ApiTestCase):
    u'''根据验证码重置密码'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')
    verification_number = Util().get_verification_number()

    def reset_password(self, plain_password, verification_number, verification_type):
        data = {"plain_password": plain_password,
                "verification_number": verification_number,
                "verification_type": verification_type
                }
        url = self.host + 'user/' + str(Util().get_user_id()) + '/reset_password'
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_reset_password_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------根据验证码重置密码：start!---------")
        result = self.reset_password(data.pwd, self.verification_number, 'phone')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_reset_password_error_hashed_password(self):
        u"""hashed_password为空"""
        LOG.info("------根据验证码重置密码：start!---------")
        result = self.reset_password('', self.verification_number, 'phone')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_reset_password_error_vc(self):
        u"""vc为空"""
        LOG.info("------根据验证码重置密码：start!---------")
        result = self.reset_password(data.pwd, '', 'phone')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidVerificationNumber)
        LOG.info("------pass!---------")

    def test_reset_password_error_seed(self):
        u"""seed为空"""
        LOG.info("------根据验证码重置密码：start!---------")
        result = self.reset_password(data.pwd, self.verification_number, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_rest_password_same(self):
        u"""新旧密码一致"""
        LOG.info("------根据验证码重置密码：start!---------")
        result = self.reset_password(data.pwd, self.verification_number, 'phone')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
