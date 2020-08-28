package models

import (
	"database/sql"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/database"
	"github.com/Jizzberry/Jizzberry_go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"reflect"
	"regexp"
	"strings"
)

func QueryBuilderGet(i interface{}, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	query := `SELECT `

	// Name of column to search by
	var searchByColumn string

	// Index of value of search-by column in interface provided
	var searchByIndex int

	queryTags := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		column := t.Field(i).Tag.Get(helpers.RowStructTag)

		// Last non-empty Field is search-by
		// Normally only one Field should be non-empty
		if !checkEmpty(v.Field(i)) {
			searchByColumn = column
			searchByIndex = i
		}

		// Add to query only if row is not empty
		// Useful when an extra - non column field is present in struct
		if column != "" {
			queryTags = append(queryTags, column)
		}
	}

	query += strings.Join(queryTags, ", ") + " FROM " + tableName
	if searchByColumn == "" {
		return query, nil
	}

	if t.Field(searchByIndex).Tag.Get("type") == "exact" {
		query += " WHERE " + searchByColumn + " = ?"
	} else if t.Field(searchByIndex).Tag.Get("type") == "like" {
		query += " WHERE " + searchByColumn + " LIKE ? COLLATE NOCASE"
	}
	args := []interface{}{v.Field(searchByIndex).Interface()}

	return query, args
}

func QueryBuilderMatch(i interface{}, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	query := `SELECT `

	// Name of column to search by
	var searchByColumn string

	// Index of value of search-by column in interface provided
	var searchByIndex int

	queryTags := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		column := t.Field(i).Tag.Get(helpers.RowStructTag)

		if !checkEmpty(v.Field(i)) {
			searchByColumn = column
			searchByIndex = i
		}

		if column != "" {
			queryTags = append(queryTags, column)
		}
	}

	query += strings.Join(queryTags, ", ") + " FROM " + tableName
	if searchByColumn == "" {
		return query, nil
	}
	query += " WHERE (" + searchByColumn + " LIKE ? COLLATE NOCASE) OR (replace(" + searchByColumn + ", ' ', '') LIKE ? COLLATE NOCASE)"
	args := []interface{}{"%" + v.Field(searchByIndex).String() + "%", v.Field(searchByIndex).Interface()}

	return query, args
}

func QueryBuilderCreate(i interface{}, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	query := `INSERT INTO ` + tableName + "("

	// Number of values to enter into DB
	var valuesCount = 0

	args := make([]interface{}, 0)

	for i := 0; i < v.NumField(); i++ {
		row := t.Field(i).Tag.Get(helpers.RowStructTag)

		if isPK(t.Field(i)) {
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
			row := t.Field(i).Tag.Get(helpers.RowStructTag)
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

	// Key at which data should be updated
	var searchBy int

	query := `UPDATE ` + tableName + " SET "
	args := make([]interface{}, 0)

	argsCount := 0
	for i := 0; i < v.NumField(); i++ {

		if isPK(t.Field(i)) {
			searchBy = i
			continue
		}

		row := t.Field(i).Tag.Get(helpers.RowStructTag)
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

	if len(args) < 0 {
		return "", nil
	}

	query += " WHERE " + t.Field(searchBy).Tag.Get(helpers.RowStructTag) + " = ?"
	args = append(args, v.Field(searchBy).Interface())

	return query, args
}

// Creates an array of provided interface,
// Scans values from row into that array
func GetIntoStruct(rows *sql.Rows, dest interface{}) {
	v := reflect.ValueOf(dest)
	direct := reflect.Indirect(v)

	if v.Kind() != reflect.Ptr {
		helpers.LogError("Destination not pointer")
		return
	}

	if direct.Kind() != reflect.Slice {
		helpers.LogError("Destination not slice")
		return
	}

	base := v.Elem().Type().Elem()
	vp := reflect.New(base)

	for rows.Next() {
		direct.Set(reflect.Append(direct, scanSingleStruct(vp, rows)))
	}
}

// Scans values from row into struct
func scanSingleStruct(dest reflect.Value, row *sql.Rows) reflect.Value {
	numfields := reflect.Indirect(dest).NumField()
	ind := reflect.Indirect(dest)

	ptrs := make([]interface{}, 0)

	for i := 0; i < numfields; i++ {
		ptrs = append(ptrs, ind.Field(i).Addr().Interface())
	}

	err := row.Scan(ptrs...)
	if err != nil {
		helpers.LogError(err.Error())
	}
	return ind
}

func IsTableEmpty(tableName string, conn *sql.DB) {
	rows, err := conn.Query(`SELECT count(name) FROM sqlite_master WHERE type='table' and name=?`, tableName)

	if err != nil {
		helpers.LogError(err.Error())
		err := database.RunMigrations()
		if err != nil {
			helpers.LogError(err.Error())
		}
		return
	}
	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			helpers.LogError(err.Error())
		}
	}

	if count < 0 {
		err := database.RunMigrations()
		if err != nil {
			helpers.LogError(err.Error())
		}
	}
}

func IsValueExists(conn *sql.DB, key interface{}, keyname string, tableName string) (bool, int64) {
	rows, err := conn.Query(fmt.Sprintf(`SELECT * FROM %s WHERE  %s=?`, tableName, keyname), key)

	if err != nil {
		helpers.LogError(err.Error())
		return false, -1
	}

	var genId int64 = -1
	for rows.Next() {
		err := rows.Scan(&genId)
		if err != nil {
			helpers.LogError(err.Error())
		}
	}

	if genId > -1 {
		return true, genId
	}

	return false, -1
}

func checkEmpty(value reflect.Value) bool {
	// Checks int
	matchedInt, err := regexp.MatchString("int", value.Type().String())
	if err != nil {
		helpers.LogError(err.Error())
		return false
	}
	if matchedInt {
		return value.IsZero()
	}

	matchedFloat, err := regexp.MatchString("float", value.Type().String())
	if err != nil {
		helpers.LogError(err.Error())
		return false
	}
	if matchedFloat {
		return value.IsZero()
	}

	//else check string
	matchedString, err := regexp.MatchString("string", value.Type().String())
	if err != nil {
		helpers.LogError(err.Error())
		return false
	}
	if matchedString {
		return value.String() == ""
	}

	//else check bool
	matchedBool, err := regexp.MatchString("bool", value.Type().String())
	if err != nil {
		helpers.LogError(err.Error())
		return false
	}
	if matchedBool {
		// Bool cant be search factor
		return true
	}

	return !value.IsValid()
}

func isPK(field reflect.StructField) bool {
	return field.Tag.Get(helpers.PKStructTag) == "auto"
}

// Runs migrations if needed before returning conn
func GetConn(table string) *sql.DB {
	conn := database.GetConn(router.GetDatabase(table))
	IsTableEmpty(table, conn)
	return conn
}
