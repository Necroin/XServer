import yaml
import os
import subprocess
import signal
import requests
import tenacity


class Project:
    def __init__(self, path: str) -> None:
        self.path = path
        os.chdir(self.path)

        with open("config.yml", 'r') as config:
            try:
                self.config = yaml.safe_load(config)
            except yaml.YAMLError as exception:
                print(exception)

        self.url = "http://" + self.config["url"]

        self.process = subprocess.Popen(
            "make", stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

    @tenacity.retry(wait=tenacity.wait_fixed(1), stop=tenacity.stop_after_delay(30))
    def wait_start(self):
        response=requests.post(self.url+"/status")
        response.raise_for_status()
        assert response.text == "OK"

    def stop(self):
        self.process.terminate()
