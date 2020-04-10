import requests

from common.check_result import ApiTestCase
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class DeleteMeasurements(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def deletemeasurements(self, record_id_list):
        data = {
            "record_id_list": record_id_list
        }
        url = self.host+'user/measurements/'+str(Util().get_user_id())+'/delete'
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        res = requests.post(url=url, json=data, headers=Util().get_token())
        return res.json()

    def test_delete_measurements_correct_params(self):
        result = self.deletemeasurements([227401, 227344])
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_delete_measurements_correct_blank_params(self):
        result = self.deletemeasurements([])
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_delete_measurements_error_params(self):
        result = self.deletemeasurements([227463])
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, 80000)
        LOG.info("------pass!---------")
