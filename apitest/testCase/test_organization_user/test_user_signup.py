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

    def signup(self, client_id, register_type, nickname, gender, birthday, height, weight):
        data = {
            "client_id": client_id,
            "register_type": register_type,
            "nickname": nickname,
            "gender": gender,
            "birthday": birthday,
            "height": height,
            "weight": weight
        }
        url = self.host + 'owner/' + str(Util().get_user_id()) + '/users/sign_up'
        res = requests.post(url=url, json=data, headers=Util().get_token())
        return res.json()

    def test_signin_correct_parameters(self):
        result = self.signup("jm-10001", "phone", "smiling", 0, "2019-06-11T06:55:35.978Z", 160, 45)
        self.assertOkResult(result)

    def test_signin_error_nickname_special_signal(self):
        result = self.signup("jm-10001", "phone", "&smiling", 0, "2019-06-11T06:55:35.978Z", 160, 45)
        self.assertErrorResult(result,const.ErrParsingRequestFailed)

    def test_signin_error_nickname_length(self):
        result = self.signup("jm-10001", "phone", "smilingaaaaaaaaaaaaaaaaaaaa", 0, "2019-06-11T06:55:35.978Z", 160, 45)
        self.assertErrorResult(result,const.ErrParsingRequestFailed)

    def test_signin_error_birthday(self):
        result = self.signup("jm-10001", "phone", "smiling", 0, "1800-06-11T06:55:35.978Z", 160, 45)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)

    def test_signin_error_height(self):
        result = self.signup("jm-10001", "phone", "smiling", 0, "2019-06-11T06:55:35.978Z", 30, 45)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)

    def test_signin_error_weight(self):
        result = self.signup("jm-10001", "phone", "smiling", 0, "2019-06-11T06:55:35.978Z", 160, 10)
        self.assertErrorResult(result, const.ErrParsingRequestFailed)