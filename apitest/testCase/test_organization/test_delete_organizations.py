import requests

from common.check_result import ApiTestCase
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class OwnerDeleteOrganization(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def delete_organization(self):
        data = {}
        LOG.info(Util().get_organization_id())
        url = self.host + "owner/organizations/" + Util().get_organization_id() + '/delete'
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        req = requests.post(url=url, json=data, headers=Util().get_token())
        return req.json()

    def test_delete_organization_correct_parameters(self):
        u"""正确参数"""
        LOG.info('------登录成功用例：start!---------')
        result = self.delete_organization()
        LOG.info('获取测试结果：%s' % result)
        self.assertErrorResult(result, const.ErrNotSupportDeleteOrganization)
        LOG.info('------pass!---------')
