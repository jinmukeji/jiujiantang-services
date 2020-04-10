from urllib.parse import urljoin

import requests

from common.check_result import ApiTestCase
from common.log import LOG
from config.read_config import ReadConfig


class Version(ApiTestCase):
    u'''查看版本信息'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def version(self):
        data = {}
        url = urljoin(self.host, "version")
        LOG.info("请求url:%s" % url)
        req = requests.get(url=url, data=data)
        LOG.info(req.json())
        return req.json()

    def test_version_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------查看版本信息：start!---------")
        result = self.version()
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
