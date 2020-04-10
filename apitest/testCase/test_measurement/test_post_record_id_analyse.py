import requests

from common.check_result import ApiTestCase
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class Analyse(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def analyse(self, str_record_id, language, question_answer, physical_dialectics, disease, dirty_dialectic):
        url = self.host + 'owner/measurements/'+str(Util().get_record_id())+'/v2/analyze'
        data = {
            "transaction_id": str_record_id,
            "language": language,
            "question_answers": question_answer,
            "physical_dialectics": physical_dialectics,
            "disease": disease,
            "dirty_dialectic": dirty_dialectic
        }
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("url是%s" % url)
        LOG.info("data是%s" % data)
        LOG.info("结果是%s" % res.json())
        return res.json()

    def test_get_analyse(self):
        question_answers= {
            "switch_pretreatment": [
                {
                    "question_key": "Q0009",
                    "answer_keys": ["QC0038"]
                },
                {
                    "question_key": "Q0009",
                    "answer_keys": ["QC0038"]
                }
            ],
            "breast_health": [
                {
                    "question_key": "Q0008",
                    "answer_keys": ["QC0033"]
                }
            ]
        },
        physical_dialectics= [
            {
                "key": "T0001"
            },
            {
                "key": "T0002"
            }
        ],
        disease= [
            {
                "key": "JB0001",
                "score": 1
            }
        ],
        dirty_dialectic= [
            {
                "key": "Z0001"
            }
        ]
        result = self.analyse(str(Util().get_record_id()), "zh-Hans", question_answers, physical_dialectics, disease, dirty_dialectic)
        self.assertOkResult(result)

    def test_get_analyse_error_language(self):
        question_answers= {
            "switch_pretreatment": [
                {
                    "question_key": "Q0009",
                    "answer_keys": ["QC0038"]
                },
                {
                    "question_key": "Q0009",
                    "answer_keys": ["QC0038"]
                }
            ],
            "breast_health": [
                {
                    "question_key": "Q0008",
                    "answer_keys": ["QC0033"]
                }
            ]
        },
        physical_dialectics= [
            {
                "key": "T0001"
            },
            {
                "key": "T0002"
            }
        ],
        disease= [
            {
                "key": "JB0001",
                "score": 1
            }
        ],
        dirty_dialectic= [
            {
                "key": "Z0001"
            }
        ]
        result = self.analyse(str(Util().get_record_id()), "chaa", question_answers, physical_dialectics, disease, dirty_dialectic)
        self.assertErrorResult(result, const.ErrInvalidValue)
