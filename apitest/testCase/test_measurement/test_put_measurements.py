import json

import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class PutMeasurements(ApiTestCase):
    u"""修改测量备注"""

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def put_measurements(self, user_id, remark):
        data = {
            "user_id": user_id,
            "remark": remark
        }
        url = self.host + 'owner/measurements/' + Util().get_record_id() + '/remark'
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        response = requests.put(url, data=json.dumps(data), headers=Util().get_token())
        return response.json()

    def test_put_organizations_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.put_measurements(data.user_id, 'test')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_put_organizations_error_user_id(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.put_measurements('', 'test')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_put_organizations_error_remark(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.put_measurements(data.user_id, 'test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213te\
                                              test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213test1213st1213test1213test1213test1213test1213test1213test1213test1213test1213test1213')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNoPermissionSubmitRemark)
        LOG.info("------pass!---------")
