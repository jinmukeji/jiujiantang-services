import unittest
from urllib.parse import urljoin

import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserSignIn(ApiTestCase):
    u"""用户登录"""

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def sign_in(self, sign_in_key, register_type, password_hash):
        data = {
            "sign_in_key": sign_in_key,
            "register_type": register_type,
            "password_hash": password_hash
        }
        url = urljoin(self.host, 'users/signin')
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        return res.json()

    def test_sign_in_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in(
            data.sign_in_key,
            'username',
            data.password_hash
        )
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_sign_in_error_parameters(self):
        u"""‘’错误参数格式''"""
        LOG.info("------登录成功用例：start!---------")
        sign_in_key = data.sign_in_key
        LOG.info("输入sign_in_key：%s" % sign_in_key)
        register_type = 'username'
        LOG.info("输入register_type：%s" % register_type)
        password_hash = data.password_hash
        LOG.info("输入password_hash：%s" % password_hash)
        result = self.sign_in(sign_in_key, register_type, password_hash)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_sign_in_error_sign_in_key(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in(
            '',
            'username',
            data.password_hash)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    # 注释原因：判断此字段为空影响当前流程
    # def test_sign_in_error_register_type(self):
    #     u"""错误参数"""
    #     LOG.info("------登录成功用例：start!---------")
    #     result = self.sign_in('1705010001', '',
    #                           'ab740d6fa580380614a2af88534f5536bea55053a42678a8f8d1e2db24c077aa')
    #     LOG.info("获取测试结果：%s" % result)
    #     self.assertErrorResult(result, const.ErrIncorrectRegisterType)
    #     LOG.info("------pass!---------")

    def test_sign_in_error_password_hash(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in(data.sign_in_key, 'username', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrUsernamePasswordNotMatch)
        LOG.info("------pass!---------")

    def test_sign_in_new_user(self):
        u"""登云正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('dengyun-10001', 'username',
                              'a78cb958be04a74210684038615b2a37eb46ae935ae16041837d2b2376837481')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

        u'''康美正确参数'''
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('kangmei-10001', 'username',
                              '122fbd67b9af5f02a05b7dddb977f2d14d7957a4a5df8cd9435f8cb2b583a9e3')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_sign_in_users_no_exit(self):
        """用户不存在"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('13700001001', 'username',
                              'ab740d6fa580380614a2af88534f5536bea55053a42678a8f8d1e2db24c077aa')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")



if __name__ == '__main__':
    unittest.TestCase()
