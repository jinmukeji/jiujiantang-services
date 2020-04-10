import unittest

import requests

from common.check_result import ApiTestCase
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class OwnerGetMeasurementsAnalyze(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def analyze(self, cid, analysis_session, question_key, values):
        data = {
            "cid": cid,
            "analysis_session": analysis_session,
            "answers":
                [{
                    "question_key": question_key,
                    "values": values
                }]
        }
        url = self.host + "owner/measurements/" + Util().get_record_id() + "/analyze"
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        res = requests.post(url=url, json=data, headers=Util().get_token())
        return res.json()

    def test_analyze_correct_parameters(self):
        u"""正确参数"""
        LOG.info('------登录成功用例：start!---------')
        result = self.analyze(206839, '', 'Q0001.0', 'QC0002.0')
        LOG.info('获取测试结果：%s' % result)
        self.assertOkResult(result)
        LOG.info('------pass!---------')

    def test_analyze_error_cid(self):
        u"""cid错误参数"""
        LOG.info('------登录成功用例：start!---------')
        result = self.analyze('1000', '', 'Q0001.0', 'QC0002.0')
        LOG.info('获取测试结果：%s' % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info('------pass!---------')

    def test_analyze_error_analysis_session(self):
        u"""analysis_session错误参数"""
        LOG.info('------登录成功用例：start!---------')
        result = self.analyze(206839, 's', 'Q0001', 'QC0002.0')
        LOG.info('获取测试结果：%s' % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info('------pass!---------')

    def test_analyze_error_question_key(self):
        u"""question_key错误参数"""
        LOG.info('------登录成功用例：start!---------')
        result = self.analyze(206839, '', '', 'QC0002.0')
        LOG.info('获取测试结果：%s' % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info('------pass!---------')

    def test_analyze_error_values(self):
        u"""values错误参数"""
        LOG.info('------登录成功用例：start!---------')
        result = self.analyze(206839, '', 'Q0001.0', '')
        LOG.info('获取测试结果：%s' % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info('------pass!---------')

    def test_analyze_error_cid_not_value(self):
        u"""cid不存在"""
        LOG.info('------登录成功用例：start!---------')
        result = self.analyze(2066786, '', 'Q0001.0', '')
        LOG.info('获取测试结果：%s' % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info('------pass!---------')


if __name__ == '__main__':
    unittest.TestCase()
