import requests

from common.check_result import ApiTestCase
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserSecureQuestionToModify(ApiTestCase):
    u'''修改密保前获取已设置密保列表'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def secure_question_to_modify(self):
        data = {}
        url = self.host + "user/" + str(Util().get_user_id()) + "/secure_question_to_modify"
        LOG.info("请求url:%s" % url)
        req = requests.get(url=url, json=data, headers=Util().get_token())
        return req.json()

    def test_secure_question_to_modify_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------修改密保前获取已设置密保列表：start!---------")
        result = self.secure_question_to_modify()
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
