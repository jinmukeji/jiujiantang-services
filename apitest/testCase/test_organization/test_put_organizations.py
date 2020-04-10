import json

import requests

from common.check_result import ApiTestCase
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class OwnerPutOrganizations(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def put_organizations(self, name, state, city, street, phone, contact, types, email, country):
        data = {"profile": {
            "name": name,
            "state": state,
            "city": city,
            "street": street,
            "phone": phone,
            "contact": contact,
            "type": types,
            "email": email,
            "country": country}}
        url = self.host + "owner/organizations/" + Util().get_organization_id()
        LOG.info("请求url:%s" % url)
        res = requests.put(url=url, data=json.dumps(data), headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_put_organizations_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.put_organizations('smile', '江苏省', '常州市', '天宁区', '13700009009', 'cc', '养生', '2270262765@qq.com',
                                        '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_put_organizations_error_name(self):
        u"""name错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.put_organizations('', '江苏省', '常州市', '天宁区', '13700009009', 'cc', '养生', '2270262765@qq.com', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_put_organizations_error_phone(self):
        u"""phone错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.put_organizations('smile', '江苏省', '常州市', '天宁区', '137000t55', 'cc', '养生', '2270262765@qq.com',
                                        '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info("------pass!---------")

    def test_put_organizations_error_email(self):
        u"""email错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.put_organizations('smile', '江苏省', '常州市', '天宁区', '13700009009', 'cc', '养生', '6706543', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info("------pass!---------")

    def test_put_organizations(self):
        u"""必填项"""
        LOG.info("------登录成功用例：start!---------")
        result = self.put_organizations('smile', '', '', '', '', '', '', '',
                                        '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
