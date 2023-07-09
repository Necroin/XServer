
import requests
import sqlite3


class Database:
    def __init__(self, server_url, storage_path="storage.db", server_db_path="/db/") -> None:
        self.server_url = server_url
        self.storage_path = storage_path
        self.server_db_path = server_db_path
        pass

    def request(self, operation: str, data):
        response = requests.post(
            self.server_url + self.server_db_path+operation, json=data)
        response.raise_for_status()
        return response.json()

    def insert(self, data: dict):
        return self.request("insert", data)

    def select(self, data: dict):
        return self.request("select", data)

    def update(self, data: dict):
        return self.request("update", data)

    def delete(self, data: dict):
        return self.request("delete", data)

    def set_schema(self, data: list):
        return self.request("set_schema", data)
        
    def clear(self):
        connection = sqlite3.connect(self.storage_path)
        connection.execute("DELETE FROM Users")
        connection.commit()
