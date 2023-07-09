from conftest import Environment


def test_go_handler(environment: Environment):
    assert environment.project.server.request(
        "/go_handler", {"a": 5, "b": 6}) == '[Go Handler] Started\n{"a": 5, "b": 6}\n'


def test_cpp_handler(environment: Environment):
    assert environment.project.server.request(
        "/cpp_handler", {"a": 5, "b": 6}) == '[C++ Handler] Started\n{"a": 5, "b": 6}\n'


def test_python_handler(environment: Environment):
    assert environment.project.server.request(
        "/python_handler", {"a": 5, "b": 6}) == '[Python Handler] Started\n{"a": 5, "b": 6}\n'


def test_lua_handler(environment: Environment):
    assert environment.project.server.request(
        "/lua_handler", {"a": 5, "b": 6}) == '[Lua Handler] Started\n{"a": 5, "b": 6}\n'
