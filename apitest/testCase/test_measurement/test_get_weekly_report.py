import requests

from common.check_result import ApiTestCase
from common.errCode import const
from common.log import LOG
from common.util import Util
from config.read_config import ReadConfig


class WeeklyReport(ApiTestCase):
    def setUp(self):
        LOG.info('测试用例开始执行')

    def tearDown(self):
        LOG.info('测试用例执行完毕')

    host = ReadConfig().get_http('url')

    def weeklyreport(self, c0, c1, c2, c3, c4, c5, c6, c7, language, user_id):
        url = self.host + 'owner/measurements/v2/weekly_report'
        data = {
            "c0": c0,
            "c1": c1,
            "c2": c2,
            "c3": c3,
            "c4": c4,
            "c5": c5,
            "c6": c6,
            "c7": c7,
            "language": language,
            "user_id": user_id
        }
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("url是%s" % url)
        LOG.info("data是%s" % data)
        LOG.info("结果是%s" % res.json())
        return res.json()

    def test_get_weekly_report_data(self):
        result = self.weeklyreport(0, 0, 0, 0, 0, 0, 0, 0, "zh-Hans", Util().get_user_id())
        self.assertOkResult(result)

    def test_get_weekly_report_data_error_c0_less(self):
        result = self.weeklyreport(-20, 0, 0, 0, 0, 0, 0, 0, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_c0_more(self):
        result = self.weeklyreport(20, 0, 0, 0, 0, 0, 0, 0, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_c1_less(self):
        result = self.weeklyreport(0, -20, 0, 0, 0, 0, 0, 0, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_c1_more(self):
        result = self.weeklyreport(0, 20, 0, 0, 0, 0, 0, 0, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_c2_less(self):
        result = self.weeklyreport(0, 0, -20, 0, 0, 0, 0, 0, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_c2_more(self):
        result = self.weeklyreport(0, 0, 20, 0, 0, 0, 0, 0, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_c3_less(self):
        result = self.weeklyreport(0, 0, 0, -20, 0, 0, 0, 0, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_c3_more(self):
        result = self.weeklyreport(0, 0, 0, 20, 0, 0, 0, 0, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_c4_less(self):
        result = self.weeklyreport(0, 0, 0, 0, -20, 0, 0, 0, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_c4_more(self):
        result = self.weeklyreport(0, 0, 0, 0, 20, 0, 0, 0, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_c5_less(self):
        result = self.weeklyreport(0, 0, 0, 0, 0, -20, 0, 0, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_c5_more(self):
        result = self.weeklyreport(0, 0, 0, 0, 0, 20, 0, 0, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_c6_less(self):
        result = self.weeklyreport(0, 0, 0, 0, 0, 0, -20, 0, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_c6_more(self):
        result = self.weeklyreport(0, 0, 0, 0, 0, 0, 20, 0, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_c7_less(self):
        result = self.weeklyreport(0, 0, 0, 0, 0, 0, 0, -20, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_c7_more(self):
        result = self.weeklyreport(0, 0, 0, 0, 0, 0, 0, 20, "zh-Hans", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_language(self):
        result = self.weeklyreport(0, 0, 0, 0, 0, 0, 0, 0, "ch", Util().get_user_id())
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_user_id_type(self):
        result = self.weeklyreport(0, 0, 0, 0, 0, 0, 0, 0, "ch", "496846465156156151151")
        self.assertErrorResult(result, const.ErrRPCInternal)

    def test_get_weekly_report_data_error_user_id_value(self):
        result = self.weeklyreport(0, 0, 0, 0, 0, 0, 0, 0, "ch", 496846465156156151151)
        self.assertErrorResult(result, const.ErrRPCInternal)
