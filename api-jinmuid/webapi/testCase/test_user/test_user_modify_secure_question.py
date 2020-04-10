import json
import unittest

import requests

from common.check_result import ApiTestCase
from common.data import data
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class UserModifySecureQuestion(ApiTestCase):
    u'''设置/修改密保'''

    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def validate_question_before_modify(self, old_question_key1, old_answer1, old_question_key2, old_answer2,
                                        old_question_key3, old_answer3, new_question_key1, new_answer1,
                                        new_question_key2, new_answer2, new_question_key3, new_answer3):
        data = {
            "old_secure_questions": [
                {
                    "question_key": old_question_key1,
                    "answer": old_answer1
                },
                {
                    "question_key": old_question_key2,
                    "answer": old_answer2
                },
                {
                    "question_key": old_question_key3,
                    "answer": old_answer3
                }
            ],
            "new_secure_questions": [
                {
                    "question_key": new_question_key1,
                    "answer": new_answer1
                }, {
                    "question_key": new_question_key2,
                    "answer": new_answer2
                },
                {
                    "question_key": new_question_key3,
                    "answer": new_answer3
                }
            ]
        }
        url = self.host + "user/" + str(Util().get_user_id()) + "/modify_secure_question"
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        req = requests.post(url=url, data=json.dumps(data), headers=Util().get_token())
        return req.json()

    def test_user_questions_correct_parameters(self):
        u"""正确参数"""
        LOG.info("------设置/修改密保：start!---------")
        result = self.validate_question_before_modify(data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3, data.new_question_key1, data.new_answer1, data.new_question_key2,
                                                      data.new_answer2, data.new_question_key3, data.new_answer3)
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_user_questions_error_key(self):
        u"""错误参数"""
        LOG.info('------设置/修改密保：start!---------')
        result = self.validate_question_before_modify('', 'a', '', 'b', '', '西红柿炒蕃茄', '', '林俊杰', '', '包租婆', '',
                                                      '请把我留在最好的时光里')
        LOG.info('获取测试结果：%s' % result)
        self.assertErrorResult(result, const.ErrWrongFormatQuestion)
        LOG.info('------pass!---------')

    def test_user_questions_error_answer(self):
        u"""错误参数"""
        LOG.info('------设置/修改密保：start!---------')
        result = self.validate_question_before_modify('1', '', '2', '', '3', '', '18', '林俊杰', '13', '包租婆', '14',
                                                      '请把我留在最好的时光里')
        LOG.info('获取测试结果：%s' % result)
        self.assertErrorResult(result, const.ErrEmptyAnswer)
        LOG.info('------pass!---------')

    def test_user_questions_answer_long(self):
        u"""答案超出长度"""
        LOG.info("------设置/修改密保：start!---------")
        result = self.validate_question_before_modify(data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3, data.new_question_key1, data.new_answer1, data.new_question_key2,
                                                      data.new_answer2, data.new_question_key3, '马克思列宁主义毛泽东思想邓小平理论三个代表重要思想科学发展观导读')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_user_questions_same(self):
        u"""新密保和旧密保一致"""
        LOG.info("------设置/修改密保：start!---------")
        result = self.validate_question_before_modify(data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3,data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3)
        LOG.info("获取测试结果：%s" % result)
        self.assertErrorResult(result, const.ErrSameSecureQuestion)
        LOG.info("------pass!---------")

    def test_user_questions_answer_format(self):
        u"""答案为特殊字符"""
        LOG.info("------设置/修改密保：start!---------")
        result = self.validate_question_before_modify(data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3, '18', '林俊杰', '13', '包租婆', '14',
                                                      '/ ')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result)
        LOG.info("------pass!---------")

    def test_user_questions_question_not_exit(self):
        u"""问题不存在"""
        LOG.info("------设置/修改密保：start!---------")
        result = self.validate_question_before_modify(data.question_key1, data.answer1, data.question_key2, data.answer2, data.question_key3,
                                                   data.answer3, '11', 'a', '20', 'b', '254',
                                                      '西红柿炒蕃茄')
        LOG.info("获取测试结果：%s" % result)
        self.assertOkResult(result, const.ErrWrongFormatQuestion)
        LOG.info("------pass!---------")

    def test_user_questions_error_paramters(self):
        u"""只有两上密保问题"""
        LOG.info('------设置/修改密保：start!---------')
        data = {
            "old_secure_questions": [
                {
                    "question_key": 1,
                    "answer": "a"
                },
                {
                    "question_key": 3,
                    "answer": "a"
                }
            ],
            "new_secure_questions": [
                {
                    "question_key": 1,
                    "answer": "b"
                }, {
                    "question_key": 2,
                    "answer": "a"
                }
            ]
        }
        url = self.host + "user/" + str(Util().get_user_id()) + "/modify_secure_question"
        LOG.info("请求url:%s" % url)
        LOG.info("请求参数:%s" % data)
        result = requests.post(url=url, data=json.dumps(data), headers=Util().get_token())
        LOG.info("请求参数:%s" % result.json())
        self.assertErrorResult(result.json(), const.ErrParsingRequestFailed)
        LOG.info('------pass!---------')

    if __name__ == '__main__':
        unittest.TestCase()
