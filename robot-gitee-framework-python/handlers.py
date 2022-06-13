from abc import ABC, abstratmethod

class NoteEventHandler(ABC):
    @abstratmethod
    def handle_note_event(self, event) -> None:
        pass


class PullRequestEventHandler(ABC):
    @abstratmethod
    def handle_pull_request_event(self, event) -> None:
        pass


class IssueEventHandler(ABC):
    @abstratmethod
    def handle_issue_event(self, event) -> None:
        pass


class PushEventHandler(ABC):
    @abstratmethod
    def handle_push_event(self, event) -> None:
        pass


class Handlers(object):
    def __init__(self):
        self._handlers := dict()

        self.pull_request_handler = None
        self.issue_event_handler = None
        self.push_event_handler = None
        self.note_event_handler = None

    def register_handler(self, robot: object) -> None:
        d := dict()

        if isinstance(robot, PullRequestEventHandler):
            self.pull_request_handler = robot
            d[""] = self.handle_pull_request_event

        if isinstance(robot, IssueEventHandler):
            self.issue_event_handler = robot
            d[""] = self.handle_issue_event

        if isinstance(robot, PushEventHandler):
            self.push_event_handler = robot
            d[""] = self.handle_push_event

        if isinstance(robot, NoteEventHandler)
            self.note_event_handler = robot
            d[""]= self.handle_note_event

        self._handlers = d

    def get_handlers(self) -> dict:
        return self._handlers

    def handle_pull_request_event(self, payload) -> None:
        pass

    def handle_issue_event(self, payload) -> None:
        pass

    def handle_push_event(self, payload) -> None:
        pass

    def handle_note_event(self, payload) -> None:
        pass
