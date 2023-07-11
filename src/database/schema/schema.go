package schema

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"xserver/src/logger"
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

type Pair[T, U any] struct {
	First  T
	Second U
}

type TableField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
}

type Table struct {
	Name       string       `json:"name"`
	Fields     []TableField `json:"fields"`
	PrimaryKey []string     `json:"primary_key"`
}

func Parse(schemaData []byte) ([]Table, error) {
	schema := &[]Table{}
	if err := json.Unmarshal(schemaData, schema); err != nil {
		return nil, err
	}

	return *schema, nil
}

func Verify(tables []Table) error {
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

func CreateTableCommand(table Table) string {
	fields := []string{}
	for _, field := range table.Fields {
		tableFieldType := fieldsTypesMap[field.Type]
		tableField := fmt.Sprintf("%s %s NOT NULL", field.Name, tableFieldType)
		if field.Nullable {
			tableField = fmt.Sprintf("%s %s NULL", field.Name, tableFieldType)
		}
		fields = append(fields, tableField)
	}

	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s(%s, PRIMARY KEY(%s))",
		table.Name,
		strings.Join(fields, ", "),
		strings.Join(table.PrimaryKey, ", "),
	)
}

func getPreviousSchema(db *sql.DB) ([]byte, error) {
	result, err := db.Query("SELECT data FROM __Schema WHERE version = 'current'")
	if err != nil {
		return nil, fmt.Errorf("[XServer] [Database] [Error] failed get previous schema from migration table: %s", err)
	}
	defer result.Close()

	data := []byte{}
	if result.Next() {
		if err := result.Scan(&data); err != nil {
			return nil, fmt.Errorf("[XServer] [Database] [Error] failed scan selected data from migration table: %s", err)
		}
	}
	if data == nil {
		return []byte{}, nil
	}

	return data, nil
}

func initMetadata(db *sql.DB, data []byte) error {
	table := Table{
		Name: "__Schema",
		Fields: []TableField{
			{
				Name:     "version",
				Type:     "string",
				Nullable: false,
			},
			{
				Name:     "data",
				Type:     "string",
				Nullable: true,
			},
		},
		PrimaryKey: []string{"version"},
	}

	if _, err := db.Exec(CreateTableCommand(table)); err != nil {
		return fmt.Errorf("[XServer] [Database] [Error] failed init migration table : %s", err)
	}

	return nil
}

func findTableFieldsDifference(currentTable Table, previousTable Table) ([]TableField, []TableField, []string) {
	currentTableFieldsMap := make(map[string]TableField)
	previousTableFieldsMap := make(map[string]TableField)

	for _, field := range previousTable.Fields {
		previousTableFieldsMap[field.Name] = field
	}

	for _, field := range currentTable.Fields {
		currentTableFieldsMap[field.Name] = field
	}

	newFields := []TableField{}
	removedFields := []TableField{}
	sameNamedFields := []Pair[TableField, TableField]{}
	sameFields := []string{}

	for _, field := range currentTable.Fields {
		previousField, ok := previousTableFieldsMap[field.Name]
		if ok {
			sameNamedFields = append(sameNamedFields, Pair[TableField, TableField]{field, previousField})
		}
		if !ok {
			newFields = append(newFields, field)
		}
	}

	for _, field := range previousTable.Fields {
		_, ok := currentTableFieldsMap[field.Name]
		if !ok {
			removedFields = append(removedFields, field)
		}
	}

	for _, fields := range sameNamedFields {
		if fields.First.Type != fields.Second.Type {
			newFields = append(newFields, fields.First)
			removedFields = append(removedFields, fields.Second)
		} else {
			sameFields = append(sameFields, fields.First.Name)
		}
	}
	return newFields, removedFields, sameFields
}

func findTablesDifference(currentSchema []Table, previousSchema []Table) ([]Table, []Table, []Pair[Table, Table]) {
	currentSchemaMap := make(map[string]Table)
	previousSchemaMap := make(map[string]Table)

	for _, table := range currentSchema {
		currentSchemaMap[table.Name] = table
	}

	for _, table := range previousSchema {
		previousSchemaMap[table.Name] = table
	}

	newTables := []Table{}
	removedTables := []Table{}
	sameTables := []Pair[Table, Table]{}
	for _, table := range currentSchema {
		previousTable, ok := previousSchemaMap[table.Name]
		if ok {
			sameTables = append(sameTables, Pair[Table, Table]{table, previousTable})
		}
		if !ok {
			newTables = append(newTables, table)
		}
	}

	for _, table := range previousSchema {
		_, ok := currentSchemaMap[table.Name]
		if !ok {
			removedTables = append(removedTables, table)
		}
	}

	return newTables, removedTables, sameTables
}

