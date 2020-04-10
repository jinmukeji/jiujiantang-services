import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserRegion(ApiTestCase):
    u'''选择区域'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def user_region(self, user_id, region):
        data = {"user_id": user_id,
                "region": region
                }
        url = self.host + 'user/region'
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_user_region_correct_parameters_china(self):
        u"""正确参数（中国大陆含港澳台）"""
        LOG.info("------选择区域：start!---------")
        result = self.user_region(data.user_id, 0)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_user_region_correct_parameters_taiwan(self):
        u"""正确参数（中国台湾）"""
        LOG.info("------选择区域：start!---------")
        result = self.user_region(data.user_id, 1)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrExsitRegion)
        LOG.info("------pass!---------")

    def test_user_region_correct_parameters_foreign(self):
        u"""正确参数（国外）"""
        LOG.info("------选择区域：start!---------")
        result = self.user_region(data.user_id, 2)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrExsitRegion)
        LOG.info("------pass!---------")

    def test_user_region_user_id_is_null(self):
        u"""user_id为空"""
        LOG.info("------选择区域：start!---------")
        result = self.user_region('', 0)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_user_region_user_id_not_signin(self):
        u"""user_id未授权登录"""
        LOG.info("------选择区域：start!---------")
        result = self.user_region(data.user_id, 0)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidUser)
        LOG.info("------pass!---------")

    def test_user_region_region_is_null(self):
        u"""region为空"""
        LOG.info("------选择区域：start!---------")
        result = self.user_region(data.user_id, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_user_region_region_not_exit(self):
        u"""region不存在"""
        LOG.info("------选择区域：start!---------")
        result = self.user_region(data.user_id, 9999)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrExsitRegion)
        LOG.info("------pass!---------")
