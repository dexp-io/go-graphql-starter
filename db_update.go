package dexp

import (
	"database/sql"
	"log"
	"reflect"
	"strings"
)

type UpdateQuery struct {
	db          *Database
	_table      string
	_fields     []string
	_whereSql   string
	_args       []interface{}
	_conditions []*QueryCondition
	_delimiter  string
}

func (q *UpdateQuery) Fields(fields map[string]interface{}) *UpdateQuery {

	for k, v := range fields {
		q._fields = append(q._fields, k)
		q._args = append(q._args, v)
	}

	return q
}

func (q *UpdateQuery) Condition(field string, value interface{}, operator string) *UpdateQuery {

	q._conditions = append(q._conditions, &QueryCondition{Field: field, Value: value, Operator: operator})

	if q._delimiter == "OR" {
		q._whereSql += " OR "
		q._delimiter = "AND"
	} else {
		if len(q._conditions) > 1 {
			q._whereSql += " AND "
		}
	}

	q._whereSql += field + " "
	switch operator {

	case "IN":

		var in []string
		switch reflect.TypeOf(value).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			for i := 0; i < s.Len(); i++ {
				in = append(in, "?")
				q._args = append(q._args, s.Index(i).String())
			}
		}

		q._whereSql += operator + "(" + strings.Join(in, ", ") + ")"

		break

	default:

		q._whereSql += operator + " ?"

		q._args = append(q._args, value)

		break
	}

	return q

}

func (q *UpdateQuery) Execute() (sql.Result, error) {

	fields := ""

	for i, field := range q._fields {
		if i > 0 && i < len(q._fields) {
			fields += ", "
		}
		fields += field + "=?"

	}
	s := "UPDATE  `" + q._table + "` SET " + fields
	if q._whereSql != "" {
		s += " WHERE " + q._whereSql
	}

	log.Println("s", s)
	return q.db.Exec(s, q._args...)
}
