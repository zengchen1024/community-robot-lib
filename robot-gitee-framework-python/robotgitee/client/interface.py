from abc import ABC, abstractmethod


class Client(ABC):
    @abstractmethod
    def create_pr_comment(self, org: str, repo: str, number: int, comment: str):
        pass


