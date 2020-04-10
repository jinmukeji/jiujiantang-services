import unittest
from urllib.parse import urljoin

import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from config.read_config import ReadConfig


class Faq(ApiTestCase):
    app_version = '1.10.0'

    def setUp(self):
        LOG.info("测试用例开始执行")

    def tearDown(self):
        LOG.info("测试用例执行完毕")

    host = ReadConfig().get_http('url')

    def get_url(self, client_id, app_version, mobile_type):
        headers = {"Content-Type": "application/json"}
        # url = self.host + "/faq/getUrl"
        url = urljoin(self.host, 'res/getUrl')
        data = {
            "client_id": client_id,
            "app_version": app_version,
            "mobile_type": mobile_type
        }
        res = requests.post(url=url, json=data, headers=headers)
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_get_url_correct_parameters_ios(self):
        u"""ios正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.get_url(data.client_id, self.app_version, 1)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_get_url_correct_parameters_android(self):
        u"""android正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.get_url(data.client_id, self.app_version, 0)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_get_url_error_client_id(self):
        u"""错误参数 ，客户端id为空"""
        LOG.info("------登录成功用例：start!---------")
        result = self.get_url('', self.app_version, 1)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_get_url_error_app_version(self):
        u"""错误参数 ，app版本为空"""
        LOG.info("------登录成功用例：start!---------")
        result = self.get_url(data.client_id, '', 1)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")


if __name__ == '__main__':
    unittest.TestCase()
