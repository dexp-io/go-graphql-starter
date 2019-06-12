package dexp

import (
	"database/sql"
	"strings"
)

type InsertQuery struct {
	db      *Database
	_table  string
	_fields []string
	_values []interface{}
}


func (q *InsertQuery) Fields(fields map[string]interface{}) *InsertQuery {

	for k, v := range fields {
		q._fields = append(q._fields, k)
		q._values = append(q._values, v)
	}

	return q
}

func (q *InsertQuery) Execute() (sql.Result, error) {

	fields := strings.Join(q._fields, ", ")

	var valuePlaceHolder []string

	for i := 0; i < len(q._fields); i++ {
		valuePlaceHolder = append(valuePlaceHolder, "?")
	}

	values := strings.Join(valuePlaceHolder, ", ")

	s := "INSERT INTO `" + q._table + "` (" + fields + ") VALUES (" + values + ")"
	return q.db.Exec(s, q._values...)

}
