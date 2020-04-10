import re
import unittest


# noinspection PyCallByClass
class ApiTestCase(unittest.TestCase):

    def assertOkResult(self, result):
        self.assertEqual(result['ok'], True)

    def assertRpcInternalErrorCode(self, result, expectedCode):
        self.assertEqual(int(re.search(r'\[errcode:(\d+)\]', result['error']['msg']).group(1)), expectedCode)

    def assertErrorResult(self, result, expectedCode):
        self.assertEqual(result['ok'], False)
        self.assertEqual(result['error']['code'], expectedCode)

    def assertErrorQuestion(self, result, expectedCode):
        self.assertEqual(result['ok'], True)
        if result['data']['code'] == result:
            self.assertIn(result['data']['wrong_question_keys'], 2)
