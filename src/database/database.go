package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"xserver/src/config"
	"xserver/src/logger"

	_ "github.com/mattn/go-sqlite3"
)

type RequestFilter struct {
	Name     string `json:"name"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type RequestField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Request struct {
	Table   string          `json:"table"`
	Fields  []RequestField  `json:"fields"`
	Filters []RequestFilter `json:"filters"`
}

type Database struct {
	db *sql.DB
}

func Create(config *config.Config) (*Database, error) {
	db, err := sql.Open("sqlite3", config.Database.Storage)
	if err != nil {
		return nil, fmt.Errorf("[XServer] [Database] [Error] failed open database: %s", err)
	}

	database := &Database{
		db: db,
	}

	schemaFile, err := os.Open(config.Database.Schema)
	if err != nil {
		return nil, fmt.Errorf("[XServer] [Database] [Error] failed open schema file: %s", err)
	}

	if err := database.SetSchema(schemaFile); err != nil {
		return nil, err
	}

	return database, nil
}

func (database *Database) Close() {
	database.db.Close()
}

func (database *Database) SetSchema(data io.Reader) error {
	schema, err := parseSchema(data)
	if err != nil {
		return fmt.Errorf("[XServer] [Database] [Error] failed parse schema file: %s", err)
	}
	if err := verifySchema(schema); err != nil {
		return fmt.Errorf("[XServer] [Database] [Error] failed verify schema: %s", err)
	}

	for _, table := range schema {
		command := createTableCommand(table)
		logger.Debug(fmt.Sprintf("[XServer] [Database] %s", command))
		_, err := database.db.Exec(command)
		if err != nil {
			return fmt.Errorf("[XServer] [Database] [Error] failed create table : %s", err)
		}
	}

	return nil
}

func (database *Database) Insert(data io.Reader, responseWriter io.Writer) error {
	request := &Request{}
	if err := json.NewDecoder(data).Decode(request); err != nil {
		return fmt.Errorf("[XServer] [Database] [Insert] [Error] failed decode json request: %s", err)
	}

	names := []string{}
	for _, field := range request.Fields {
		names = append(names, field.Name)
	}

	values := []string{}
	for _, field := range request.Fields {
		values = append(values, field.Value)
	}

	sqlCommand := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", request.Table, strings.Join(names, ", "), strings.Join(values, ", "))
	logger.Debug(fmt.Sprintf("[XServer] [Database] [Insert] sql request: %s", sqlCommand))

	_, err := database.db.Exec(sqlCommand)
	if err != nil {
		return fmt.Errorf("[XServer] [Database] [Insert] [Error] failed database request: %s", err)
	}

	responseWriter.Write([]byte(`{"result": true}`))
	return nil
}

func (database *Database) Select(data io.Reader, responseWriter io.Writer) error {
	request := &Request{}
	if err := json.NewDecoder(data).Decode(request); err != nil {
		return fmt.Errorf("[XServer] [Database] [Select] [Error] failed decode json request: %s", err)
	}

	sqlCommand := fmt.Sprintf("SELECT * FROM %s", request.Table)

	if len(request.Fields) != 0 {
		fields := []string{}
		for _, field := range request.Fields {
			fields = append(fields, field.Name)
		}
		sqlCommand = fmt.Sprintf("SELECT %s FROM %s", strings.Join(fields, ", "), request.Table)
	}

	if len(request.Filters) != 0 {
		filters := []string{}
		for _, filter := range request.Filters {
			filters = append(filters, fmt.Sprintf("%s %s %s", filter.Name, filter.Operator, filter.Value))
		}
		sqlFilters := " WHERE " + strings.Join(filters, " AND ")
		sqlCommand = sqlCommand + sqlFilters
	}
	logger.Debug(fmt.Sprintf("[XServer] [Database] [Select] sql request: %s", sqlCommand))

	result, err := database.db.Query(sqlCommand)
	if err != nil {
		return fmt.Errorf("[XServer] [Database] [Select] [Error] failed database request: %s", err)
	}
	defer result.Close()

	columns, err := result.Columns()
	if err != nil {
		return fmt.Errorf("[XServer] [Database] [Select] [Error] failed get result columns: %s", err)
	}

	records := []string{}

	for result.Next() {
		values := make([]string, len(columns))
		valuesPointers := make([]interface{}, len(columns))
		for i := range values {
			valuesPointers[i] = &values[i]
		}

		if err := result.Scan(valuesPointers...); err != nil {
			return fmt.Errorf("[XServer] [Database] [Select] [Error] failed scan row values: %s", err)
		}

		record := []string{}

		for i, column := range columns {
			record = append(record, fmt.Sprintf(`"%s": "%s"`, column, values[i]))
		}

		records = append(records, fmt.Sprintf("{%s}", strings.Join(record, ", ")))
	}

	responseWriter.Write([]byte(fmt.Sprintf(`{"result": [%s]}`, strings.Join(records, ", "))))

	return nil
}

func (database *Database) Update(data io.Reader, responseWriter io.Writer) error {
	request := &Request{}
	if err := json.NewDecoder(data).Decode(request); err != nil {
		return fmt.Errorf("[XServer] [Database] [Update] [Error] failed decode json request: %s", err)
	}

	sqlCommand := fmt.Sprintf("UPDATE %s SET ", request.Table)

	fields := []string{}
	for _, field := range request.Fields {
		fields = append(fields, fmt.Sprintf("%s = %s", field.Name, field.Value))
	}
	sqlCommand = sqlCommand + strings.Join(fields, ", ")

	if len(request.Filters) != 0 {
		filters := []string{}
		for _, filter := range request.Filters {
			filters = append(filters, fmt.Sprintf("%s %s %s", filter.Name, filter.Operator, filter.Value))
		}
		sqlFilters := " WHERE " + strings.Join(filters, " AND ")
		sqlCommand = sqlCommand + sqlFilters
	}
	logger.Debug(fmt.Sprintf("[XServer] [Database] [Update] sql request: %s", sqlCommand))

	_, err := database.db.Exec(sqlCommand)
	if err != nil {
		return fmt.Errorf("[XServer] [Database] [Update] [Error] failed database request: %s", err)
	}

	responseWriter.Write([]byte(`{"result": true}`))

	return nil
}

func (database *Database) Delete(data io.Reader, responseWriter io.Writer) error {
	request := &Request{}
	if err := json.NewDecoder(data).Decode(request); err != nil {
		return fmt.Errorf("[XServer] [Database] [Delete] [Error] failed decode json request: %s", err)
	}

	sqlCommand := fmt.Sprintf("DELETE FROM %s", request.Table)

	if len(request.Filters) != 0 {
		filters := []string{}
		for _, filter := range request.Filters {
			filters = append(filters, fmt.Sprintf("%s %s %s", filter.Name, filter.Operator, filter.Value))
		}
		sqlFilters := " WHERE " + strings.Join(filters, " AND ")
		sqlCommand = sqlCommand + sqlFilters
	}
	logger.Debug(fmt.Sprintf("[XServer] [Database] [Delete] sql request: %s", sqlCommand))

	_, err := database.db.Exec(sqlCommand)
	if err != nil {
		return fmt.Errorf("[XServer] [Database] [Delete] [Error] failed database request: %s", err)
	}

	responseWriter.Write([]byte(`{"result": true}`))

	return nil
}
