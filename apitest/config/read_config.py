import codecs
import configparser
import os

CUR_DIR = os.path.dirname(os.path.realpath(__file__))
CONFIG_NAME = 'conf.ini'
CONFIG_FILE = os.path.join(CUR_DIR, CONFIG_NAME)


class ReadConfig:
    """
    Read configuration file
    """

    def __init__(self):
        config = open(CONFIG_FILE)
        data = config.read()

        if data[:3] == codecs.BOM_UTF8:
            data = data[3:]
            file = codecs.open(CONFIG_FILE, "w")
            file.write(data)
            file.close()
        config.close()

        self.configparser = configparser.ConfigParser()
        self.configparser.read(CONFIG_FILE)

    def get_email(self, name):
        value = self.configparser.get("EMAIL", name)
        return value

    def get_http(self, name):
        value = self.configparser.get("HTTP", name)
        return value

    def get_db(self, name):
        value = self.configparser.get("DATABASE", name)
        return value
