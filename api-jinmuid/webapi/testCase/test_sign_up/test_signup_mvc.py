import unittest

import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class SignupMvc(ApiTestCase):
    u'''使用验证码注册'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')
    number = Util().get_verification_number_signup()

    def signup_mvc(self, phone, nation_code, verification_number, plain_password,nickname,birthday,gender,height,weight):
        data = {
            "phone": phone,
            "nation_code": nation_code,
            "verification_number": verification_number,
            "plain_password": plain_password,
            "user_profile": {
                "nickname": nickname,
                "birthday": birthday,
                "gender": gender,
                "height":height,
                "weight": weight
            },
            "language": "zh-Hans"
        }
        url = self.host + "signup/verification_number"
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        req = requests.post(url=url, json=data, headers=Util().get_authorization())
        return req.json()

    def test_signup_mvc_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------使用验证码注册：start!---------")
        result = self.signup_mvc(data.phone, data.nation_code, self.number, "a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signup_mvc_error_mvc(self):
        u"""mvc错误"""
        LOG.info("------使用验证码注册：start!---------")
        result = self.signup_mvc(data.phone, data.nation_code, self.number,"a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result,const.InvalidVcRecord)
        LOG.info("------pass!---------")

    def test_signup_mvc_error_phone(self):
        u"""phone错误"""
        LOG.info("------使用验证码注册：start!---------")
        result = self.signup_mvc('', data.nation_code, self.number, "a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result)
        LOG.info("------pass!---------")

    def test_signup_mvc_error_serial_number(self):
        u"""serial_number错误"""
        LOG.info("------使用验证码注册：start!---------")
        result = self.signup_mvc(data.phone, data.nation_code, '', "a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrExsitRegisteredPhone)
        LOG.info("------pass!---------")

    def test_signup_mvc_error_nation_code(self):
        u"""nation_code错误"""
        LOG.info("------使用验证码注册：start!---------")
        result = self.signup_mvc(data.phone, '', self.number, "a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result,const.ErrNationCode)
        LOG.info("------pass!---------")

    def test_signup_mvc_phone_format_error(self):
        u"""手机格式错误"""
        LOG.info("------使用验证码注册：start!---------")
        result = self.signup_mvc('1322143', data.nation_code, self.number,"a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result,const.InvalidVcRecord)
        LOG.info("------pass!---------")

    def test_signup_mvc_phone_not_correspond_nation_code(self):
        u"""手机与验证码对应不上"""
        LOG.info("------使用验证码注册：start!---------")
        verification_number_usa=Util().get_verification_number_signup_usa()
        result = self.signup_mvc(data.phone, '+1', verification_number_usa,"a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result,const.InvalidVcRecord)
        LOG.info("------pass!---------")

    def test_singup_mvc_phone_is_exit(self):
        u"""手机已注册"""
        LOG.info("------使用验证码注册：start!---------")
        result = self.signup_mvc('13700007474', '+1', self.number,"a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.InvalidVcRecord)
        LOG.info("------pass!---------")

    def test_signup_mvc_phone_usa(self):
        u"""美国手机号"""
        LOG.info("------使用验证码注册：start!---------")
        verification_number_usa=Util().get_verification_number_signup_usa()
        result = self.signup_mvc('7022682192', '+1', verification_number_usa,"a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")


    def test_signup_mvc_phone_HK(self):
        u"""香港手机号"""
        LOG.info("------使用验证码注册：start!---------")
        verification_number_usa=Util().get_verification_number_signup_usa()
        result = self.signup_mvc('7022682192', '+852', verification_number_usa,"a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signup_mvc_phone_AOMEN(self):
        u"""澳门手机号"""
        LOG.info("------使用验证码注册：start!---------")
        verification_number_usa=Util().get_verification_number_signup_usa()
        result = self.signup_mvc('7022682192', '+853', verification_number_usa,"a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signup_mvc_phone_Taiwan(self):
        u"""台湾手机号"""
        LOG.info("------使用验证码注册：start!---------")
        verification_number_usa=Util().get_verification_number_signup_usa()
        result = self.signup_mvc('7022682192', '+886', verification_number_usa,"a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")


    def test_signup_mvc_phone_USA(self):
        u"""美国手机号"""
        LOG.info("------使用验证码注册：start!---------")
        verification_number_usa=Util().get_verification_number_signup_usa()
        result = self.signup_mvc('7022682192', '+1', verification_number_usa,"a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")


    def test_signup_mvc_phone_Jianada(self):
        u"""加拿大手机号"""
        LOG.info("------使用验证码注册：start!---------")
        verification_number_usa=Util().get_verification_number_signup_usa()
        result = self.signup_mvc('7022682192', '+1', verification_number_usa,"a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")


    def test_signup_mvc_phone_England(self):
        u"""英国手机号"""
        LOG.info("------使用验证码注册：start!---------")
        verification_number_usa=Util().get_verification_number_signup_usa()
        result = self.signup_mvc('7022682192', '+44', verification_number_usa,"a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")


    def test_signup_mvc_phone_Japan(self):
        u"""日本手机号"""
        LOG.info("------使用验证码注册：start!---------")
        verification_number_usa=Util().get_verification_number_signup_usa()
        result = self.signup_mvc('7022682192', '+81', verification_number_usa,"a1234567",data.nickname,"2019-04-11T09:08:29.568Z",0,data.height,data.weight)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")
if __name__ == '__main__':
    unittest.TestCase()
