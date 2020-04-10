import re
import unittest

from common.db import Database


# noinspection PyCallByClass
class ApiTestCase(unittest.TestCase):
    organization_id = 786
    name = 31
    user_id = 786
    nickname = 'smile'

    def assertOkResult(self, result):
        self.assertEqual(result['ok'], True)

    def assertRpcInternalErrorCode(self, result, expectedCode):
        self.assertEqual(int(re.search(r'\[errcode:(\d+)\]', result['error']['msg']).group(1)), expectedCode)

    def assertErrorResult(self, result, expectedCode):
        self.assertEqual(result['ok'], False)
        self.assertEqual(result['error']['code'], expectedCode)

    def assertOkResult_get_organization(self, result):
        self.assertEqual(result['ok'], True)

        ex = Database().execute("select * from organization where organization_id = '%s' " % (self.organization_id))
        name = Database().fetch_one(ex)[2]
        self.assertEqual(result['data'][0]['profile']['name'], name)

    # noinspection PyCallByClass
    def assertOkResult_organization(self, result):
        self.assertEqual(result['ok'], True)
        ex = Database.execute("select * from user where name = '%s' " % (self.name))
        print(Database.fetch_one(ex)[2])
        self.assertEqual(result['data']['profile']['name'], Database.fetch_one(ex[2]))

    def assertOkResult_signup(self, result):
        self.assertEqual(result['ok'], True)
        ex = Database().execute("select * from user where nickname = '%s' " % (self.nickname))
        nickname = Database().fetch_one(ex)[9]
        self.assertEqual(result['data']['profile']['nickname'], nickname)

    def assertOkResult_profile(self, result):
        self.assertEqual(result['ok'], True)
        ex = Database().execute("select * from user_profile where user_id = '%s' " % (self.user_id))
        username = Database().fetch_one(ex)[2]
        self.assertEqual(result['data']['username'], username)
        self.assertEqual(result['data']['profile']['nickname'], 'smile')

    def assertOkResult_preferences(self, result):
        self.assertEqual(result['ok'], True)
        ex = Database().execute("select * from user_preferences where user_id = '%s' " % (self.user_id))
        enable_heart_rate_chart = Database().fetch_one(ex)[1]
        self.assertEqual(result['data']['enable_choose_status'], enable_heart_rate_chart)
        self.assertEqual(result['data'], {'enable_choose_status': True, 'enable_comment': True,
                                          'enable_constitution_differentiation': True,
                                          'enable_health_trending': True, 'enable_heart_rate_chart': True,
                                          'enable_meridian_bar_graph': True, 'enable_pulse_wave_chart': True,
                                          'enable_syndrome_differentiation': True, 'enable_warm_prompt': True,
                                          'enable_western_medicine_analysis': True})

    def delete_organization(self):
        sql = "UPDATE organization SET deleted_at =  cast(getdate() as datetime) where organization_id=786; "
        ex = Database().execute(sql)
