
import requests


class Server:
    def __init__(self, url) -> None:
        self.url = url

    def status(self):
        response = requests.post(self.url+"/status")
        response.raise_for_status()
        return response.text

    def request(self, path, data):
        response = requests.post(self.url + path, json=data)
        return response.text