func schemaUpgrade(db *sql.DB, currentSchema []Table, previousSchema []Table) error {
	newTables, removedTables, sameTables := findTablesDifference(currentSchema, previousSchema)

	logger.Info("[XServer] [Database] [Schema] [Migration] transaction begin")

	if _, err := db.Exec("BEGIN TRANSACTION"); err != nil {
		return fmt.Errorf("[XServer] [Database] [Error] failed start schema upgrade transaction: %s", err)
	}

	for _, table := range newTables {
		logger.Info(fmt.Sprintf(`[XServer] [Database] [Schema] [Migration] create table "%s"`, table.Name))

		if _, err := db.Exec(CreateTableCommand(table)); err != nil {
			return fmt.Errorf("[XServer] [Database] [Error] failed create new table: %s", err)
		}
	}

	for _, table := range removedTables {
		logger.Info(fmt.Sprintf(`[XServer] [Database] [Schema] [Migration] remove table "%s"`, table.Name))

		if _, err := db.Exec(fmt.Sprintf("DROP TABLE %s", table.Name)); err != nil {
			logger.Error(fmt.Sprintf("[XServer] [Database] [Error] failed delete removed table: %s", err))
		}
	}

	for _, tables := range sameTables {
		table := tables.First
		table.Name = "_new_" + table.Name

		newFields, removedFields, sameFields := findTableFieldsDifference(tables.First, tables.Second)
		tableIsChanged := len(newFields) != 0 || len(removedFields) != 0

		if tableIsChanged {
			if _, err := db.Exec(CreateTableCommand(table)); err != nil {
				return fmt.Errorf(`[XServer] [Database] [Error] failed create "%s" table: %s`, table.Name, err)
			}

			selectFields := strings.Join(sameFields, ", ")
			if _, err := db.Exec(fmt.Sprintf("INSERT INTO %s (%s) SELECT %s FROM %s", table.Name, selectFields, selectFields, tables.Second.Name)); err != nil {
				return fmt.Errorf(`[XServer] [Database] [Error] failed copy data from "%s" table to "%s" table: %s`, tables.Second.Name, table.Name, err)
			}

			if _, err := db.Exec(fmt.Sprintf("DROP TABLE %s", tables.Second.Name)); err != nil {
				return fmt.Errorf(`[XServer] [Database] [Error] failed delete "%s" table: %s`, tables.Second.Name, err)
			}

			if _, err := db.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", table.Name, tables.Second.Name)); err != nil {
				return fmt.Errorf(`[XServer] [Database] [Error] failed rename "%s" table to "%s" table: %s`, table.Name, tables.Second.Name, err)
			}
		}
	}

	logger.Info("[XServer] [Database] [Schema] [Migration] transaction commit")

	if _, err := db.Exec("COMMIT TRANSACTION"); err != nil {
		return fmt.Errorf("[XServer] [Database] [Error] failed commit schema upgrade transaction: %s", err)
	}

	return nil
}

func Migration(db *sql.DB, currentSchemaData []byte) error {
	logger.Info("[XServer] [Database] [Schema] [Migration] start schema upgrade")

	if err := initMetadata(db, currentSchemaData); err != nil {
		return err
	}

	previousSchemaData, err := getPreviousSchema(db)
	if err != nil {
		return err
	}

	currentSchema, err := Parse(currentSchemaData)
	if err != nil {
		return fmt.Errorf("[XServer] [Database] [Error] failed parse current schema: %s", err)
	}

	previousSchema, err := Parse(previousSchemaData)
	if err != nil {
		return fmt.Errorf("[XServer] [Database] [Error] failed parse previous schema: %s", err)
	}

	if err := schemaUpgrade(db, currentSchema, previousSchema); err != nil {
		if _, err := db.Exec("ROLLBACK TRANSACTION"); err != nil {
			return fmt.Errorf("[XServer] [Database] [Error] failed rallback schema upgrade changes: %s", err)
		}
		return err
	}

	logger.Info("[XServer] [Database] [Schema] [Migration] save schema metadata")

	if _, err := db.Exec("DELETE FROM __Schema WHERE version = 'current'"); err != nil {
		return fmt.Errorf("[XServer] [Database] [Error] failed delete current schema from migration table: %s", err)
	}

	if _, err := db.Exec("INSERT INTO __Schema (version, data) VALUES ('current', $1)", currentSchemaData); err != nil {
		return fmt.Errorf("[XServer] [Database] [Error] failed insert current schema to migration table: %s", err)
	}

	logger.Info("[XServer] [Database] [Schema] [Migration] schema upgrade finished")

	return nil
}
