import yaml
import os
import subprocess
import tenacity
from modules.database import Database
from modules.server import Server


class Project:
    def __init__(self, path: str = "example") -> None:
        self.path = path
        os.chdir(self.path)

        with open("config.yml", 'r') as config:
            try:
                self.config = yaml.safe_load(config)
            except yaml.YAMLError as exception:
                print(exception)

        self.url = "http://" + self.config["url"]

        self.database: Database = Database(self.url)
        self.database.clear()
        self.server: Server = Server(self.url)

        self.process = subprocess.Popen(
            "make", stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

    @tenacity.retry(wait=tenacity.wait_fixed(1), stop=tenacity.stop_after_delay(30))
    def wait_start(self):
        assert self.server.status() == "OK"

    def stop(self):
        self.process.terminate()
