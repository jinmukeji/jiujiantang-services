import pymysql

from common.log import LOG
from config import read_config

CUR_CONFIG = read_config.ReadConfig()


class Database:

    def __init__(self):
        self.conn = None
        self.cursor = None

    def connect(self):
        """
        连接数据库
        :return:
        """
        host = CUR_CONFIG.get_db("host")
        username = CUR_CONFIG.get_db("username")
        password = CUR_CONFIG.get_db("password")
        port = CUR_CONFIG.get_db("port")
        database = CUR_CONFIG.get_db("database")

        try:
            # connect to DB
            conn = pymysql.connect(host=host, port=int(port), user=username, password=password, db=database)
            self.cursor = conn.cursor()
            LOG.info("Connect DB successfully!")
        except ConnectionError as ex:
            LOG.error(ex)

    def close(self):
        """
        关闭数据库连接
        :param self:
        :return:
        """
        self.conn.close()
        LOG.info("Database closed!")

    def execute(self, sql):
        """
        操作数据
        :param sql:
        :param self:
        :return:
        """
        self.connect()
        self.cursor.execute(sql)
        return self.cursor

    @staticmethod
    def fetch_all(cursor):
        """
        获取所有数据
        :param cursor:
        :return:
        """
        value = cursor.fetchall()
        return value

    @staticmethod
    def fetch_one(cursor):
        """
        获取单条数据
        :param cursor:
        :return:
        """
        result = cursor.fetchone()
        return result
