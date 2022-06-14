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
        self.pull_request_handler = None
        self.issue_event_handler = None
        self.push_event_handler = None
        self.note_event_handler = None

    def register_handler(self, robot: object) -> None:
        if isinstance(robot, PullRequestEventHandler):
            self.pull_request_handler = robot

        if isinstance(robot, IssueEventHandler):
            self.issue_event_handler = robot

        if isinstance(robot, PushEventHandler):
            self.push_event_handler = robot

        if isinstance(robot, NoteEventHandler)
            self.note_event_handler = robot

    def get_handlers(self) -> dict:
        d := dict()

        if self.pull_request_handler != None:
            d[""] = self.handle_pull_request_event

        if self.issue_event_handler != None:
            d[""] = self.handle_issue_event

        if self.push_event_handler != None:
            d[""] = self.handle_push_event

        if self.note_event_handler != None:
            d[""]= self.handle_note_event

        return d

    def handle_pull_request_event(self, payload) -> None:
        pass

    def handle_issue_event(self, payload) -> None:
        pass

    def handle_push_event(self, payload) -> None:
        pass

    def handle_note_event(self, payload) -> None:
        pass
