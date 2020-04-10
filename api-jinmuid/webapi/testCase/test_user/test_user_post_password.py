import time

import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserPasswordPost(ApiTestCase):
    u'''用户设置密码'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def user_password(self, plain_password):
        data = {"plain_password": plain_password}
        url = self.host + 'user/' + str(Util().get_user_id_set_password()) + '/password'
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_token_set_password())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_user_password_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------用户设置密码：start!---------")
        result = self.user_password(data.pwd)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_user_password_error_password(self):
        u"""password错误参数"""
        LOG.info("------用户设置密码：start!---------")
        result = self.user_password('')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrEmptyPassword)
        LOG.info("------pass!---------")

    def test_user_password_length(self):
        u"""密码长度小于8位"""
        LOG.info("------用户设置密码：start!---------")
        result = self.user_password('123')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrWrongFormatOfPassword)
        LOG.info("------pass!---------")
        u'''密码长度大于8位'''
        LOG.info("------用户设置密码：start!---------")
        result = self.user_password('12313413412341234213413333333')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrWrongFormatOfPassword)
        LOG.info("------pass!---------")

    def test_user_password_is_set(self):
        u"""密码长度已被卧设置 """
        LOG.info("------用户设置密码：start!---------")
        result = self.user_password(data.pwd)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrExistPassword)
        LOG.info("------pass!---------")

