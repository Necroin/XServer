from time import sleep
import pytest

from modules.project import Project


class Environment():
    def __init__(self, project: Project) -> None:
        self.project: Project = project


@pytest.fixture(scope="session", name="project")
def project():
    project_instance = Project(
        path="example"
    )
    project_instance.wait_start()
    yield project_instance
    project_instance.stop()


@pytest.fixture(scope="session", name="environment")
def environment(project):
    yield Environment(
        project=project
    )
