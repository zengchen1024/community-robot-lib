def run(robot: object, port: int, timeout: int) -> None:
    h = Handlers()
    h.register_handler(robot)

    hs = self.handlers.get_handlers()
    if len(hs) == 0:
        raise Exception("it is not a robot")

    d = Dispatcher(('', port), hs)
    d.run()
