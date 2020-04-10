from urllib.parse import urljoin

import requests

from common.check_result import ApiTestCase
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class SignOut(ApiTestCase):
    u"""用户登出"""

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def sign_out(self):
        data = {}
        url = urljoin(self.host, 'user/signout')
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        res = requests.post(url=url, json=data, headers=Util().get_token())
        return res.json()

    def test_sign_out(self):
        LOG.info("------用户登出：start!---------")
        result = self.sign_out()
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
