from urllib.parse import urljoin

import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class Feedback(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def feedback(self, contact_way, content):
        data = {
            'content': content,
            'contact_way': contact_way
        }
        # url = self.host + "/feedback"
        url = urljoin(self.host, 'feedback')
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        req = requests.post(url=url, json=data, headers=Util().get_token())
        return req.json()

    def test_feedback_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.feedback(data.phone, "good")
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_feedback_correct_contact_way(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.feedback(data.email, "good")
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_feedback_error_contact(self):
        u"""错误参数  联系方式不正确,联系方式不作验证"""
        LOG.info("------登录成功用例：start!---------")
        result = self.feedback("137", "good")
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_feedback_error_contact_way(self):
        u"""错误参数  联系方式为空"""
        LOG.info("------登录成功用例：start!---------")
        result = self.feedback("", "good")
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_feedback_error_content(self):
        u"""错误参数  反馈内容为空"""
        LOG.info("------登录成功用例：start!---------")
        result = self.feedback(data.email, "")
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")
