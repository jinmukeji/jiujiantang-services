import unittest

import requests

from common.check_result import ApiTestCase
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserSecureQuestionList(ApiTestCase):
    u'''获取所有密保问题列表'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def secure_question_list(self):
        data = {}
        url = self.host + "user/" + str(Util().get_user_id()) + "/secure_question_list"
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        req = requests.get(url=url, json=data, headers=Util().get_token())
        return req.json()

    def test_user_questions_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------获取所有密保问题列表：start!---------")
        result = self.secure_question_list()
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")


if __name__ == '__main__':
    unittest.TestCase()
