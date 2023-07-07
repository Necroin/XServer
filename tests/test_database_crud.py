import requests
import sqlite3
from conftest import Environment


def clear_database():
    connection = sqlite3.connect("storage.db")
    connection.execute("DELETE FROM Users")
    connection.commit()


def test_insert(environment: Environment):
    clear_database()

    response = requests.post(
        environment.project.url + "/db/insert",
        json={
            "table": "Users",
            "fields": [
                {
                    "name": "name",
                    "value": "'Me'"
                },
                {
                    "name": "age",
                    "value": "20"
                }
            ]
        }
    )

    assert response.json()["result"] == True
    assert response.json().get("error", None) is None

    response = requests.post(
        environment.project.url + "/db/insert",
        json={
            "table": "Users",
            "fields": [
                {
                    "name": "name",
                    "value": "'Other'"
                },
                {
                    "name": "age",
                    "value": "50"
                }
            ]
        }
    )

    assert response.json()["result"] == True
    assert response.json().get("error", None) is None


def test_select_all(environment: Environment):
    response = requests.post(
        environment.project.url + "/db/select",
        json={
            "table": "Users",
            "fields": [{"name": "name"}, {"name": "age"}]
        }
    )

    expected = {
        "result": [
            {
                "name": "Me",
                "age": "20"
            },
            {
                "name": "Other",
                "age": "50"
            }
        ]
    }

    assert response.json() == expected


def test_select_with_filter(environment: Environment):
    response = requests.post(
        environment.project.url + "/db/select",
        json={
            "table": "Users",
            "fields": [{"name": "name"}, {"name": "age"}],
            "filters": [
                {
                    "name": "name",
                    "operator": "=",
                    "value": "'Me'"
                }
            ]
        }
    )

    expected = {
        "result": [
            {
                "name": "Me",
                "age": "20"
            }
        ]
    }

    assert response.json() == expected


def test_update_all(environment: Environment):
    response = requests.post(
        environment.project.url + "/db/update",
        json={
            "table": "Users",
            "fields": [
                {
                    "name": "age",
                    "value": "100"
                }
            ]
        }
    )

    assert response.json()["result"] == True
    assert response.json().get("error", None) is None

    response = requests.post(
        environment.project.url + "/db/select",
        json={
            "table": "Users",
            "fields": [{"name": "name"}, {"name": "age"}]
        }
    )

    expected = {
        "result": [
            {
                "name": "Me",
                "age": "100"
            },
            {
                "name": "Other",
                "age": "100"
            }
        ]
    }

    assert response.json() == expected


def test_update_with_filter(environment: Environment):
    response = requests.post(
        environment.project.url + "/db/update",
        json={
            "table": "Users",
            "filters": [
                {
                    "name": "name",
                    "operator": "=",
                    "value": "'Me'"
                }
            ],
            "fields": [
                {
                    "name": "age",
                    "value": "200"
                }
            ]
        }
    )

    assert response.json()["result"] == True
    assert response.json().get("error", None) is None

    response = requests.post(
        environment.project.url + "/db/select",
        json={
            "table": "Users",
            "fields": [{"name": "name"}, {"name": "age"}]
        }
    )

    expected = {
        "result": [
            {
                "name": "Me",
                "age": "200"
            },
            {
                "name": "Other",
                "age": "100"
            }
        ]
    }

    assert response.json() == expected


def test_delete(environment: Environment):
    response = requests.post(
        "http://localhost:3301/db/delete",
        json={
            "table": "Users",
            "filters": [
                {
                    "name": "name",
                    "operator": "=",
                    "value": "'Other'"
                }
            ]
        }
    )

    assert response.json()["result"] == True
    assert response.json().get("error", None) is None

    response = requests.post(
        "http://localhost:3301/db/select",
        json={
            "table": "Users",
            "fields": [{"name": "name"}, {"name": "age"}]
        }
    )

    expected = {
        "result": [
            {
                "name": "Me",
                "age": "200"
            }
        ]
    }

    assert response.json() == expected
