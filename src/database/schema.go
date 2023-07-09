package database

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

var (
	fieldsTypesMap = map[string]string{
		"null":      "null",
		"int":       "integer",
		"integer":   "integer",
		"float":     "float",
		"string":    "text",
		"timestamp": "timestamp",
		"datetime":  "datetime",
	}
)

type TableField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Table struct {
	Name       string       `json:"name"`
	Fields     []TableField `json:"fields"`
	PrimaryKey []string     `json:"primary_key"`
}

func parseSchema(schemaFile io.Reader) ([]Table, error) {
	schema := &[]Table{}
	if err := json.NewDecoder(schemaFile).Decode(schema); err != nil {
		return nil, err
	}

	return *schema, nil
}

func verifySchema(tables []Table) error {
	for _, table := range tables {
		if table.Name == "" {
			return fmt.Errorf("table name is empty")
		}

		if len(table.Fields) == 0 {
			return fmt.Errorf(`missed "fields" section`)
		}

		fieldsMap := make(map[string]bool)
		for _, field := range table.Fields {
			if field.Name == "" {
				return fmt.Errorf(`empty field name in "%s" table`, table.Name)
			}
			_, ok := fieldsTypesMap[field.Type]
			if !ok {
				return fmt.Errorf(`unknown type for "%s" field in "%s" table`, field.Name, table.Name)
			}
			fieldsMap[field.Name] = true
		}

		if len(table.PrimaryKey) == 0 {
			return fmt.Errorf(`primary key for "%s" table is empty`, table.Name)
		}

		for _, fieldName := range table.PrimaryKey {
			_, ok := fieldsMap[fieldName]
			if !ok {
				return fmt.Errorf(`unknown field "%s" in primary key for "%s" table`, fieldName, table.Name)
			}
		}

	}

	return nil
}

func createTableCommand(table Table) string {
	fields := []string{}
	for _, field := range table.Fields {
		tableFieldType := fieldsTypesMap[field.Type]
		tableField := fmt.Sprintf("%s %s NOT NULL", field.Name, tableFieldType)
		fields = append(fields, tableField)
	}

	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s(%s, PRIMARY KEY(%s))",
		table.Name,
		strings.Join(fields, ", "),
		strings.Join(table.PrimaryKey, ", "),
	)
}
