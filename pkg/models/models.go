package models

import (
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"reflect"
	"regexp"
)

const component = "QueryBuilder"

func QueryBuilderGet(i interface{}, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	query := `SELECT `

	var searchByRow string
	var searchByIndex int
	for i := 0; i < v.NumField(); i++ {
		row := t.Field(i).Tag.Get("row")

		if !checkEmpty(v.Field(i)) {
			searchByRow = row
			searchByIndex = i
		}

		if row != "" {
			if i < t.NumField()-1 {
				query += row + ", "
			} else {
				query += row
			}
		}
	}

	query += " FROM " + tableName
	if searchByRow == "" {
		return query, nil
	}

	if t.Field(searchByIndex).Tag.Get("type") == "exact" {
		query += " WHERE " + searchByRow + " = ?"
	} else if t.Field(searchByIndex).Tag.Get("type") == "like" {
		query += " WHERE " + searchByRow + " LIKE ? COLLATE NOCASE"
	}
	args := []interface{}{v.Field(searchByIndex).Interface()}
	return query, args
}

func QueryBuilderCreate(i interface{}, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	query := `INSERT INTO ` + tableName + "("

	var valuesCount = 0
	args := make([]interface{}, 0)

	for i := 0; i < v.NumField(); i++ {
		row := t.Field(i).Tag.Get("row")

		if checkPK(t.Field(i)) {
			continue
		}

		if row != "" {
			if valuesCount != 0 {
				query += ", " + row
			} else {
				query += row
			}
			args = append(args, v.Field(i).Interface())
			valuesCount++
		}
	}

	query += ") values("
	for i := 0; i < valuesCount; i++ {
		if i < valuesCount-1 {
			query += "?, "
		} else {
			query += "?"
		}
	}

	query += ")"

	return query, args
}

func QueryBuilderDelete(i interface{}, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	query := `DELETE FROM ` + tableName + " WHERE "

	args := make([]interface{}, 0)

	for i := 0; i < v.NumField(); i++ {

		if !checkEmpty(v.Field(i)) {
			row := t.Field(i).Tag.Get("row")
			if row != "" {
				query += row + " = ?"
				args = append(args, v.Field(i).Interface())
				return query, args
			}
		}
	}
	return "", nil
}

func QueryBuilderUpdate(i interface{}, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	var searchBy int
	query := `UPDATE ` + tableName + " SET "
	args := make([]interface{}, 0)

	argsCount := 0
	for i := 0; i < v.NumField(); i++ {

		if checkPK(t.Field(i)) {
			searchBy = i
			continue
		}

		if !checkEmpty(v.Field(i)) {
			row := t.Field(i).Tag.Get("row")
			if row != "" {
				if argsCount < 1 {
					query += row + " = ?"
				} else {
					query += " ," + row + " = ?"
				}
				args = append(args, v.Field(i).Interface())
				argsCount++
			}
		}
	}

	if len(args) < 0 {
		return "", nil
	}

	query += " WHERE " + t.Field(searchBy).Tag.Get("row") + " = ?"
	args = append(args, v.Field(searchBy).Interface())

	return query, args
}

func checkEmpty(value reflect.Value) bool {
	// Checks int
	matchedInt, err := regexp.MatchString("int", value.Type().String())
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
	if matchedInt {
		return value.IsZero()
	}

	//else check string
	matchedString, err := regexp.MatchString("string", value.Type().String())
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
	if matchedString {
		return value.String() == ""
	}

	return !value.IsValid()
}

func checkPK(field reflect.StructField) bool {
	return field.Tag.Get("pk") != ""
}
