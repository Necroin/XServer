from conftest import Environment


def test_empty_schema(environment: Environment):
    response = environment.project.database.set_schema([])
    assert response["result"] == True
    assert response.get("error", None) is None


def test_wrong_schema_empty_table_name(environment: Environment):
    response = environment.project.database.set_schema([{}])
    assert response["result"] == False
    assert response["error"] == '[XServer] [Database] [Error] failed verify schema: table name is empty'


def test_wrong_schema_empty_fields(environment: Environment):
    response = environment.project.database.set_schema([{"name": "Test"}])
    assert response["result"] == False
    assert response["error"] == '[XServer] [Database] [Error] failed verify schema: missed "fields" section'

    response = environment.project.database.set_schema([
        {
            "name": "Test",
            "fields": []
        }
    ])
    assert response["result"] == False
    assert response["error"] == '[XServer] [Database] [Error] failed verify schema: missed "fields" section'


def test_wrong_schema_empty_field_name(environment: Environment):
    response = environment.project.database.set_schema([
        {
            "name": "Test",
            "fields": [
                {}
            ]
        }
    ])
    assert response["result"] == False
    assert response["error"] == '[XServer] [Database] [Error] failed verify schema: empty field name in "Test" table'


def test_wrong_schema_empty_field_type(environment: Environment):
    response = environment.project.database.set_schema([
        {
            "name": "Test",
            "fields": [
                {
                    "name": "field1"
                }
            ]
        }
    ])
    assert response["result"] == False
    assert response["error"] == '[XServer] [Database] [Error] failed verify schema: unknown type for "field1" field in "Test" table'


def test_wrong_schema_unknown_field_type(environment: Environment):
    response = environment.project.database.set_schema([
        {
            "name": "Test",
            "fields": [
                {
                    "name": "field1",
                    "type": "unknown"
                }
            ]
        }
    ])
    assert response["result"] == False
    assert response["error"] == '[XServer] [Database] [Error] failed verify schema: unknown type for "field1" field in "Test" table'


def test_wrong_schema_empty_primary_key(environment: Environment):
    response = environment.project.database.set_schema([
        {
            "name": "Test",
            "fields": [
                {
                    "name": "field1",
                    "type": "string"
                }
            ]
        }
    ])
    assert response["result"] == False
    assert response["error"] == '[XServer] [Database] [Error] failed verify schema: primary key for "Test" table is empty'

    response = environment.project.database.set_schema([
        {
            "name": "Test",
            "fields": [
                {
                    "name": "field1",
                    "type": "string"
                }
            ],
            "primary_key": []
        }
    ])
    assert response["result"] == False
    assert response["error"] == '[XServer] [Database] [Error] failed verify schema: primary key for "Test" table is empty'


def test_wrong_schema_primary_key_unknown_field(environment: Environment):
    response = environment.project.database.set_schema([
        {
            "name": "Test",
            "fields": [
                {
                    "name": "field1",
                    "type": "string"
                }
            ],
            "primary_key": ["unknown_field"]
        }
    ])
    assert response["result"] == False
    assert response["error"] == '[XServer] [Database] [Error] failed verify schema: unknown field "unknown_field" in primary key for "Test" table'


def test_schema(environment: Environment):
    response = environment.project.database.set_schema([
        {
            "name": "Test",
            "fields": [
                {
                    "name": "field1",
                    "type": "string"
                }
            ],
            "primary_key": ["field1"]
        }
    ])
    assert response["result"] == True
    assert response.get("error", None) is None


def test_schema_migration_new_field(environment: Environment):
    response = environment.project.database.set_schema([
        {
            "name": "Test",
            "fields": [
                {
                    "name": "field1",
                    "type": "string"
                },
                {
                    "name": "field2",
                    "type": "string"
                }
            ],
            "primary_key": ["field1"]
        }
    ])
    assert response["result"] == True
    assert response.get("error", None) is None

    response = environment.project.database.select({
        "table": "Test",
        "fields": [{"name": "field1"}, {"name": "field2"}]
    })
    
    assert response["result"] == []
    assert response.get("error", None) is None
