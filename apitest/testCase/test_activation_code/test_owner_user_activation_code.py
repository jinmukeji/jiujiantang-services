import unittest

import requests

from common.check_result import ApiTestCase
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class OwnerUserActivationCode(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def activation_code(self, code):
        data = {"code": code}
        url = self.host + 'owner/user/'+str(Util().get_user_id())+'/activation_code'
        LOG.info("请求url:%s" % url)
        req = requests.post(url=url, json=data, headers=Util().get_token())
        return req.json()

    def test_activation_code_correct_parameters(self):
        u"""正确参数"""
        LOG.info('------登录成功用例：start!---------')
        result = self.activation_code("3CFL6VEWDLRFQYTE")
        LOG.info('获取测试结果：%s' % result)
        self.assertOkResult(result)
        LOG.info('------pass!---------')

    def test_activation_code_error_parameters(self):
        u"""正确参数"""
        LOG.info('------登录成功用例：start!---------')
        result = self.activation_code(" ")
        LOG.info('获取测试结果：%s' % result)
        self.assertErrorResult(result,const.ErrDatabase)
        LOG.info('------pass!---------')


if __name__ == '__main__':
    unittest.TestCase()
