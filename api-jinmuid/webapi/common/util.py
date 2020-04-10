import json
from urllib.parse import urljoin

import requests

from common.data import data
from common.log import LOG
from config.read_config import ReadConfig


class Util:
    def __init__(self):
        self.headers = {"Content-Type": "Application/json"}
        self.client_id = data.client_id
        self.secret_key_hash = data.secret_key_hash
        self.seed = ''
        self.sign_in_key = data.username
        self.register_type = 'username'
        self.password_hash = data.password
        self.sign_in_method = 'phone_password'
        self.username = ''
        self.phone = data.phone
        self.mvc = ''
        self.hashed_password = data.hashed_password
        self.seed_new = data.seed_phone
        self.sms_signin = 'sign_in'
        self.language = data.language
        self.nation_code = data.nation_code
        self.sms_signup = "sign_up"
        self.sms_reset_password = "reset_password"
        self.sms_set_phone_number = "set_phone_number"
        self.sms_modify_phone_number = "modify_phone_number"
        self.email_set_secure_email = "set_secure_email"
        self.email_unset_secure_email = "unset_secure_email"
        self.email_modify_secure_email = "modify_secure_email"
        self.email_find_username = "find_username"
        self.email_reset_password = "reset_password"
        self.email = data.email
        self.phone_usa="7022682192"
        self.nation_code_usa="+1"
        self.password=data.password
        self.seed_phone=data.seed_phone
        self.email=data.email

    host = ReadConfig().get_http('url')

    def get_authorization(self):
        data = {
            "client_id": self.client_id,
            "secret_key_hash": self.secret_key_hash,
            "seed": self.seed
        }
        url = urljoin(self.host, "client/auth")
        response = requests.post(url=url, json=data, headers=self.headers)
        return {"Content-Type": "application/json", "Authorization": response.json()['data']['authorization']}

    def get_token(self):
        data = {
            "client_id": self.client_id,
            "secret_key_hash": self.secret_key_hash,
            "seed": self.seed
        }
        url = urljoin(self.host, "client/auth")
        serial_number = Util().get_signin_serial_number()
        response = requests.post(url=url, json=data, headers=self.headers)
        data =  {"sign_in_method": 'phone_password',
                 "phone": self.phone,
                 "hashed_password": self.password,
                 "seed": self.seed_phone,
                 "nation_code": self.nation_code
                 }
        url = urljoin(self.host, 'signin')
        result = requests.post(url=url, json=data, headers=Util().get_authorization())
        print(data)

        print(result.json())
        response = {"Content-Type": "application/json", "authorization": response.json()['data']['authorization'],
                    "X-Access-Token": result.json()['data']['access_token']}
        print(response)
        return response

    def get_token_user(self):
        data = {
            "client_id": self.client_id,
            "secret_key_hash": self.secret_key_hash,
            "seed": self.seed
        }
        url = urljoin(self.host, "client/auth")
        serial_number = Util().get_signin_serial_number()
        response = requests.post(url=url, json=data, headers=self.headers)
        data ={"sign_in_method": 'username_password',
                "username": '31',
                "hashed_password": '9626c7444717aab7a3bbdd509bcafa35a7491e9478d421b38e539a621f695edd',
                "seed": ''
                }
        url = urljoin(self.host, 'signin')
        result = requests.post(url=url, json=data, headers=Util().get_authorization())
        print(data)

        print(result.json())
        response = {"Content-Type": "application/json", "authorization": response.json()['data']['authorization'],
                    "X-Access-Token": result.json()['data']['access_token']}
        print(response)
        return response

    def get_token_set_password(self):
        data = {
            "client_id": self.client_id,
            "secret_key_hash": self.secret_key_hash,
            "seed": self.seed
        }
        url = urljoin(self.host, "client/auth")
        serial_number = Util().get_signin_serial_number()
        response = requests.post(url=url, json=data, headers=self.headers)
        data = {"sign_in_method":  'phone_password',
                "username":'',
                "phone": self.phone,
                "mvc": "",
                "hashed_password": self.password,
                "seed": self.seed_new,
                "serial_number": serial_number,
                "nation_code": self.nation_code
                }
        url = urljoin(self.host, 'signin')
        result = requests.post(url=url, json=data, headers=Util().get_authorization())
        print(data)

        print(result.json())
        response = {"Content-Type": "application/json", "authorization": response.json()['data']['authorization'],
                    "X-Access-Token": result.json()['data']['access_token']}
        print(response)
        return response

    def get_user_id_set_password(self):
        serial_number = Util().get_signin_serial_number()
        data ={"sign_in_method": self.sign_in_method,
                "username": self.username,
                "phone": self.phone,
                "mvc": Util().phone_verification_code(),
                "hashed_password": self.password,
                "seed": self.seed_new,
                "serial_number": serial_number,
                "nation_code": self.nation_code
                }
        url = urljoin(self.host, 'signin')
        LOG.info("请求url:%s" % url)
        result = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % result.json())
        response = result.json()['data']['user_id']
        LOG.info(result.json())
        return response

    def get_user_id(self):
        #设置密码
        data =  {"sign_in_method": 'phone_password',
                 "phone": self.phone,
                 "hashed_password": self.password,
                 "seed": self.seed_phone,
                 "nation_code": self.nation_code
                 }
        url = urljoin(self.host, 'signin')
        LOG.info("请求url:%s" % url)
        result = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % result.json())
        response = result.json()['data']['user_id']
        LOG.info(result.json())
        return response

    def get_signin_serial_number(self):
        data = {"sms_Notification_type": self.sms_signin,
                "phone": self.phone,
                "language": self.language,
                "nation_code": self.nation_code  # 国家区号
                }
        url = self.host + "notification/sms"
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['serial_number']

    def get_signup_serial_number(self):
        data = {"sms_Notification_type": self.sms_signup,
                "phone": self.phone,
                "language": self.language,
                "nation_code": self.nation_code  # 国家区号
                }
        url = self.host + "notification/sms"
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['serial_number']

    def get_reset_password_serial_number(self):
        data = {"sms_Notification_type": self.sms_reset_password,
                "phone": self.phone,
                "language": self.language,
                "nation_code": self.nation_code  # 国家区号
                }
        url = self.host + "notification/sms"
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['serial_number']

    def get_set_phone_number_serial_number(self):
        data = {"sms_Notification_type": self.sms_set_phone_number,
                "phone": self.phone,
                "language": self.language,
                "nation_code": self.nation_code  # 国家区号
                }
        url = self.host + "notification/sms"
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['serial_number']

    def get_modify_phone_number_serial_number(self):
        data = {"sms_Notification_type": self.sms_modify_phone_number,
                "phone": self.phone,
                "language": self.language,
                "nation_code": self.nation_code  # 国家区号
                }
        url = self.host + "notification/sms"
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['serial_number']

    def phone_verification_code(self):
        data = {
            "send_information": [
                {
                    "send_via": "phone",
                    "email": "",
                    "phone": self.phone,
                    "nation_code": self.nation_code
                }]
        }
        url = self.host + "_debug/user/latest_verification_code"
        LOG.info("请求url:%s" % url)
        result = requests.post(url=url, data=json.dumps(data), headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % result.json())
        response = result.json()['data'][0]['verification_code']
        return response

    def email_verification_code(self):
        data = {
            "send_information": [
                {
                    "send_via": "email",
                    "email": self.email,
                    "phone": "",
                    "nation_code": ""
                }]
        }
        url = self.host + "_debug/user/latest_verification_code"
        LOG.info("请求url:%s" % url)
        result = requests.post(url=url, data=json.dumps(data), headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % result.json())
        response = result.json()['data'][0]['verification_code']
        return response

    def email_verification_code_new(self):
        data = {
            "send_information": [
                {
                    "send_via": "email",
                    "email": "601224464@qq.com",
                    "phone": "",
                    "nation_code": ""
                }]
        }
        url = self.host + "_debug/user/latest_verification_code"
        LOG.info("请求url:%s" % url)
        result = requests.post(url=url, data=json.dumps(data), headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % result.json())
        response = result.json()['data'][0]['verification_code']
        return response

    def get_verification_number(self):
        serial_number = Util().get_reset_password_serial_number()
        mvc = Util().phone_verification_code()
        data = {
            "phone": self.phone,
            "nation_code": self.nation_code,
            "mvc": mvc,
            "serial_number": serial_number
        }
        url = self.host + 'validate_signin_phone'
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['verification_number']

    def get_verification_number_email_modify_email(self):
        serial_number =Util().get_email_modify_secure_email_serial_number_old()
        mvc = Util().email_verification_code()
        data = {
            "email":'cuimin@jinmuhealth.com',
            "verification_code": mvc,
            "serial_number": serial_number,
            "verification_type": "modify_secure_email"
        }
        url = self.host + 'user/validate_email_verification_code'
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['verification_number']

    def get_email_set_secure_email_serial_number(self):
        url = self.host + "notification/email"
        LOG.info("请求url:%s" % url)
        data = {"type": self.email_set_secure_email,
                "email": self.email,
                "language": self.language,
                "user_id": self.get_user_id(),
                "send_to_new_if_modify": False
                }
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['serial_number']

    def get_email_unset_secure_email_serial_number(self):
        url = self.host + "notification/email"
        LOG.info("请求url:%s" % url)
        data = {"type": self.email_unset_secure_email,
                "email": self.email,
                "language": self.language,
                "user_id": self.get_user_id(),
                "send_to_new_if_modify": False
                }
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['serial_number']

    def get_email_modify_secure_email_serial_number(self):
        url = self.host + "notification/email"
        LOG.info("请求url:%s" % url)
        data = {"type": self.email_modify_secure_email,
                "email": "601224464@qq.com",
                "language": self.language,
                "user_id": self.get_user_id(),
                "send_to_new_if_modify": True
                }
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['serial_number']

    def get_email_modify_secure_email_serial_number_old(self):
        url = self.host + "notification/email"
        LOG.info("请求url:%s" % url)
        data = {"type": self.email_modify_secure_email,
                "email": self.email,
                "language": self.language,
                "user_id": self.get_user_id(),
                "send_to_new_if_modify": False
                }
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['serial_number']

    def get_email_find_username_serial_number(self):
        url = self.host + "notification/email"
        LOG.info("请求url:%s" % url)
        data = {"type": self.email_find_username,
                "email": self.email,
                "language": self.language,
                "user_id": self.get_user_id(),
                "send_to_new_if_modify": False
                }
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['serial_number']

    def get_email_reset_password_serial_number(self):
        url = self.host + "notification/email"
        LOG.info("请求url:%s" % url)
        data = {"type": self.email_reset_password,
                "email": self.email,
                "language": self.language,
                "user_id": self.get_user_id(),
                "send_to_new_if_modify": False
                }
        res = requests.post(url=url, json=data, headers=Util().get_token())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['serial_number']

    def get_verification_number_signup(self):
        serial_number =Util().get_signup_serial_number()
        mvc = Util().phone_verification_code()
        data = {
            "phone": self.phone,
            "nation_code": self.nation_code,
            "mvc": mvc,
            "serial_number":serial_number
        }
        url = self.host + 'user/validate_phone_verification_code'
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['verification_number']

    def get_signup_serial_number_usa(self):
        data = {"sms_Notification_type": self.sms_signup,
                "phone": self.phone_usa,
                "language": self.language,
                "nation_code": self.nation_code_usa  # 国家区号
                }
        url = self.host + "notification/sms"
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['serial_number']

    def get_verification_number_signup_usa(self):
        serial_number_usa =Util().get_signup_serial_number_usa()
        mvc = Util().phone_verification_code()
        data = {
            "phone": self.phone_usa,
            "nation_code": self.nation_code_usa,
            "mvc": mvc,
            "serial_number":serial_number_usa
        }
        url = self.host + 'user/validate_phone_verification_code'
        LOG.info("请求url:%s" % url)
        res = requests.post(url=url, json=data, headers=Util().get_authorization())
        LOG.info("请求参数:%s" % data)
        LOG.info("响应结果：%s" % res.json())
        return res.json()['data']['verification_number']
