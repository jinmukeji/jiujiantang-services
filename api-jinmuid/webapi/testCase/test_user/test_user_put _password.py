from urllib.parse import urljoin

import requests

from common.check_result import ApiTestCase
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserPasswordPut(ApiTestCase):
    u'''用户修改密码'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def user_password_put(self, old_hashed_password, new_hashed_password, seed):
        data = {"old_hashed_password": old_hashed_password,
                "new_plain_password": new_hashed_password,
                "seed": seed
                }
        url = self.host + 'user/' + str(Util().get_user_id()) + '/password'
        LOG.info("请求url:%s" % url)
        res = requests.put(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        return res.json()

    def test_user_password_put_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------用户修改密码：start!---------")
        result = self.user_password_put('2a2ae8e0d4be04d08802f4f1dbeb2bc151eda934555af8004440a7a950ce8dda', '654321',
                                        'a123')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_user_password_put_error_old_hashed_password(self):
        u"""old_hashed_password错误参数"""
        LOG.info("------用户修改密码：start!---------")
        result = self.user_password_put('', '654321', 'a123')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info("------pass!---------")

    def test_user_password_error_new_hashed_password(self):
        u"""new_hashed_password错误参数"""
        LOG.info("------用户修改密码：start!---------")
        result = self.user_password_put('2a2ae8e0d4be04d08802f4f1dbeb2bc151eda934555af8004440a7a950ce8dda', '', 'a123')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info("------pass!---------")

    def test_user_password_error_seed(self):
        u"""seed错误参数"""
        LOG.info("------用户修改密码：start!---------")
        result = self.user_password_put('2a2ae8e0d4be04d08802f4f1dbeb2bc151eda934555af8004440a7a950ce8dda', '654321',
                                        '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info("------pass!---------")

    def test_user_password_error_password(self):
        u"""新旧密码一致"""
        LOG.info("------用户修改密码：start!---------")
        result = self.user_password_put('2a2ae8e0d4be04d08802f4f1dbeb2bc151eda934555af8004440a7a950ce8dda', '123456',
                                        'a123')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrInvalidValue)
        LOG.info("------pass!---------")

    def test_user_password_old(self):
        u"""恢复成原始密码"""
        data = {"sign_in_method": 'username_password',
                "username": 'hello',
                "hashed_password": '98c9402bc506326eea9c507921fe16e06b2175217ae498eeb313de2efa8b55ff',
                "seed": 'a123'
                }
        url = urljoin(self.host, 'signin')
        result = requests.post(url=url, json=data, headers=Util().get_authorization())
        req = result.json()['data']['user_id']
        data_pwd = {"old_hashed_password": '98c9402bc506326eea9c507921fe16e06b2175217ae498eeb313de2efa8b55ff',
                    "new_plain_password": '123456',
                    "seed": 'a123'
                    }
        url_pwd = self.host + 'user/' + str(req) + '/password'
        LOG.info("请求url:%s" % url)
        res = requests.put(url=url_pwd, json=data_pwd, headers=Util().get_token())
        LOG.info("获取测试结果：%s" % res.json())
