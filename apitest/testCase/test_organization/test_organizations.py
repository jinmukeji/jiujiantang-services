import unittest
from urllib.parse import urljoin

import requests

from common.check_result import ApiTestCase
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class OwnerOrganizations(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def organizations(self, name, state, city, street, phone, contact, types, email, country):
        data = {
            "profile":
                {"name": name,
                 "state": state,
                 "city": city,
                 "street": street,
                 "phone": phone,
                 "contact": contact,
                 "type": types,
                 "email": email,
                 "country": country
                 }
        }
        # url = self.host + "owner/organizations"
        url = urljoin(self.host, 'owner/organizations')
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_organizations_correct_required_parameters(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        self.delete_organization()
        result = self.organizations('张三')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult_organization(result)
        LOG.info("------pass!---------")

    def test_organizations_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.organizations('张三', '江苏省', '常州市', '天宁区', '13700009009', 'cc', '养生', '2270262765@qq.com',
                                    '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult_organization(result)
        LOG.info("------pass!---------")

    def test_organizations_error_name(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.organizations('', '江苏省', '常州市', '天宁区', '13700009009', 'cc', '养生', '2270262765@qq.com', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_organizations_error_state(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.organizations('smile', '', '常州市', '天宁区', '13700009009', 'cc', '养生', '2270262765@qq.com', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_organizations_error_city(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.organizations('smile', '江苏省', '', '天宁区', '13700009009', 'cc', '养生', '2270262765@qq.com', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_organizations_error_street(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.organizations('smile', '江苏省', '常州市', '', '13700009009', 'cc', '养生', '2270262765@qq.com', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_organizations_error_phone(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.organizations('smile', '江苏省', '常州市', '天宁区', '', 'cc', '养生', '2270262765@qq.com', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_organizations_error_contact(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.organizations('smile', '江苏省', '常州市', '天宁区', '13700009009', '', '养生', '2270262765@qq.com', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_organizations_error_type(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.organizations('smile', '江苏省', '常州市', '天宁区', '13700009009', 'cc', '', '2270262765@qq.com', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_organizations_error_email(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.organizations('smile', '江苏省', '常州市', '天宁区', '13700009009', 'cc', '养生', '', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")

    def test_organizations_error_country(self):
        u"""错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.organizations('smile', '江苏省', '常州市', '天宁区', '13700009009', 'cc', '养生', '2270262765@qq.com', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRPCInternal)
        LOG.info("------pass!---------")


if __name__ == '__main__':
    unittest.TestCase()
