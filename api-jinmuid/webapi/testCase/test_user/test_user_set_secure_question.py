import json
import unittest

import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserSetSecureQuestion(ApiTestCase):
    u'''设置/修改密保'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def set_secure_question(self, question_key1, answer1, question_key2, answer2, question_key3, answer3):
        data = {"secure_questions": [
            {
                "question_key": question_key1,
                "answer": answer1
            }, {
                "question_key": question_key2,
                "answer": answer2
            }, {
                "question_key": question_key3,
                "answer": answer3
            }]}
        url = self.host + "user/" + str(Util().get_user_id()) + "/set_secure_question"
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        req = requests.post(url=url, data=json.dumps(data), headers=Util().get_token())
        return req.json()

    def test_user_questions_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------设置/修改密保：start!---------")
        result = self.set_secure_question(data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3, data.answer3)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_user_questions_error_key(self):
        u"""错误参数"""
        LOG.info('------设置/修改密保：start!---------')
        result = self.set_secure_question('', 'a', '', 'b', '', '西红柿炒蕃茄')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrWrongFormatQuestion)
        LOG.info('------pass!---------')

    def test_user_questions_error_answer(self):
        u"""错误参数"""
        LOG.info('------设置/修改密保：start!---------')
        result = self.set_secure_question('1', '', '2', '', '3', '')
        LOG.info('获取测试结果：%s' % result)
        self.assertErrorResult(result, const.ErrEmptyAnswer)
        LOG.info('------pass!---------')

    def test_user_questions_answer_long(self):
        u"""答案超出长度"""
        LOG.info("------设置/修改密保：start!---------")
        result = self.set_secure_question('1', '西红柿炒蕃茄炒萝卜炒西蓝花炒肉炒青菠', '2', '西红柿炒蕃茄', '3', '西红柿炒蕃茄')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrSecureQuestionExist)
        LOG.info("------pass!---------")

    def test_user_questions_answer_same(self):
        u"""问题和答案一致"""
        LOG.info("------设置/修改密保：start!---------")
        result = self.set_secure_question('1', '西红柿炒蕃茄', '1', '西红柿炒蕃茄', '1', '西红柿炒蕃茄')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrRepeatedQuestion)
        LOG.info("------pass!---------")

    def test_user_questions_first(self):
        u"""只填写一个问题"""
        LOG.info("------设置/修改密保：start!---------")
        result = self.set_secure_question('1', '西红柿炒蕃茄', '', '', '', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrEmptyAnswer)
        LOG.info("------pass!---------")

    def test_user_questions_second(self):
        u"""只填写两个问题"""
        LOG.info("------设置/修改密保：start!---------")
        result = self.set_secure_question('1', '西红柿炒蕃茄', '2', 'a', '', '')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrEmptyAnswer)
        LOG.info("------pass!---------")

    def test_user_questions_answer_format(self):
        u"""答案为特殊字符"""
        LOG.info("------设置/修改密保：start!---------")
        result = self.set_secure_question('1', '西红柿炒蕃茄', '2', 'a', '3', '#a/')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrSecureQuestionExist)
        LOG.info("------pass!---------")

    def test_user_questions_question_not_exit(self):
        u"""问题不存在"""
        LOG.info("------设置/修改密保：start!---------")
        result = self.set_secure_question('17', '西红柿炒蕃茄', '20', 'a', '3', '1')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrWrongFormatQuestion)
        LOG.info("------pass!---------")

    def test_user_questions_error_paramters(self):
        u"""错误参数,只传两个"""
        LOG.info('------设置/修改密保：start!---------')
        data = {"secure_questions": [
            {
                "question_key": '1',
                "answer": '1'
            }, {
                "question_key": '2',
                "answer": '1'
            }]}
        url = self.host + "user/" + str(Util().get_user_id()) + "/set_secure_question"
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        result = requests.post(url=url, data=json.dumps(data), headers=Util().get_token())
        LOG.info(result.json())
        self.assertErrorResult(result.json(), const.ErrWrongSecureQuestionCount)
        LOG.info('------pass!---------')

    def test_user_question_ist_exist(self):
        u"""密保问题已设置 """
        LOG.info("------设置/修改密保：start!---------")
        result = self.set_secure_question('1', 'a', '2', 'b', '3', '西红柿炒蕃茄')
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrSecureQuestionExist)
        LOG.info("------pass!---------")


if __name__ == '__main__':
    unittest.TestCase()
