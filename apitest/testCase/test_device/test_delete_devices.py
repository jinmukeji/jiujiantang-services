import json
import unittest

import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


# noinspection PyPep8Naming
class DeleteDevices(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def delete_device(self, mac):
        data = [{'mac': mac}]
        url = self.host + 'owner/organizations/' + Util().get_organization_id() + '/devices/delete'
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        req = requests.post(url=url, data=json.dumps(data), headers=Util().get_token())
        return req.json()

    def test_delete_device_correct_parameters(self):
        u"""正确参数"""
        LOG.info('------登录成功用例：start!---------')
        result = self.delete_device(data.mac)
        LOG.info('获取测试结果：%s' % result)
        self.assertOkResult(result)
        LOG.info('------pass!---------')

    def test_delete_device_error_parameters(self):
        u"""错误参数"""
        LOG.info('------登录成功用例：start!---------')
        result = self.delete_device('38D269ED6')
        LOG.info('获取测试结果：%s' % result)
        self.assertErrorResult(result, const.ErrDeviceNotBelongToOrganization)
        LOG.info('------pass!---------')

    def test_delete_device_error_mac(self):
        u"""mac错误参数"""
        LOG.info('------登录成功用例：start!---------')
        result = self.delete_device('')
        LOG.info('获取测试结果：%s' % result)
        self.assertErrorResult(result, const.ErrInvalidDevice)
        LOG.info('------pass!---------')


if __name__ == '__main__':
    unittest.TestCase()
