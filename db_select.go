package dexp

import (
	"database/sql"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type SelectQueryTable struct {
	Name  string
	Alias string
}

type SelectQueryField struct {
	Alias  string
	Fields []string
}

type SelectQuery struct {
	db          *Database
	table       *SelectQueryTable
	_fields     []*SelectQueryField
	_delimiter  string
	_conditions []*QueryCondition
	_whereSql   string
	_args       []interface{}
	_sql        string
	_orders     [] *QueryOrderBy
	_limit      *int
	_offset     *int
	_joins      [] *QueryJoin
}

func (q *SelectQuery) Join(table, alias, condition string) *SelectQuery {

	q._joins = append(q._joins, &QueryJoin{Table: table, Alias: alias, Condition: condition, Type: "INNER"})

	return q
}

func (q *SelectQuery) LeftJoin(table, alias, condition string) *SelectQuery {

	q._joins = append(q._joins, &QueryJoin{Table: table, Alias: alias, Condition: condition, Type: "LEFT"})

	return q
}

func (q *SelectQuery) RightJoin(table, alias, condition string) *SelectQuery {

	q._joins = append(q._joins, &QueryJoin{Table: table, Alias: alias, Condition: condition, Type: "RIGHT"})

	return q
}

func (q *SelectQuery) Fields(alias string, fields []string) *SelectQuery {

	q._fields = append(q._fields, &SelectQueryField{Alias: alias, Fields: fields})
	return q
}

func (q *SelectQuery) Range(limit, offset int) {

	q._limit = &limit
	q._offset = &offset
}

func (q *SelectQuery) OrderBy(field, value string) {
	q._orders = append(q._orders, &QueryOrderBy{
		Field: field,
		Value: value,
	})
}

func (q *SelectQuery) Or() {
	q._delimiter = "OR"
}

func (q *SelectQuery) And() {
	q._delimiter = "AND"
}

func (q *SelectQuery) Condition(field string, value interface{}, operator string) *SelectQuery {

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
				q._args = append(q._args, s.Index(i))
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

func (q *SelectQuery) exec() {

	if q._sql != "" {
		return
	}

	var selectFields [] string

	for _, field := range q._fields {
		for i := 0; i < len(field.Fields); i++ {
			selectFields = append(selectFields, field.Alias+"."+field.Fields[i])
		}

	}

	s := "SELECT " + strings.Join(selectFields, ", ") + " FROM `" + q.table.Name + "` AS " + q.table.Alias

	var joins []string

	if len(q._joins) > 0 {

		for _, j := range q._joins {
			joins = append(joins, j.Type+" JOIN `"+j.Table+"` AS "+j.Alias+" ON "+j.Condition)
		}

		s += " " + strings.Join(joins, " ")
	}

	if q._whereSql != "" {
		s += " WHERE " + q._whereSql
	}

	var orders []string

	if len(q._orders) > 0 {
		for _, o := range q._orders {
			orders = append(orders, o.Field+" "+o.Value)
		}
		s += " ORDER BY " + strings.Join(orders, ", ")
	}

	if q._limit != nil {
		s += " LIMIT " + strconv.Itoa(*q._limit) + " OFFSET " + strconv.Itoa(*q._offset)
	}
	q._sql = s

}

func (q *SelectQuery) ToSql() string {
	q.exec()
	return q._sql
}

func (q *SelectQuery) FetchOne() *sql.Row {

	q.exec()

	if len(q._args) > 0 {
		return q.db.QueryRow(q._sql, q._args...)
	}

	log.Println("sql", q._sql)
	return q.db.QueryRow(q._sql)
}

func (q *SelectQuery) FetchAll() (*sql.Rows, error) {

	q.exec()
	if len(q._args) > 0 {
		return q.db.Query(q._sql, q._args...)
	}

	return q.db.Query(q._sql)

}
