import requests

from common.check_result import ApiTestCase
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UsersPutProfile(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def put_profile(self, nickname, birthday, gender, height, weight, phone, email, remark, user_defined_code,
                    state, city, street, country):
        data = {
            "profile":
                {
                    "nickname": nickname,
                    "birthday": birthday,
                    "gender": gender,
                    "height": height,
                    "weight": weight,
                    "phone": phone,
                    "email": email,
                    "remark": remark,
                    "user_defined_code": user_defined_code,
                    "state": state,
                    "city": city,
                    "street": street,
                    "country": country
                }
        }
        url = self.host + 'owner/users/' + str(Util().get_user_id()) + '/profile'
        LOG.info("请求url:%s" % url)
        res = requests.put(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_put_profile_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.put_profile('smile', '2018-10-15T02:41:31Z', 0, 160, 45, '13700007003', '2334@qq.com', 'a',
                                  '1', '江苏省',
                                  '常州市', '天宁区', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_put_profile(self):
        data = {
            "profile":
                {
                    "nickname": "哈哈",
                    "birthday":  '2018-10-15T02:41:31Z',
                    "gender": 1,
                    "height": 160,
                    "weight": 45,
                    "phone": "",
                    "email": "",
                    "remark": "",
                    "user_defined_code": "",
                    "state": "",
                    "city": "",
                    "street": "",
                    "country": ""
                }
        }
        url = self.host + 'owner/users/118526/profile'
        LOG.info("请求url:%s" % url)
        res = requests.put(url=url, json=data, headers=Util().get_token())
        LOG.info("获取测试结果：%s" % res.json())
        self.assertOkResult(res.json())
        LOG.info("------pass!---------")



    def test_put_profile_error_nickname(self):
        u"""nickname错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.put_profile('', '2018-10-15T02:41:31Z', 0, 160, 45, '13700007003', '2334@qq.com', 'a', '1',
                                  '江苏省',
                                  '常州市', '天宁区', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info("------pass!---------")

    def test_put_profile_error_birthday(self):
        u"""birthday错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.put_profile('smile', '2018-10', 0, 160, 45, '13700007003', '2334@qq.com', 'a', '1', '江苏省',
                                  '常州市', '天宁区', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    # def test_put_profile_error_gender(self):
    #     u"""gender错误参数"""
    #     LOG.info("------登录成功用例：start!---------")
    #     result = self.put_profile('smile', '2018-10-15T02:41:31Z', 8, 160, 45, '13700007003', '2334@qq.com', 'a',
    #                               '1', '江苏省',
    #                               '常州市', '天宁区', '中国')
    #     LOG.info("获取测试结果：%s" % result)
    #     self.assertErrorResult(result, const.ErrRPCInternal)
    #     LOG.info("------pass!---------")

    def test_put_profile_error_height(self):
        u"""height错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.put_profile('smile', '2018-10-15T02:41:31Z', 0, '', 45, '13700007003', '2334@qq.com', 'a',
                                  '1', '江苏省',
                                  '常州市', '天宁区', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_put_profile_error_weight(self):
        u"""weight错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.put_profile('smile', '2018-10-15T02:41:31Z', 0, 160, '', '13700007003', '2334@qq.com', 'a',
                                  '1', '江苏省',
                                  '常州市', '天宁区', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")
