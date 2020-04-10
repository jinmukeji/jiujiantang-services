import os
from functools import wraps

import logbook
from logbook.more import ColorizedStderrHandler

# 创建日志路径
CUR_DIR = os.path.dirname(os.path.realpath(__file__))
LOG_DIR = os.path.join(CUR_DIR, 'log')
FILE_STREAM = False
if not os.path.exists(LOG_DIR):
    os.makedirs(LOG_DIR)
    FILE_STREAM = True


def get_logger(name='interface', level=''):
    """ get logger Factory function """
    logbook.set_datetime_format('local')
    # 打印到屏幕句柄

    ColorizedStderrHandler(bubble=False, level=level).push_thread()
    # 打印到文件句柄
    logbook.TimedRotatingFileHandler(
        os.path.join(LOG_DIR, '%s.log' % name),
        date_format='%Y-%m-%d-%H', bubble=True, encoding='utf-8').push_thread()
    return logbook.Logger(name)


LOG = get_logger(level='INFO')


def logger(param):
    """ function from logger meta """

    def wrap(function):
        """ logger wrapper """

        @wraps(function)
        def _wrap(*args, **kwargs):
            """ wrap tool """
            LOG.info("运行位置 {}".format(param))
            # LOG.info("全部args参数参数信息 , {}".format(str(args)))
            # LOG.info("全部kwargs参数信息 , {}".format(str(kwargs)))
            return function(*args, **kwargs)

        return _wrap

    return wrap
