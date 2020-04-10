import requests

from common.check_result import ApiTestCase
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class WxUserConnectDelete(ApiTestCase):
    u'''获取解绑微信二维码'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def connect_qrcode_delete(self):
        data = {}
        url = self.host + 'wx/user/' + str(Util().get_user_id()) + '/connect/qrcode'
        LOG.info("请求url:%s" % url)
        res = requests.delete(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_connect_qrcode_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------获取解绑微信二维码：start!---------")
        result = self.connect_qrcode_delete()
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
