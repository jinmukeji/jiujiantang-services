import unittest
from urllib.parse import urljoin

import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from config.read_config import ReadConfig


class ClientAuthor(ApiTestCase):
    u"""提交客户端授权"""

    def setUp(self):
        LOG.info("测试用例开始执行")

    def tearDown(self):
        LOG.info("测试用例执行完毕")

    host = ReadConfig().get_http('url')

    def client(self, client_id, secret_key_hash, seed):
        headers = {"Content-Type": "application/json"}
        url = urljoin(self.host, 'client/auth')
        LOG.info("请求url:%s" % url)
        data = {"client_id": client_id,
                "secret_key_hash": secret_key_hash,
                "seed": seed}
        res = requests.post(url=url, json=data, headers=headers)
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_client_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------提交客户端授权用例：start!---------")
        result = self.client(data.client_id, data.secret_key_hash, data.seed)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_client_error_secret_key_hash(self):
        u"""错误密码"""
        LOG.info("------提交客户端授权用例：start!---------")
        result = self.client(data.client_id, '3+o6Q2y7+g.tzF,U=4qpy7orzd9@(}X8', data.seed)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidSecretKey)
        LOG.info("------pass!--------")

    def test_client_error_client_id(self):
        u"""client_id为空."""
        LOG.info("------提交客户端授权用例：start!---------")
        result = self.client('', data.secret_key_hash, data.seed)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_client_error_secret_key(self):
        u"""secret_key_hash参数为空"""
        LOG.info("------提交客户端授权用例：start!---------")
        result = self.client(data.client_id, '', data.seed)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidSecretKey)
        LOG.info("------pass!---------")

    def test_client_client_id_no_exit(self):
        u"""client_id不存在"""
        LOG.info("------提交客户端授权用例：start!---------")
        result = self.client('123',data.secret_key_hash, data.seed)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")


if __name__ == '__main__':
    unittest.TestCase()
