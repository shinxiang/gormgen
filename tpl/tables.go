package tpl

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"unicode"

	"github.com/jinzhu/inflection"
	"github.com/pkg/errors"
)

const (
	showTablesSQL      = "SHOW TABLES"
	showCreateTableSQL = "SHOW CREATE TABLE %s"
	primaryKeyMark     = "PRIMARY KEY"
	commentMark        = "COMMENT"
	unsignedKeyFlag    = "unsigned"
)

var tableToGoStruct = map[string]string{
	"bigint":    "int64",
	"int":       "int",
	"smallint":  "int",
	"tinyint":   "int",
	"mediumint": "int",
	"decimal":   "float64",
	"numeric":   "float64",
	"float":     "float64",
	"datetime":  "time.Time",
	"date":      "time.Time",
	"timestamp": "time.Time",
	"varchar":   "string",
	"char":      "string",
	"text":      "string",
	"blob":      "[]byte",
	"json":      "datatypes.JSON",
}

func GetAllTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query(showTablesSQL)
	if err != nil {
		return nil, errors.Wrapf(err, "query tables failed")
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return nil, errors.Wrapf(err, "scan result failed")
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func Generate(db *sql.DB, table string, matchInfo TableInfo) (StructLevel, error) {
	var structData = StructLevel{
		TableName: table,
	}
	structName := camelCase(matchInfo.GoStruct)
	camelStructName := string(unicode.ToUpper(rune(structName[0]))) + structName[1:]
	structData.StructName.UpperS = inflection.Singular(structName)
	structData.StructName.UpperP = inflection.Plural(structName)

	camelStructName = string(unicode.ToLower(rune(structName[0]))) + structName[1:]
	structData.StructName.LowerS = inflection.Singular(camelStructName)
	structData.StructName.LowerP = inflection.Plural(camelStructName)

	var createSQL string
	if db != nil {
		rows, err := db.Query(fmt.Sprintf(showCreateTableSQL, table))
		if err != nil {
			return structData, errors.Wrapf(err, "show table failed")
		}
		defer rows.Close()

		var t string
		for rows.Next() {
			err = rows.Scan(&t, &createSQL)
			if err != nil {
				return structData, err
			}
		}
	} else {
		return structData, errors.New("connect database failed")
	}
	structData.Columns = parseTable(createSQL)

	for _, v := range structData.Columns {
		switch v.GormName {
		case matchInfo.CreateTime:
			structData.TableCreateTime = matchInfo.CreateTime
			structData.FieldCreateTime = v.FieldName
			structData.IsTimestamp = matchInfo.IsTimestamp
		case matchInfo.UpdateTime:
			structData.TableUpdateTime = matchInfo.UpdateTime
			structData.FieldUpdateTime = v.FieldName
			structData.IsTimestamp = matchInfo.IsTimestamp
		case matchInfo.SoftDeleteKey:
			structData.TableSoftDeleteKey = matchInfo.SoftDeleteKey
			structData.TableSoftDeleteValue = matchInfo.SoftDeleteValue
			structData.FieldSoftDeleteKey = v.FieldName
		}
	}

	return structData, nil
}

func parseTable(s string) []FieldLevel {
	lines := strings.Split(s, "\n")

	isPrimaryKey := getPrimaryKey(lines)

	var columns []FieldLevel
	for _, line := range lines {
		line = strings.Trim(line, " ")
		if strings.HasPrefix(line, "`") {
			p := strings.Split(line, " ")
			name := strings.Trim(p[0], "`")
			dataType := p[1]
			isUnsigned := strings.Contains(line, unsignedKeyFlag)
			primaryKey := ""
			if isPrimaryKey[name] {
				primaryKey = ";primary_key"
			}
			columns = append(columns, FieldLevel{
				FieldName:  camelCase(name),
				FieldType:  fieldType(dataType, isUnsigned),
				PrimaryKey: primaryKey,
				GormName:   name,
				JsonName:   jsonCamelCase(name),
				Comment:    fieldComment(p),
			})
		}
	}
	return columns
}

func getPrimaryKey(lines []string) map[string]bool {
	priKey := make(map[string]bool)
	for _, line := range lines {
		line = strings.Trim(line, " ")
		if strings.HasPrefix(line, primaryKeyMark) {
			line = strings.Trim(line, primaryKeyMark)
			line = strings.TrimLeft(line, " ")
			p := strings.Split(line, " ")
			if len(p) > 0 {
				s := strings.Trim(p[0], "(")
				s = strings.Trim(s, ")")
				s = strings.ReplaceAll(s, "`", "")
				columns := strings.Split(s, ",")
				for _, name := range columns {
					priKey[name] = true
				}
				return priKey
			}
		}
	}
	return priKey
}

func fieldType(dataType string, isUnsigned bool) string {
	for tableType, goType := range tableToGoStruct {
		if strings.HasPrefix(dataType, tableType) {
			if isUnsigned && strings.HasPrefix(goType, "int") {
				return "u" + goType
			}
			return goType
		}
	}
	return "unknown"
}

func fieldComment(p []string) (comment string) {
	for i := len(p); i > 2; i-- {
		if strings.ToUpper(p[i-2]) == commentMark {
			comment = strings.Join(p[i-1:], " ")
			comment = strings.Trim(strings.Trim(comment, ","), "'")
			if comment != "" {
				break
			}
		}
	}
	return
}

func camelCase(s string) string {
	var buf bytes.Buffer
	var flag = false
	for i, c := range s {
		if c == '_' {
			flag = true
			continue
		}
		if i == 0 {
			buf.WriteRune(unicode.ToUpper(c))
			continue
		}
		if flag {
			buf.WriteRune(unicode.ToUpper(c))
			flag = false
			continue
		}
		buf.WriteRune(c)
	}
	return buf.String()
}

func jsonCamelCase(s string) string {
	var buf bytes.Buffer
	var flag = false
	for i, c := range s {
		if c == '_' {
			flag = true
			continue
		}
		if i == 0 {
			buf.WriteRune(unicode.ToLower(c))
			continue
		}
		if flag {
			buf.WriteRune(unicode.ToUpper(c))
			flag = false
			continue
		}
		buf.WriteRune(c)
	}
	return buf.String()
}
