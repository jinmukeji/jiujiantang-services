import requests

from common.check_result import ApiTestCase
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UsersDelete(ApiTestCase):

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def users_delete(self, user_id_list):
        data = {"user_id_list": [user_id_list]}
        url = self.host + 'owner/organizations/' + Util().get_organization_id() + '/users/delete'
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_users_delete_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.users_delete(95002)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_users_delete_error_parameters(self):
        u"""传值类型不正确"""
        LOG.info("------登录成功用例：start!---------")
        result = self.users_delete('a')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_users_delete_error_users(self):
        u"""删除不存在的用户 """
        LOG.info("------登录成功用例：start!---------")
        result = self.users_delete(100000000)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
