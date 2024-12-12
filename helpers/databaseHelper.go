package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	db "pluto_remastered/config"
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

func ExecuteGORMQueryWithoutResult(query string, wg *sync.WaitGroup) {
	defer wg.Done()

	db.DB.Exec(query)
}

func ExecuteGORMQueryIndexString(query string, resultsChan chan<- map[string][]map[string]interface{}, index string, wg *sync.WaitGroup) {
	defer wg.Done()

	results, _ := ExecuteQuery(query)

	var res []map[string]interface{}

	queries := fmt.Sprintf(`SELECT JSON_AGG(data) as data FROM (%s) AS data`, query)

	if err := db.DB.Exec(queries).Scan(&res).Error; err != nil {
		fmt.Println("a")
	}

	fmt.Println(db.DB.Exec(queries))

	for _, body := range results {

		for key, value := range body {
			if key == "id" || key == "customer_id" || key == "penjualan_id" ||
				key == "piutang_id" || key == "pengembalian_id" || key == "kunjungan_id" ||
				key == "pembayaran_piutang_id" || key == "payment_id" {
				// if strings.Contains(key, "id") && (key != "user_id" || key != "user_id_subtitute")  {
				switch v := value.(type) {
				case json.Number:
					// Convert float64 to an integer, then to a string
					body[key] = v.String()
				default:
					// Convert other types to a string
					body[key] = fmt.Sprintf("%v", value)
				}
			}

		}
	}

	resultsChan <- map[string][]map[string]interface{}{index: results}
}
