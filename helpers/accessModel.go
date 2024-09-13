package helpers

import (
	"fmt"
	"reflect"

	// "go_sales_api/internal/model"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var structMap = map[string]reflect.Type{
	// "Cashiers":   reflect.TypeOf(model.Cashiers{}),
	// "Categories": reflect.TypeOf(model.Categories{}),
	// Add more struct types as needed
}

func GetStructMap(tablenames string) (reflect.Type, bool) {
	val, ok := structMap[tablenames]
	return val, ok
}

// func CreateStructInstance(structName string) (interface{}, error) {
// 	// Get the type of the struct based on its name
// 	structType, err := StructNameToType(structName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create a new instance of the struct
// 	instance := reflect.New(structType).Elem().Interface()
// 	return instance, nil
// }

func CreateStructInstance(structName string) (interface{}, error) {
	// Get the type of the struct based on its name
	structType, err := StructNameToType(structName)
	if err != nil {
		return nil, err
	}

	// Create a new instance of the struct as a pointer
	instancePtr := reflect.New(structType).Interface()
	return instancePtr, nil
}

func StructNameToType(structName string) (reflect.Type, error) {
	structName = cases.Title(language.English, cases.NoLower).String(structName)

	// Get the type from structMap
	structType, ok := structMap[structName]
	if !ok {
		return nil, fmt.Errorf("struct type %s not found", structName)
	}
	return structType, nil
}

func ValidateStructInstance(instance interface{}) error {
	// Get the value of the instance
	value := reflect.ValueOf(instance)

	// Check if the instance has a Validate method
	validateMethod := value.MethodByName("Validate")
	if !validateMethod.IsValid() {
		// If Validate method doesn't exist, return nil error
		return nil
	}

	// Call the Validate method and get the result
	result := validateMethod.Call(nil)
	if len(result) > 0 {
		// Check if the result has an error
		errValue := result[0]
		if !errValue.IsNil() {
			// If there's an error, return it
			err := errValue.Interface().(error)
			return err
		}
	}

	// If everything is successful, return nil error
	return nil
}
