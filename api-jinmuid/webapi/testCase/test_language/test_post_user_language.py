import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class PostUserLanguage(ApiTestCase):
    u""" 设置显示语言"""

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def user_language(self, language):
        data = {"language": language}
        url = self.host + "user/" + str(Util().get_user_id()) + "/language"
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        res = requests.post(url=url, json=data, headers=Util().get_token())
        return res.json()

    def test_user_language(self):
        u"""正确参数"""
        LOG.info("------设置显示语言用例：start!---------")
        result = self.user_language(data.language)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_user_language_error_language(self):
        u"""language为空"""
        LOG.info("------设置显示语言用例：start!---------")
        result = self.user_language(" ")
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result,const.ErrInvalidValue)
        LOG.info("------pass!---------")

    def test_user_language_no_exit(self):
        u"""language不存在"""
        LOG.info("------设置显示语言用例：start!---------")
        result = self.user_language("zh")
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result,const.ErrInvalidValue)
        LOG.info("------pass!---------")
