import requests

from common.check_result import ApiTestCase
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserName(ApiTestCase):
    u'''用户设置用户名'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def username(self, username):
        data = {"username": username}
        url = self.host + 'user/' + str(Util().get_user_id()) + '/username'
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_username_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------用户设置用户名：start!---------")
        result = self.username('june')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_username_error_username(self):
        u"""username为空"""
        LOG.info("------用户设置用户名：start!---------")
        result = self.username('')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info("------pass!---------")

    def test_username_exit(self):
        u"""username已存在"""
        LOG.info("------用户设置用户名：start!---------")
        result = self.username('june')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info("------pass!---------")
