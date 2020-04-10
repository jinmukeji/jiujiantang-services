import unittest

import requests

from common.check_result import ApiTestCase
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserDevice(ApiTestCase):
    u''' 用户得到使用的device'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def user_device(self):
        data = {}
        url = self.host + "user/" + str(Util().get_user_id()) + "/device"
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        req = requests.get(url=url, data=data, headers=Util().get_token())
        return req.json()

    def test_user_device_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------用户得到使用的device'：start!---------")
        result = self.user_device()
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")


if __name__ == '__main__':
    unittest.TestCase()
