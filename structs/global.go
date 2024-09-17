package structs

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

func GetStructInstanceByTableName(tableName string) (interface{}, error) {
	// Map table names to struct types
	tableStructMap := map[string]reflect.Type{
		TableNameCustomerHistory: reflect.TypeOf(CustomerHistory{}),
		TableNameCustomerKtp:     reflect.TypeOf(CustomerKtp{}),
		TableNameCustomerTokoku:  reflect.TypeOf(CustomerTokoku{}),
	}

	if structType, exists := tableStructMap[tableName]; exists {
		// Create a new instance of the struct and return it
		return reflect.New(structType).Interface(), nil
	}
	return nil, fmt.Errorf("no struct found for table name: %s", tableName)
}

type Int32Array []int32

// Value converts Int32Array to a PostgreSQL array-compatible format.
func (a Int32Array) Value() (driver.Value, error) {
	// Convert []int32 to []interface{}
	var arr = make([]interface{}, len(a))
	for i, v := range a {
		arr[i] = v
	}
	return arr, nil
}

func (a Int32Array) Value2() (driver.Value, error) {
	// Convert to []int32 to []int64 for pq.Array
	int64Array := make([]int64, len(a))
	for i, v := range a {
		int64Array[i] = int64(v)
	}
	return pq.Array(int64Array), nil
}

// Scan converts a PostgreSQL array to Int32Array.
func (a *Int32Array) Scan(value interface{}) error {
	var ints []int32

	switch v := value.(type) {
	case string:
		// Handle the case where the array is returned as a string
		trimmed := strings.Trim(v, "{}")
		if len(trimmed) == 0 {
			*a = []int32{}
			return nil
		}
		strElements := strings.Split(trimmed, ",")
		for _, strElem := range strElements {
			i, err := strconv.Atoi(strElem)
			if err != nil {
				return err
			}
			ints = append(ints, int32(i))
		}
	case []byte:
		// Handle the case where the array is returned as []byte
		trimmed := strings.Trim(string(v), "{}")
		if len(trimmed) == 0 {
			*a = []int32{}
			return nil
		}
		strElements := strings.Split(trimmed, ",")
		for _, strElem := range strElements {
			i, err := strconv.Atoi(strElem)
			if err != nil {
				return err
			}
			ints = append(ints, int32(i))
		}
	default:
		return fmt.Errorf("unsupported data type: %T", v)
	}

	*a = ints
	return nil
}
