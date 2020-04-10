import requests

from common.check_result import ApiTestCase
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class PostUserLanguage(ApiTestCase):
    u"""得到显示语言"""

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def user_language(self):
        data = {}
        url = self.host + "user/" + str(Util().get_user_id()) + "/language"
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        res = requests.get(url=url, json=data, headers=Util().get_token())
        return res.json()

    def test_user_language(self):
        u"""正确参数"""
        LOG.info("------得到显示语言：start!---------")
        result = self.user_language()
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
