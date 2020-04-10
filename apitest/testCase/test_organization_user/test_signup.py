from urllib.parse import urljoin

import requests

from common.check_result import ApiTestCase
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class SignUp(ApiTestCase):

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def signup(self, client_id, password, register_type, nickname, birthday, gender, height, weight, phone, email,
               remark, user_defined_code, state, city, street, country):
        data = {
            "client_id": client_id,
            "password": password,
            "register_type": register_type,
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
        # url = self.host + "owner/users/signup"
        url = urljoin(self.host, 'owner/users/signup')
        LOG.info("请求url:%s" % url)
        print(url)
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_signup_correct_required_parameters(self):
        u"""只填必填项"""
        LOG.info("------登录成功用例：start!---------")
        url = urljoin(self.host, 'owner/users/signup')
        data = {
            "client_id": "jm-10001",
            "password": "28acc304c6a7252e0796e0f49f60f31bc436de66fed042e97b11619e1ee1a986",
            "register_type": "username",
            "profile":
                {
                    "nickname": "smile",
                    "birthday": "2018-10-15T02:41:31Z",
                    "gender": 1,
                    "height": 160,
                    "weight": 45,
                }
        }
        result = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("获取测试结果：%s" % result.json())
        self.assertOkResult(result.json())
        LOG.info("------pass!---------")

    def test_signup_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.signup('jm-10001', '28acc304c6a7252e0796e0f49f60f31bc436de66fed042e97b11619e1ee1a986',
                             'username', 'smile', '2018-10-15T02:41:31Z', 1, 160, 45, '13700009001', '22389@qq.com',
                             'aaa', '1', '江苏省', '常州市', '天宁区', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_signup_error_parameters(self):
        u"""‘’错误参数格式''"""
        LOG.info("------登录成功用例：start!---------")
        client_id = 'jm-10001',
        password = '28acc304c6a7252e0796e0f49f60f31bc436de66fed042e97b11619e1ee1a986',
        register_type = 'username',
        nickname = 'smile',
        birthday = '2018-10-15T02:41:31Z',
        gender = '1',
        height = '160'
        weight = '45',
        phone = '13700009001',
        email = '22389@qq.com'
        remark = 'aaa',
        user_defined_code = '1',
        state = '江苏省',
        city = '常州市',
        street = '天宁区',
        country = '中国'
        result = self.signup(client_id, password, register_type, nickname, birthday, gender, height, weight, phone,
                             email, remark, user_defined_code, state, city, street, country)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_signup_error_client_id(self):
        u"""client_id错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.signup('', '28acc304c6a7252e0796e0f49f60f31bc436de66fed042e97b11619e1ee1a986',
                             'username', 'smile', '2018-10-15T02:41:31Z', 1, 160, 45, '', '',
                             'aaa', '1', '江苏省', '常州市', '天宁区', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrIncorrectClientID)
        LOG.info("------pass!---------")

    def test_signup_error_password(self):
        u"""password错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.signup('jm-10001', '',
                             'username', 'smile', '2018-10-15T02:41:31Z', 1, 160, 45, '', '',
                             'aaa', '1', '江苏省', '常州市', '天宁区', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidPassword)
        LOG.info("------pass!---------")

    def test_signup_error_register_type(self):
        u"""register_type错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.signup('jm-10001', '28acc304c6a7252e0796e0f49f60f31bc436de66fed042e97b11619e1ee1a986',
                             '', 'smile', '2018-10-15T02:41:31Z', 1, 160, 45, '', '',
                             'aaa', '1', '江苏省', '常州市', '天宁区', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info("------pass!---------")

    def test_signup_error_nickname(self):
        u"""nickname错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.signup('jm-10001', '28acc304c6a7252e0796e0f49f60f31bc436de66fed042e97b11619e1ee1a986',
                             'username', '', '2018-10-15T02:41:31Z', 1, 160, 45, '', '',
                             'aaa', '1', '江苏省', '常州市', '天宁区', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info("------pass!---------")

    def test_signup_error_birthday(self):
        u"""birthday错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.signup('jm-10001', '28acc304c6a7252e0796e0f49f60f31bc436de66fed042e97b11619e1ee1a986',
                             'username', 'smile', '', 1, 160, 45, '', '',
                             'aaa', '1', '江苏省', '常州市', '天宁区', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_signup_error_gender(self):
        u"""gender错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.signup('jm-10001', '28acc304c6a7252e0796e0f49f60f31bc436de66fed042e97b11619e1ee1a986',
                             'username', 'smile', '2018-10-15T02:41:31Z', '', 160, 45, '', '',
                             'aaa', '1', '江苏省', '常州市', '天宁区', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_signup_error_height(self):
        u"""height错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.signup('jm-10001', '28acc304c6a7252e0796e0f49f60f31bc436de66fed042e97b11619e1ee1a986',
                             'username', 'smile', '2018-10-15T02:41:31Z', 1, '', 45, '', '',
                             'aaa', '1', '江苏省', '常州市', '天宁区', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")

    def test_signup_error_weight(self):
        u"""weight错误参数"""
        LOG.info("------登录成功用例：start!---------")
        result = self.signup('jm-10001', '28acc304c6a7252e0796e0f49f60f31bc436de66fed042e97b11619e1ee1a986',
                             'username', 'smile', '2018-10-15T02:41:31Z', 1, 160, '', '', '',
                             'aaa', '1', '江苏省', '常州市', '天宁区', '中国')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)
        LOG.info("------pass!---------")
