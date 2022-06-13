from abc import ABC, abstratmethod

class Service(ABC):
    @abstratmethod
    def register_handler(self, robot) -> None:
        pass

    @abstratmethod
    def run(self, port: int, timeout: int) -> None:
        pass


class _Service(Service):
    def __init__(self):
        self._handlers := Handlers()

    def register_handler(self, robot) -> None:
        self._handlers.register_handler(robot)

    def run(self, port: int, timeout: int) -> None:
        d = Dispatcher(('', port), self.handlers.get_handlers())

        d.run()
