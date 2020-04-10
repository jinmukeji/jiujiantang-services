import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class Validate(ApiTestCase):
    u'''根据密保问题修改密码前验证手机号码或者用户名是否存在'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def validate_username_or_phone(self, validation_type, username, phone, nation_code):
        data = {
            "validation_type": validation_type,
            "username": username,
            "phone": phone,
            "nation_code": nation_code
        }
        url = self.host + "user/validate_username_or_phone"
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        req = requests.post(url=url, json=data, headers=Util().get_token())
        return req.json()

    def test_user_questions_correct_parameters_username(self):
        u"""用户名正确参数"""
        LOG.info("------根据密保问题修改密码前验证手机号码或者用户名是否存在：start!---------")
        result = self.validate_username_or_phone('username', data.username, '', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_user_questions_correct_parameters_phone(self):
        u"""手机号正确参数"""
        LOG.info("------根据密保问题修改密码前验证手机号码或者用户名是否存在：start!---------")
        result = self.validate_username_or_phone('phone', '', data.phone, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_user_questions_username_type_is_null(self):
        u"""type为空"""
        LOG.info("------根据密保问题修改密码前验证手机号码或者用户名是否存在：start!---------")
        result = self.validate_username_or_phone('', data.username, '', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidSecureQuestionValidationMethod)
        LOG.info("------pass!---------")

    def test_user_questions_phone_type_is_null(self):
        u"""type为空"""
        LOG.info("------根据密保问题修改密码前验证手机号码或者用户名是否存在：start!---------")
        result = self.validate_username_or_phone('', '', data.phone, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidSecureQuestionValidationMethod)
        LOG.info("------pass!---------")

    def test_user_questions_phone_format(self):
        u"""手机格式不正确"""
        LOG.info("------根据密保问题修改密码前验证手机号码或者用户名是否存在：start!---------")
        result = self.validate_username_or_phone('phone', '', '13221058', data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNonexistentPhone)
        LOG.info("------pass!---------")

    def test_user_questions_phone_not_nation_code(self):
        u'''手机格式和区号不一致'''
        LOG.info("------根据密保问题修改密码前验证手机号码或者用户名是否存在：start!---------")
        result = self.validate_username_or_phone('phone', '', data.phone, '+1')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNonexistentPhone)
        LOG.info("------pass!---------")

    def test_user_question_phone_not_exit(self):
        u'''手机号码不存在'''
        LOG.info("------根据密保问题修改密码前验证手机号码或者用户名是否存在：start!---------")
        result = self.validate_username_or_phone('phone', '', data.phone_not_exist, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNonexistentPhone)
        LOG.info("------pass!---------")

    def test_user_question_username_not_exit(self):
        u'''用户名不存在'''
        LOG.info("------根据密保问题修改密码前验证手机号码或者用户名是否存在：start!---------")
        result = self.validate_username_or_phone('username', 'lili', '', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNonexistentUsername)
        LOG.info("------pass!---------")
