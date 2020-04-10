from urllib.parse import urljoin

import requests

from common.check_result import ApiTestCase
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class Tips(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def tip(self):
        data = {}
        # url = self.host + "/tips"
        url = urljoin(self.host, 'tips')
        LOG.info("请求url:%s" % url)
        headers = Util().get_token()
        res = requests.get(url=url, data=data, headers=headers)
        return res.json()

    def test_tip_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.tip()
        LOG.info("获取测试结果：%s" % result)
        self.assertEqual(result['ok'], True)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
