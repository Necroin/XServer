import requests
from conftest import Environment


def test_insert(environment: Environment):
    environment.project.database.clear()

    response = environment.project.database.insert(
        {
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

    assert response["result"] == True
    assert response.get("error", None) is None

    response = environment.project.database.insert(
        {
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

    assert response["result"] == True
    assert response.get("error", None) is None


def test_select_all(environment: Environment):
    response = environment.project.database.select(
        {
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

    assert response == expected


def test_select_with_filter(environment: Environment):
    response = environment.project.database.select({
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

    assert response == expected


def test_update_all(environment: Environment):
    response = environment.project.database.update(
        {
            "table": "Users",
            "fields": [
                {
                    "name": "age",
                    "value": "100"
                }
            ]
        }
    )

    assert response["result"] == True
    assert response.get("error", None) is None

    response = environment.project.database.select(
        {
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

    assert response == expected


def test_update_with_filter(environment: Environment):
    response = environment.project.database.update(
        {
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

    assert response["result"] == True
    assert response.get("error", None) is None

    response = environment.project.database.select(
        {
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

    assert response == expected


def test_delete(environment: Environment):
    response = environment.project.database.delete(
        {
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

    assert response["result"] == True
    assert response.get("error", None) is None

    response = environment.project.database.select(
        {
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

    assert response == expected
