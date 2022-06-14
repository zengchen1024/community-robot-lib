import logging


class Log(logging.Logger):

    def __init__(self, name, log_level=logging.DEBUG, log_format=None, data_format=None):
        super().__init__(name)

        if not log_format:
            log_format='[%(asctime)s][%(levelname)s][%(filename)s:%(lineno)s]%(message)s',

        if not data_format:
            data_format='%Y-%m-%d %H:%M:%S'

        formatter = logging.Formatter(log_format, data_format)

        console = logging.StreamHandler()
        console.setFormatter(formatter)
        console.setLevel(log_level)

        self.addHandler(console)
