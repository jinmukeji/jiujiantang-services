import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserModifyPasswordViaQuestion(ApiTestCase):
    u'''根据密保问题修改密码'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def modify_password_via_question(self, validation_type, username, phone, nation_code, password, question_key1,
                                     answer1,
                                     question_key2, answer2, question_key3, answer3):
        data = {
            "validation_type": validation_type,
            "username": username,
            "phone": phone,
            "nation_code": nation_code,
            "password": password,
            "secure_questions": [
                {
                    "question_key": question_key1,
                    "answer": answer1
                },
                {
                    "question_key": question_key2,
                    "answer": answer2
                },
                {
                    "question_key": question_key3,
                    "answer": answer3
                }
            ]
        }
        url = self.host + "user/modify_password_via_question"
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        req = requests.post(url=url, json=data, headers=Util().get_authorization())
        return req.json()

    def test_user_questions_correct_parameters_username(self):
        u"""用户名正确参数"""
        LOG.info("------根据密保问题修改密码：start!---------")
        result = self.modify_password_via_question('username', data.username, '', '', data.pwd, data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_user_questions_correct_parameters_phone(self):
        u"""手机号正确参数"""
        LOG.info("------根据密保问题修改密码：start!---------")
        result = self.modify_password_via_question('phone', '', data.phone, data.nation_code, data.pwd,  data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_user_questions_username_type_is_null(self):
        u"""type为空"""
        LOG.info("------根据密保问题修改密码：start!---------")
        result = self.modify_password_via_question('', data.username, '', '', data.pwd,  data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result,const.ErrInvalidSecureQuestionValidationMethod)
        LOG.info("------pass!---------")

    def test_user_questions_phone_type_is_null(self):
        u"""type为空"""
        LOG.info("------根据密保问题修改密码：start!---------")
        result = self.modify_password_via_question('', '',data.phone, data.nation_code, data.pwd, data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNonexistentUsername)
        LOG.info("------pass!---------")

    def test_user_questions_phone_format(self):
        u"""手机格式不正确"""
        LOG.info("------根据密保问题修改密码：start!---------")
        result = self.modify_password_via_question('phone', '', '13221058', data.nation_code, data.pwd, data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNonexistentPhone)
        LOG.info("------pass!---------")

    def test_user_questions_phone_not_nation_code(self):
        u'''手机格式和区号不一致'''
        LOG.info("------根据密保问题修改密码：start!---------")
        result = self.modify_password_via_question('phone', '', data.phone, '+1', data.pwd,  data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNonexistentPhone)
        LOG.info("------pass!---------")

    def test_user_questions_phone_not_exit(self):
        u'''手机号码不存在'''
        LOG.info("------根据密保问题修改密码：start!---------")
        result = self.modify_password_via_question('phone', '', data.phone_not_exist, data.nation_code, data.pwd, data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNonexistentPhone)
        LOG.info("------pass!---------")

    def test_user_questions_username_not_exit(self):
        u'''用户名不存在'''
        LOG.info("------根据密保问题修改密码：start!---------")
        result = self.modify_password_via_question('username', 'lili', '', '', data.pwd, data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNonexistentUsername)
        LOG.info("------pass!---------")

    def test_user_questions_question_error(self):
        u"""密保答案错误"""
        LOG.info("------根据密保问题修改密码：start!---------")
        result = self.modify_password_via_question('username', data.username, '', '', data.pwd,  data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   "123456")
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNonexistentUsername)
        LOG.info("------pass!---------")

    def test_user_questions_password_is_null(self):
        u"""密码为空"""
        LOG.info("------根据密保问题修改密码：start!---------")
        result = self.modify_password_via_question('username', data.username, '', '', '',  data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrNonexistentUsername)
        LOG.info("------pass!---------")

    def test_user_questions_password_error(self):
        u"""密码存在特殊字符"""
        LOG.info("------根据密保问题修改密码：start!---------")
        result = self.modify_password_via_question('username',  data.username, '', '', '/123456',  data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result,const.ErrNonexistentUsername)
        LOG.info("------pass!---------")
