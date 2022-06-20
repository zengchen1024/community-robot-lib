import gitee

from robotgitee.client.interface import Client


class _Client(Client):
    def __init__(self, token: str):
        cli = gitee.ApiClient(header_name="Authorization",
                              header_value="Bearer " + token)

        self.prapi = gitee.PullRequestsApi(cli)

    def create_pr_comment(self, org: str, repo: str, number: int,
                          comment: str):
        body = gitee.PullRequestCommentPostParam(body=comment)

        self.prapi.post_v5_repos_owner_repo_pulls_number_comments(
            org, repo, number, body)


def new_client(token: str) -> Client:
    return _Client(token)
