import unittest
from urllib.parse import urljoin

import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class SignIn(ApiTestCase):
    u"""用户登录"""

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')
    nation = Util().get_signin_serial_number()
    mvc = Util().phone_verification_code()

    def sign_in(self, sign_in_method, username, phone, mvc, hashed_password, seed, serial_number, nation_code):
        data = {"sign_in_method": sign_in_method,
                "username": username,
                "phone": phone,
                "mvc": mvc,
                "hashed_password": hashed_password,
                "seed": seed,
                "serial_number": serial_number,
                "nation_code": nation_code
                }
        url = urljoin(self.host, 'signin')
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        return res.json()

    def test_sign_in_phone_mvc(self):
        u"""手机验证码登录"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('phone_mvc', '', data.phone, self.mvc, '', data.seed_phone, self.nation, data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_sign_in_username_password(self):
        u"""用户名密码登录"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('username_password', data.username, '', '',
                              data.hashed_password, data.seed, '', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_sign_in_phone_password(self):
        u"""手机密码登录"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('phone_password', '', data.phone, '',
                              data.password, 'a123', '', '+86')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_sign_in_required_username(self):
        u"""‘’用户名密码登录必填''"""
        LOG.info("------登录成功用例：start!---------")
        data = {"sign_in_method": 'username_password',
                "username": '31',
                "hashed_password": '9626c7444717aab7a3bbdd509bcafa35a7491e9478d421b38e539a621f695edd',
                "seed": ''
                }
        url = urljoin(self.host, 'signin')
        result = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("获取测试结果：%s" % result.json())
        self.assertOkResult(result.json())
        LOG.info("------pass!---------")

    def test_sign_in_required_mvc(self):
        u"""手机验证码登录必填"""
        LOG.info("------登录成功用例：start!---------")
        data = {"sign_in_method": 'phone_mvc',
                "phone": '13221058643',
                "mvc": self.mvc,
                "seed": 'a123',
                'serial_number': self.nation,
                "nation_code": '+86'
                }
        url = urljoin(self.host, 'signin')
        result = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("获取测试结果：%s" % result.json())
        self.assertErrorResult(result.json(),const.InvalidVcRecord)
        LOG.info("------pass!---------")

    def test_sign_in_required_phone(self):
        u"""手机密码登录必填"""
        LOG.info("------登录成功用例：start!---------")
        data = {"sign_in_method": 'phone_password',
                "phone": '13221058643',
                "hashed_password": '501b5144c6bb247e820948b097a8ce41be6e164798519fc67a410c4fd74346c1',
                "seed": 'a123',
                "nation_code": '+86'
                }
        url = urljoin(self.host, 'signin')
        result = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("获取测试结果：%s" % result.json())
        self.assertOkResult(result.json())
        LOG.info("------pass!---------")

    def test_sign_in_error_sign_in_method(self):
        u"""sign_in为空"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('', '31', '', '', data.hashed_password,
                              data.seed, '', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrUsernamePasswordNotMatch)
        LOG.info("------pass!---------")

    def test_sign_in_error_username(self):
        u"""username为空"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('username_password', '', '', '',
                               data.password, data.seed, '', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNonexistentPassword)
        LOG.info("------pass!---------")

    def test_sign_in_error_phone(self):
        u"""phone为空"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('phone_password', '', '', '',
                              data.password, data.seed_phone, '', data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNonexistentPhone)
        LOG.info("------pass!---------")

    def test_sign_in_error_mvc(self):
        u"""mvc为空"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('phone_mvc', '', data.phone, '', '', data.seed_phone, '', data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.InvalidVcRecord)
        LOG.info("------pass!---------")

    def test_sign_in_error_hashed_password(self):
        u"""hashed_password为空"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('username_password', data.username, '', '', '', data.seed_phone, '', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrUsernamePasswordNotMatch)
        LOG.info("------pass!---------")

    def test_sign_in_error_seed(self):
        u"""seed为空"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('username_password', data.username, '', '',
                              data.hashed_password, '', '', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_sign_in_error_serial_number(self):
        u"""序列号错误"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('phone_mvc', '', data.phone, '123456', '', data.seed_phone, '', data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result,const.InvalidVcRecord)
        LOG.info("------pass!---------")

    def test_sign_in_error_nation_code(self):
        u"""区号错误"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('phone_password', '', data.phone, '',
                              data.password, data.seed_phone, '', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result,const.ErrNationCode)
        LOG.info("------pass!---------")

    def test_sign_inc_phone_not_correspond_nation_code(self):
        u"""手机与验证码对应不上"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('phone_mvc', '', data.phone, '123456', '', data.seed_phone, self.nation, "+1")
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result,const.ErrNonexistentPhone)
        LOG.info("------pass!---------")

    def test_sign_in_users_no_exit(self):
        """用户不存在"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('username_password', 'aa', '', '',
                              '2a2ae8e0d4be04d08802f4f1dbeb2bc151eda934555af8004440a7a950ce8dda', 'a123', '', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNonexistentUsername)
        LOG.info("------pass!---------")

    def test_sign_in_error_phone_format(self):
        """手机格式不正确"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('phone_password', '', '132210', '',
                              '2a2ae8e0d4be04d08802f4f1dbeb2bc151eda934555af8004440a7a950ce8dda', 'a123', '', data.nation_code)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNonexistentPhone)
        LOG.info("------pass!---------")

    def test_sign_in_olduser(self):
        """老用户登录"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('username_password', data.username, '', '',
                              data.hashed_password, '', '', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_sign_in_olduser_error(self):
        """老用户密码错误"""
        LOG.info("------登录成功用例：start!---------")
        result = self.sign_in('username_password', data.username, '', '', '666862446769', '', '', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrUsernamePasswordNotMatch)
        LOG.info("------pass!---------")

    def test_sign_in_required_phone_not_set_password(self):
        u"""手机密码登录必填"""
        LOG.info("------登录成功用例：start!---------")
        data = {"sign_in_method": 'phone_password',
                "phone": '13221058643',
                "hashed_password": '',
                "seed": 'a123',
                "nation_code": '+86'
                }
        url = urljoin(self.host, 'signin')
        result = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("获取测试结果：%s" % result.json())
        self.assertErrorResult(result, const.ErrPhonePasswordNotMatch)


if __name__ == '__main__':
    unittest.TestCase()
