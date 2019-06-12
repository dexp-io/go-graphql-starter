package generate

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type Field struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Title string `json:"title"`
}
type EntityModel struct {
	Entity string  `json:"entity"`
	Table  string  `json:"table"`
	Fields []Field `json:"fields"`
}

func getFieldTitle(s string) string {

	strLen := len(s)
	for i := 0; i < strLen; i++ {
		if string(s[i]) == "_" {
			if i+1 < strLen {
				s = s[:i+1] + strings.ToUpper(string(s[i+1])) + s[i+2:]
			}
		}
	}
	s = strings.ToUpper(string(s[0])) + s[1:]
	s = strings.ReplaceAll(s, "_", "")
	return s
}

func Entities(projectPath, projectPackage string) {

	var packages = map[string]string{
		"database/sql":  "database/sql",
		"encoding/json": "encoding/json",
		"strconv":       "strconv",
		"strings":       "strings",
		"time":          "time",
		"errors":        "errors",
	}

	f, err := os.Create(projectPath + "/entities.go")
	if err != nil {
		panic(err)
	}

	jsonFile, err := os.Open(projectPath + "/generate/entities.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var entities []EntityModel
	_ = json.Unmarshal(byteValue, &entities)

	t, err := template.ParseFiles(projectPath + "/generate/template.go.tpl")

	var data []map[string]interface{}

	for _, entity := range entities {
		var fields []string
		var insertParams []string
		var insertValues [] string
		var selectFields [] string
		var scanFields []string

		shortEntityName := string(entity.Entity[0])

		for index, field := range entity.Fields {
			fieldTitle := getFieldTitle(field.Name)
			entity.Fields[index].Title = fieldTitle

			fields = append(fields, field.Name)
			selectFields = append(selectFields, shortEntityName+"."+field.Name)
			scanFields = append(scanFields, "&"+shortEntityName+"."+fieldTitle)
			insertParams = append(insertParams, "?")
			insertValues = append(insertValues, shortEntityName+"."+fieldTitle)

			switch field.Type {

			case "sql.NullInt64":
				packages["github.com/go-sql-driver/mysql"] = "github.com/go-sql-driver/mysql"
				break

			case "mysql.NullTime":
				packages["database/sql"] = "database/sql"
				break

			default:
				break
			}

		}

		// Execute the template to the file.

		insertFields := strings.Join(fields, ", ")
		updateFields := strings.Join(fields, " = ? , ")

		params := map[string]interface{}{
			"entity_title":  getFieldTitle(entity.Entity),
			"entity_name":   entity.Entity,
			"entity_table":  entity.Table,
			"entity_short":  shortEntityName,
			"fields":        entity.Fields,
			"insert_fields": insertFields,
			"insert_params": strings.Join(insertParams, ", "),
			"insert_values": strings.Join(insertValues, ", "),
			"update_fields": updateFields,
			"select_fields": strings.Join(selectFields, ", "),
			"scan_fields":   strings.Join(scanFields, ", "),
		}

		data = append(data, params)

	}

	params := map[string]interface{}{
		"package":  projectPackage,
		"packages": packages,
		"entities": data,
	}

	err = t.Execute(f, params)

	if err != nil {
		panic(err)
	}
	// Close the file when done.
	f.Close()

}
