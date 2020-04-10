import unittest
from urllib.parse import urljoin

import requests

from common.check_result import ApiTestCase
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class OwnerGetOrganizations(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def get_organizations(self):
        data = {}
        # url = self.host + 'owner/organizations'
        url = urljoin(self.host, 'owner/organizations')
        LOG.info("请求url:%s" % url)
        res = requests.get(url=url, json=data, headers=Util().get_token())
        return res.json()

    def test_get_organizations_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.get_organizations()
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")


if __name__ == '__main__':
    unittest.TestCase()
