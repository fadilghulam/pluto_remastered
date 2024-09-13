package helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	db "pluto_remastered/config"

	"github.com/iancoleman/orderedmap"
)

func NewExecuteQuery(query string) ([]*orderedmap.OrderedMap, error) {

	queries := fmt.Sprintf(`SELECT JSON_AGG(data) as data FROM (%s) AS data`, query)

	rows, err := db.DB.Raw(queries).Rows()
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results, err := NewJsonDecode(rows, columns)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func NewJsonDecode(rows *sql.Rows, columns []string) ([]*orderedmap.OrderedMap, error) {
	// var result []map[string]interface{}
	var result []*orderedmap.OrderedMap
	var resultReturn []*orderedmap.OrderedMap
	for rows.Next() {
		// rowData := make(map[string]interface{})
		rowData := orderedmap.New()

		// Create a slice of interface{} to hold the values for Scan
		values := make([]interface{}, len(columns))
		for i := range columns {
			var value interface{}
			values[i] = &value
		}

		// Scan the row into the slice of interface{}
		if err := rows.Scan(values...); err != nil {
			log.Fatal(err)
		}

		// Transfer values from slice to map
		for i, col := range columns {
			// rowData[col] = *values[i].(*interface{})
			rowData.Set(col, *values[i].(*interface{}))
		}

		// Append the map to the result slice
		result = append(result, rowData)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, m := range result {
		getData, _ := m.Get("data")
		if getData != nil {
			stringData := string(getData.([]byte))

			var datas []*orderedmap.OrderedMap

			err := json.Unmarshal([]byte(stringData), &datas)
			if err != nil {
				log.Fatal(err)
			}

			resultReturn = append(resultReturn, datas...)
		}

	}

	return resultReturn, nil
}
