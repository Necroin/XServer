import requests

response = requests.post("http://localhost:3301/go_handler", json={"a": 5, "b": 6})
print(response.text)