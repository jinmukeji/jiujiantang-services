import requests

from common.check_result import ApiTestCase
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserToken(ApiTestCase):
    u"""获取最新token"""

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def get_token(self):
        data = {}
        url = self.host + "user/" + str(Util().get_user_id()) + "/token"
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        res = requests.get(url=url, json=data, headers=Util().get_token())
        return res.json()

    def test_get_token_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.get_token()
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
