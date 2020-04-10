import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserPutProfile(ApiTestCase):
    u"""修改个人档案"""

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def put_profile(self, nickname, birthday, gender, height, weight):
        data = {"nickname": nickname,
                "birthday": birthday,
                "gender": gender,
                "height": height,
                "weight": weight
                }
        url = self.host + 'user/' + str(Util().get_user_id()) + '/profile'
        LOG.info("请求url:%s" % url)
        res = requests.put(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_put_profile_correct_parameters_boy(self):
        u"""正确参数男"""
        LOG.info("------修改个人档案：start!---------")
        result = self.put_profile(data.nickname, '2018-10-15T02:41:31Z', 0, data.height, data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_put_profile_correct_parameters_girl(self):
        u"""正确参数女"""
        LOG.info("------修改个人档案：start!---------")
        result = self.put_profile(data.nickname, '2018-10-15T02:41:31Z', 1, data.height, data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_put_profile_error_nickname(self):
        u"""nickname错误参数"""
        LOG.info("------修改个人档案：start!---------")
        result = self.put_profile('', '2018-10-15T02:41:31Z', 0, data.height, data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info("------pass!---------")

    def test_put_profile_error_birthday(self):
        u"""birthday错误参数"""
        LOG.info("------修改个人档案：start!---------")
        result = self.put_profile(data.nickname, '2018-10', 0, data.height, data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_put_profile_error_gender(self):
        u"""gender错误参数"""
        LOG.info("------修改个人档案：start!---------")
        result = self.put_profile(data.nickname, '2018-10-15T02:41:31Z', 8, data.height, data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info("------pass!---------")

    def test_put_profile_error_height(self):
        u"""height错误参数"""
        LOG.info("------修改个人档案：start!---------")
        result = self.put_profile(data.nickname, '2018-10-15T02:41:31Z', 0, '', data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_put_profile_height_zero(self):
        u"""height等于0"""
        LOG.info("------修改个人档案：start!---------")
        result = self.put_profile(data.nickname, '2018-10-15T02:41:31Z', 0, 0, data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_put_profile_height_max(self):
        u"""height超过最大限制"""
        LOG.info("------修改个人档案：start!---------")
        result = self.put_profile(data.nickname, '2018-10-15T02:41:31Z', 0, 300, data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_put_profile_error_weight(self):
        u"""weight错误参数"""
        LOG.info("------修改个人档案：start!---------")
        result = self.put_profile(data.nickname, '2018-10-15T02:41:31Z', 0, data.height, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_put_profile_weight_zero(self):
        u"""weight等于0"""
        LOG.info("------修改个人档案：start!---------")
        result = self.put_profile(data.nickname, '2018-10-15T02:41:31Z', 0, data.height, 0)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_put_profile_weight_max(self):
        u"""weight超过最大限制"""
        LOG.info("------修改个人档案：start!---------")
        result = self.put_profile(data.nickname, '2018-10-15T02:41:31Z', 0, data.height, 550)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_put_profile_nickname_max(self):
        u"""nickname超过最大限制"""
        LOG.info("------修改个人档案：start!---------")
        result = self.put_profile('smileaefdasfefqadfadefweqfdwqafew', '2018-10-15T02:41:31Z', 0, data.height, data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")
