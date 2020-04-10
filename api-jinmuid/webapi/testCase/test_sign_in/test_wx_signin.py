from urllib.parse import urljoin

import requests

from common.check_result import ApiTestCase
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class WxSignIn(ApiTestCase):
    u"""获取登录二维码"""

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def signin_qrcode(self):
        data = {}
        url = urljoin(self.host, 'wx/signin/qrcode')
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        return res.json()

    def test_signin_qrcode_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.signin_qrcode()
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
