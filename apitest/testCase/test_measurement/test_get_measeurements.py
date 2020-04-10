import unittest

import requests

from common.check_result import ApiTestCase
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class OwnerGetMeasurements(ApiTestCase):
    u"""查看测量历史记录"""

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def get_measurement(self, start, end, offset, size, user_id):
        data = {
            "start": start,
            "end": end,
            "offset": offset,
            "size": size,
            "user_id": user_id
        }
        url = self.host + "owner/measurements"
        # url = urljoin(self.host,'owner/measurements')
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        res = requests.get(url=url, params=data, headers=Util().get_token())
        return res.json()

    def test_get_measurement_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.get_measurement('2019-07-24T00:00:00Z', '2019-8-30T17:59:59Z', 0, 100, Util().get_user_id())
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_get_measurement_error_start(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.get_measurement('', '2018-11-10T08:01:27Z ', 0, 10, 76091)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_get_measurement_error_end(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.get_measurement('2018-09-01T08:01:27Z ', '', 0, 10, 76091)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_get_measurement_error_offset(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.get_measurement('2018-09-01T08:01:27Z ', '2018-11-10T08:01:27Z ', '', 10, 76091)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_get_measurement_error_size(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.get_measurement('2018-09-01T08:01:27Z ', '2018-11-10T08:01:27Z ', 0, '', 76091)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_get_measurement_error_user_id(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.get_measurement('2018-09-01T08:01:27Z ', '2018-11-10T08:01:27Z ', 0, 10, '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_get_measurement_error_parameters(self):
        u"""错误参数，结束时间小于开始时间"""
        LOG.info("------登录成功用例：start!---------")
        result = self.get_measurement('2018-09- 10T09:01:27Z', '2018-11-10T08:01:27Z', 0, 10, 76091)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")


if __name__ == '__main__':
    unittest.TestCase()
