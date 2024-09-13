package helpers

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"text/template"

	newOrderedmap "github.com/iancoleman/orderedmap"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func PrepareQuery(query string, args map[string]interface{}) (string, error) {
	tmpl, err := template.New("sqlQuery").Parse(query)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return "", err
	}

	var queryBuffer bytes.Buffer
	err = tmpl.Execute(&queryBuffer, args)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return "", err
	}

	finalQuery := strings.Replace(queryBuffer.String(), "<no value>", "", -1)
	return finalQuery, nil
}

func ExecuteGORMQuery(query string, resultsChan chan<- map[int][]map[string]interface{}, index int, wg *sync.WaitGroup) {
	defer wg.Done()

	results, _ := ExecuteQuery(query)

	resultsChan <- map[int][]map[string]interface{}{index: results}
}
func ExecuteGORMQuery2(query string, resultsChan chan<- map[int][]*orderedmap.OrderedMap[string, interface{}], index int, wg *sync.WaitGroup, specialCondition string) {
	defer wg.Done()

	results, _ := ExecuteQuery2(query, specialCondition)

	resultsChan <- map[int][]*orderedmap.OrderedMap[string, interface{}]{index: results}
}

func ExecuteGORMQueryOrdered(query string, resultsChan chan<- map[int][]*newOrderedmap.OrderedMap, index int, wg *sync.WaitGroup) {
	defer wg.Done()

	results, _ := NewExecuteQuery(query)

	resultsChan <- map[int][]*newOrderedmap.OrderedMap{index: results}
}
