import json
import unittest

import requests

from common.check_result import ApiTestCase
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserPreferences(ApiTestCase):
    u'''设置/修改密保'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def user_preferences(self):
        data = {}
        url = self.host + "user/" + str(Util().get_user_id()) + "/preferences"
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        req = requests.get(url=url, data=json.dumps(data), headers=Util().get_token())
        return req.json()

    def test_user_questions_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------设置/修改密保：start!---------")
        result = self.user_preferences()
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")


if __name__ == '__main__':
    unittest.TestCase()
