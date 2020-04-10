import os
import smtplib
import time
import unittest
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText

from common.html_runner import Htmlrunner
from common.log import LOG

CUR_DIR = os.path.dirname(os.path.realpath(__file__))


class AllTest:
    def __init__(self):

        self.caseListFile = os.path.join(CUR_DIR, "caselist")
        LOG.info(self.caseListFile)
        self.caseFile = os.path.join(CUR_DIR, "testCase")
        # self.caseFile = None
        self.caseList = []

    def load_case_list(self):
        """
        set case list
        :return:
        """
        fb = open(self.caseListFile)
        for value in fb.readlines():
            data = str(value)
            if data != '' and not data.startswith("#"):
                self.caseList.append(data.replace("\n", ""))
        fb.close()

    def load_testsuite(self):
        """
        set case suite
        :return:
        """
        self.load_case_list()
        test_suite = unittest.TestSuite()
        suite_module = []

        for case in self.caseList:
            case_name = case.split("/")[-1]
            print(case_name + ".py")
            discover = unittest.defaultTestLoader.discover(self.caseFile, pattern=case_name + '.py', top_level_dir=None)
            suite_module.append(discover)
        if len(suite_module) > 0:

            for suite in suite_module:
                for test_name in suite:
                    test_suite.addTest(test_name)
        else:
            return None

        return test_suite

    def run_testsuite(self, job_id):
        """
        run test
        :return:
        """
        try:
            suit = self.load_testsuite()
            if suit is not None:
                # LOG.info("********TEST START********")
                now = time.strftime("%Y_%m_%d_%H_%M_%S")
                result_path = os.path.join(CUR_DIR, "report")
                if not os.path.exists(result_path):
                    os.mkdir(result_path)

                if job_id == '':
                    report_abspath = os.path.join(result_path, now + "_result.html")
                else:
                    report_abspath = os.path.join(result_path, job_id + "_result.html")
                print("report path:%s" % report_abspath)
                report = open(report_abspath, "wb")
                runner = Htmlrunner(stream=report, title=u'自动化测试报告,测试结果如下：', description=u'用例执行情况：')

                runner.run(suit)
            else:
                LOG.info("Have no case to test.")
        except Exception as ex:
            LOG.error(str(ex))
        # finally:
        # LOG.info("*********TEST END*********")

    @staticmethod
    def report_name(report_path):
        u"""第三步：获取最新的测试报告"""
        lists = os.listdir(report_path)
        lists.sort(key=lambda fn: os.path.getmtime(os.path.join(report_path, fn)))
        print(u'最新测试生成的报告： ' + lists[-1])
        # 找到最新生成的报告文件
        report_file = os.path.join(report_path, lists[-1])
        return report_file

    # noinspection PyBroadException
    @staticmethod
    def send_mail(sender, psw, receiver, smtpserver, report_file, port):
        u"""第四步：发送最新的测试报告内容"""
        with open(report_file, "rb") as f:
            mail_body = f.read()
        # 定义邮件内容
        msg = MIMEMultipart()
        body = MIMEText(mail_body, _subtype='html', _charset='utf-8')
        msg['Subject'] = u"自动化测试报告"
        msg["from"] = sender
        msg["to"] = psw
        msg.attach(body)
        # 添加附件
        att = MIMEText(open(report_file, "rb").read(), "base64", "utf-8")
        att["Content-Type"] = "application/octet-stream"
        att["Content-Disposition"] = 'attachment; filename= "result.html"'
        msg.attach(att)
        # noinspection PyBroadException
        try:
            smtp = smtplib.SMTP_SSL(smtpserver, port)
        except:
            smtp = smtplib.SMTP()
            smtp.connect(smtpserver, port)
        # 用户名密码
        smtp.login(sender, psw)
        smtp.sendmail(sender, receiver, msg.as_string())
        smtp.quit()
        print('test report email has send out !')


if __name__ == '__main__':
    test = AllTest()
    job_id = str(os.getenv("BUILD_NUMBER"))
    test.run_testsuite(job_id)
    # 邮箱配置
    # regport_path = os.path.join(CUR_DIR, "report")
    # print(report_path)
    # report_file =obj.report_name(report_path)  # 3获取最新的测试报告
    # from config.read_config import Read_config
    # sender = Read_config().get_email('sender')
    # print(sender)
    # mail_password = Read_config().get_email('mail_password')
    # smtp_server = Read_config().get_email('mail_host')
    # port =Read_config().get_email('mail_port')
    # receiver = Read_config().get_email('receiver')
    # test.send_mail(sender, psw, receiver, smtp_server, report_file, port)  # 4最后一步发送报告
