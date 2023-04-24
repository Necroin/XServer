import requests

def test_go_handler():
    response = requests.post("http://localhost:3301/go_handler", json={"a": 5, "b": 6})
    assert response.text == '[Go Handler] Started\n{"a": 5, "b": 6}\n'

def test_cpp_handler():
    response = requests.post("http://localhost:3301/cpp_handler", json={"a": 5, "b": 6})
    assert response.text == '[C++ Handler] Started\n{"a": 5, "b": 6}\n'
    
def test_python_handler():
    response = requests.post("http://localhost:3301/python_handler", json={"a": 5, "b": 6})
    assert response.text == '[Python Handler] Started\n{"a": 5, "b": 6}\n'
    
if __name__ == "__main__":
    test_go_handler()
    test_cpp_handler()
    test_python_handler()