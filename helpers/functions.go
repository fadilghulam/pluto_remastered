package helpers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	db "pluto_remastered/config"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func ParamArray(param [][]byte) []string {
	var arr []string

	for _, idStr := range param {
		id := string(idStr)
		arr = append(arr, id)
	}

	return arr
}

func ArraySum(nums []float64) float64 {
	var sum float64
	for _, num := range nums {
		sum += num
	}
	return sum
}

func ArrayValues(data []map[string]interface{}) []interface{} {
	var values []interface{}
	for _, m := range data {
		for _, v := range m {
			values = append(values, v)
		}
	}
	return values
}

func ArrayKeys(data []map[string]interface{}) []string {
	var values []string
	for _, m := range data {
		for k := range m {
			values = append(values, k)
		}
	}
	return values
}

func parseMonthYear(s string) (time.Month, int) {
	layout := "Jan-2006"
	t, err := time.Parse(layout, s)
	if err != nil {
		return time.Month(0), 0
	}
	return t.Month(), t.Year()
}

func ParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func ParseDate(s string) time.Time {
	layout := "2006-01-02"
	t, err := time.Parse(layout, s)
	if err != nil {
		return time.Now()
	}
	return t
}

// Sort a slice of month-year strings in chronological order
func SortMonthYear(monthYears []string) []string {
	sort.SliceStable(monthYears, func(i, j int) bool {
		monthI, yearI := parseMonthYear(monthYears[i])
		monthJ, yearJ := parseMonthYear(monthYears[j])
		if yearI != yearJ {
			return yearI < yearJ
		}
		return monthI < monthJ
	})
	return monthYears
}

func GetSortedValues(m map[string]interface{}) (keys []string, values []float64) {
	// Extract keys
	for k := range m {
		keys = append(keys, k)
	}

	// Sort keys
	sort.Strings(keys)

	// Extract values in the sorted key order
	for _, k := range keys {
		values = append(values, m[k].(float64))
	}

	return
}

func ReorderData(data []map[string]interface{}) []map[string]interface{} {
	var orderedData []map[string]interface{}

	// Process each map in the slice
	for _, item := range data {
		// Create a new map to hold ordered data
		orderedMap := make(map[string]interface{})

		// Extract product and month keys
		monthKeys := []string{}
		var product string

		for key := range item {
			if key == "product" {
				product = item[key].(string)
			} else {
				monthKeys = append(monthKeys, key)
			}
		}

		// Sort the month keys chronologically
		sort.Slice(monthKeys, func(i, j int) bool {
			layout := "Jan-2006"
			t1, _ := time.Parse(layout, monthKeys[i])
			t2, _ := time.Parse(layout, monthKeys[j])
			return t1.Before(t2)
		})

		// First add the product key to the ordered map
		orderedMap["product"] = product

		// Add the sorted months to the ordered map
		for _, month := range monthKeys {
			orderedMap[month] = item[month]
		}

		// Append the ordered map to the result slice
		orderedData = append(orderedData, orderedMap)
	}

	// fmt.Println(orderedData)

	return orderedData
}

func ExecuteQuery(query string) ([]map[string]interface{}, error) {

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

	// datas, err := SaveRowToDynamicStruct(rows, columns)
	// if err != nil {
	// 	return nil, err
	// }

	results, err := JsonDecode(rows, columns)
	if err != nil {
		return nil, err
	}

	if results[0]["data"] == nil {
		return nil, nil
	}

	return results[0]["data"].([]map[string]interface{}), nil
}

func ExecuteQuery2(query string, specialCondition string) ([]*orderedmap.OrderedMap[string, interface{}], error) {

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

	results, err := JsonDecode2(rows, columns, specialCondition)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func InsertDataDynamic(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {

	tx := db.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	transactionData := data["transactions"].(map[string]interface{})

	dataInsert := make(map[string]interface{})
	dataDelete := make(map[string]interface{})

	insertedIds := make(map[string][]interface{})

	for tableName, data := range transactionData {

		if IsDeletedIds(tableName) {
			pattern := regexp.MustCompile(`_deleted_ids`)

			// Perform the replacement
			tableName = pattern.ReplaceAllString(tableName, "")

			dataDelete[tableName] = data
		} else {
			dataSlice := data.([]interface{})

			for i := 0; i < len(dataSlice); i++ {
				for key, value := range dataSlice[i].(map[string]interface{}) {
					if b, ok := value.(bool); ok {
						if b {
							dataSlice[i].(map[string]interface{})[key] = 1
						} else {
							dataSlice[i].(map[string]interface{})[key] = 0
						}
					}
				}
			}

			dataInsert[tableName] = dataSlice
		}
	}

	for tablenames, data := range dataInsert {

		tempStruct, err := CreateStructInstance(tablenames)
		if err != nil {
			return nil, err
		}

		tempTable := tablenames
		tempSchema := ""
		whereSchema := ""
		tempTableName := strings.Split(tablenames, ".")
		if len(tempTableName) > 1 {
			tempTable = tempTableName[1]
			tempSchema = tempTableName[0]
			whereSchema = fmt.Sprintf(" AND table_schema = '%s'", tempSchema)
		}

		query := fmt.Sprintf(`SELECT 1
								FROM information_schema.columns
								WHERE table_name = '%s' AND column_name = 'dtm_upd' %s
								ORDER BY ordinal_position`, tempTable, whereSchema)

		var count int64

		if err := db.DB.Raw(query).Count(&count).Error; err != nil {
			panic(err.Error())
		}

		var columnTime string

		// Check if the query returned any rows
		if count == 0 {
			columnTime = "updated_at"
		} else {
			columnTime = "dtm_upd"
		}

		// fmt.Println(columnTime)
		param := fmt.Sprintf("id , (%s - '5 day'::interval) as %s ", columnTime, columnTime)
		dataSlice := data.([]interface{})

		// fmt.Println(dataSlice)
		for i := 0; i < len(dataSlice); i++ {

			element := dataSlice[i].(map[string]interface{})
			// fmt.Println(element["id"])
			temp, err := db.DB.Raw(fmt.Sprintf("SELECT %s FROM %s WHERE id = %v", param, tablenames, element["id"])).Rows()
			if err != nil {
				panic(err.Error())
			}

			columns, err := temp.Columns()
			if err != nil {
				panic(err.Error())
			}

			defer temp.Close()
			checkCol, err := SaveRowToDynamicStruct(temp, columns)
			if err != nil {
				panic(err.Error())
			}

			if len(checkCol) > 0 {
				for _, m := range checkCol {

					var tempRow map[string]interface{}
					id := m["id"]
					coltime := m[columnTime]

					var timeVal time.Time

					if dataSlice[i].(map[string]interface{})[columnTime] == nil {
						dataSlice[i].(map[string]interface{})[columnTime] = time.Now()
						timeVal = time.Now()
					} else {
						timeStr := dataSlice[i].(map[string]interface{})[columnTime].(string)
						if dataSlice[i].(map[string]interface{})[columnTime] == nil {
							dataSlice[i].(map[string]interface{})[columnTime] = time.Now()
						}

						timeVal, err = time.Parse("2006-01-02 15:04:05", timeStr)
						if err != nil {
							fmt.Println("Error parsing time:", err)
							continue
						}
					}

					if coltime != nil {

						// fmt.Println(timeVal, coltime.(time.Time))

						if timeVal.After(coltime.(time.Time)) {
							tempRow = dataSlice[i].(map[string]interface{})
						}
					} else {
						dataSlice[i].(map[string]interface{})[columnTime] = time.Now()
						tempRow = dataSlice[i].(map[string]interface{})
					}
					if len(tempRow) > 0 {

						idData := dataSlice[i].(map[string]interface{})["id"]
						insertedIds[tablenames] = append(insertedIds[tablenames], idData)

						structValuePtr := reflect.ValueOf(tempStruct)
						if structValuePtr.Kind() != reflect.Ptr || structValuePtr.Elem().Kind() != reflect.Struct {
							return nil, fmt.Errorf("tempStruct is not a pointer to a struct")
						}

						// Dereference the pointer to get the struct value
						structValue := structValuePtr.Elem()

						dataMap := dataSlice[i].(map[string]interface{})
						err := mapstructure.Decode(dataMap, structValue.Addr().Interface())
						if err != nil {
							return nil, err
						}

						// Validate the tempStruct
						err = ValidateStructInstance(tempStruct)
						if err != nil {
							return nil, err
						}

						delete(dataSlice[i].(map[string]interface{}), "id")
						db.DB.Model(tempStruct).Where("id = ?", id).Updates(dataSlice[i].(map[string]interface{}))
					}
				}
			} else {

				structValuePtr := reflect.ValueOf(tempStruct)
				if structValuePtr.Kind() != reflect.Ptr || structValuePtr.Elem().Kind() != reflect.Struct {
					return nil, fmt.Errorf("tempStruct is not a pointer to a struct")
				}

				// Dereference the pointer to get the struct value
				structValue := structValuePtr.Elem()

				dataMap := dataSlice[i].(map[string]interface{})
				err := mapstructure.Decode(dataMap, structValue.Addr().Interface())
				if err != nil {
					return nil, err
				}

				// Validate the tempStruct
				err = ValidateStructInstance(tempStruct)
				if err != nil {
					return nil, err
				}

				db.DB.Model(tempStruct).Create(dataSlice[i].(map[string]interface{}))
				idData := dataSlice[i].(map[string]interface{})["id"]
				insertedIds[tablenames] = append(insertedIds[tablenames], idData)
			}
		}
	}

	returnedData := make(map[string]interface{})
	for tablenames, ids := range insertedIds {
		tempIds := Implode(ids)
		rows, err := db.DB.Raw("SELECT JSON_AGG(data) as data FROM (SELECT * FROM "+tablenames+" WHERE id IN (?) ) data", tempIds).Rows()
		if err != nil {
			return nil, err
		}

		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		datas, err := JsonDecode(rows, columns)
		if err != nil {
			return nil, err
		}
		returnedData[tablenames] = datas
	}

	for tempname, value := range returnedData {
		for _, v := range value.([]map[string]interface{}) {
			for _, vdata := range v {
				returnedData[tempname] = vdata
			}
		}
	}

	for tablenames, data := range dataDelete {
		tempStruct, err := CreateStructInstance(tablenames)
		if err != nil {
			return nil, err
		}
		db.DB.Where("id IN (?)", Implode(data.([]interface{}))).Delete(tempStruct)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return returnedData, nil
}

func JsonDecode(rows *sql.Rows, columns []string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	for rows.Next() {
		rowData := make(map[string]interface{})

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
			rowData[col] = *values[i].(*interface{})
		}

		// Append the map to the result slice
		result = append(result, rowData)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	for i, m := range result {
		for key, value := range m {
			if bytes, isBytes := value.([]byte); isBytes {
				// fmt.Println(isBytes)
				var decodedData []map[string]interface{}
				if err := json.Unmarshal(bytes, &decodedData); err != nil {
					log.Fatal(err)
				} else {
					result[i][key] = decodedData
					// fmt.Println(decodedData)
				}
			}
		}
	}

	return result, nil
}

func ConvertMapToOrderedMap(input map[string]interface{}) *orderedmap.OrderedMap[string, interface{}] {
	om := orderedmap.New[string, interface{}]()
	for k, v := range input {
		switch val := v.(type) {
		case map[string]interface{}:
			// Recursively convert nested maps
			om.Set(k, ConvertMapToOrderedMap(val))
		default:
			om.Set(k, v)
		}
	}
	return om
}

func JsonDecode2(rows *sql.Rows, columns []string, specialCondition string) ([]*orderedmap.OrderedMap[string, interface{}], error) {
	// var result []map[string]interface{}
	var result []*orderedmap.OrderedMap[string, interface{}]
	var resultReturn []*orderedmap.OrderedMap[string, interface{}]
	for rows.Next() {
		// rowData := make(map[string]interface{})
		rowData := orderedmap.New[string, interface{}]()

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

		stringData := string(m.Value("data").([]byte))

		if stringData[0] == '[' {
			stringData = stringData[1:]
		}
		if last := len(stringData) - 1; last >= 0 && stringData[last] == ']' {
			stringData = stringData[:last]
		}
		// fmt.Println(stringData)

		// var rawData []map[string]interface{}
		// err := json.Unmarshal([]byte(stringData), &rawData)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// // Convert each map into an ordered map and recursively handle nested maps
		// for _, item := range rawData {
		// 	orderedRow := ConvertMapToOrderedMap(item)
		// 	resultReturn = append(resultReturn, orderedRow)
		// }

		// fmt.Println("=====================================================")

		if specialCondition == "" {
			specialCondition = "},"
		}

		mapStringData := strings.Split(stringData, specialCondition)
		// fmt.Println(mapStringData)
		for _, v := range mapStringData {
			om := orderedmap.New[string, interface{}]()

			_ = om.UnmarshalJSON([]byte(v + "}"))
			resultReturn = append(resultReturn, om)
		}

	}

	return resultReturn, nil
}

func JsonDecodeMap(input map[string]interface{}) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	// Create a map to hold the data
	rowData := make(map[string]interface{})

	// Iterate over the input map
	for key, value := range input {
		if bytes, isBytes := value.([]byte); isBytes {
			// Unmarshal bytes into a map
			var decodedData []map[string]interface{}
			if err := json.Unmarshal(bytes, &decodedData); err != nil {
				return nil, err
			}
			rowData[key] = decodedData
		} else {
			rowData[key] = value
		}
	}

	// Append the map to the result slice
	result = append(result, rowData)

	return result, nil
}

func JoinStrings(strings []string, separator string) string {
	result := ""
	for i, s := range strings {
		if i > 0 {
			result += separator
		}
		result += s
	}
	return result
}

func SplitToString(a []int, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}

func Implode(interfaceSlice []interface{}) []int64 {
	intSlice := make([]int64, len(interfaceSlice))
	for i, v := range interfaceSlice {

		if val, ok := v.(float64); ok {
			intSlice[i] = int64(val)
		} else if val, ok := v.(string); ok {
			temp, _ := strconv.Atoi(val)
			intSlice[i] = int64(temp)
		} else if val, ok := v.(int64); ok {
			intSlice[i] = int64(val)
		} else {
			fmt.Println("unknown type", v, reflect.TypeOf(v))
		}
	}
	return intSlice
}

func IsDeletedIds(s string) bool {
	pattern := `.*_deleted_ids`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

func SaveRowToDynamicStruct(rows *sql.Rows, columns []string) ([]map[string]interface{}, error) {

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	var results []map[string]interface{}

	// Iterate through the rows and store in the slice
	for rows.Next() {
		// Scan the values into the value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Create a map for the row
		rowMap := make(map[string]interface{})

		// Fill the map with column name and corresponding value
		for i, col := range columns {
			val := values[i]
			rowMap[col] = val
		}

		// Append the row map to the results slice
		results = append(results, rowMap)
	}

	return results, nil
}

func ByteResponse(responseBody []byte) (map[string]interface{}, error) {
	var responses map[string]interface{}
	if err := json.Unmarshal(responseBody, &responses); err != nil {
		// fmt.Println("Error:", err)
		return nil, err
	}

	responseData := make(map[string]interface{})
	for key, val := range responses {
		responseData[key] = val
	}

	return responseData, nil
}

func PostBody(body []byte) (map[string]interface{}, error) {
	// bodyBytes := c.Body()

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return data, nil
}

func ConvertStringToInt64(i string) int64 {

	temp, _ := strconv.Atoi(i)
	temp2 := int64(temp)
	return temp2
}
func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func TrimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

func ArrayColumn(data []any, column string) []any {

	var result []any
	for _, v := range data {
		result = append(result, v.(map[string]interface{})[column])
	}

	return result
}

// func SendNotification(title string, body string, userIds int, customerId int64, dataSend map[string]string, c *fiber.Ctx) error {

// 	tokenFCM := new([]model.TokenFcm)

// 	// fmt.Println(userIds)

// 	err := db.DB.Where("user_id IN ? AND app_name = ? AND customer_id = ?", userIds, "tokoku", customerId).Find(&tokenFCM).Error

// 	if err != nil {
// 		log.Println(err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(ResponseWithoutData{
// 			Message: "Something went wrong",
// 			Success: false,
// 		})
// 	}

// 	var tokens []string
// 	for _, token := range *tokenFCM {
// 		tokens = append(tokens, token.Token)
// 	}

// 	type NotificationRequest struct {
// 		Title  string            `json:"title"`
// 		Body   string            `json:"body"`
// 		Tokens []string          `json:"tokens"`
// 		Data   map[string]string `json:"data"` // additional data
// 	}

// 	var req NotificationRequest

// 	req.Title = title
// 	req.Body = body
// 	req.Tokens = tokens
// 	req.Data = dataSend

// 	// Load the service account key JSON file
// 	opt := option.WithCredentialsFile("middleware/tokoku.json")

// 	// Initialize the Firebase app
// 	ctx := context.Background()
// 	app, err := firebase.NewApp(ctx, nil, opt)
// 	if err != nil {
// 		log.Fatalf("error initializing app: %v\n", err)
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to initialize Firebase app",
// 		})
// 	}

// 	// Initialize the FCM client
// 	client, err := app.Messaging(ctx)
// 	if err != nil {
// 		log.Fatalf("error getting Messaging client: %v\n", err)
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to initialize FCM client",
// 		})
// 	}

// 	// Create the multicast message to send
// 	message := &messaging.MulticastMessage{
// 		Tokens: req.Tokens,
// 		Notification: &messaging.Notification{
// 			Title: req.Title,
// 			Body:  req.Body,
// 		},
// 		Data: req.Data,
// 	}

// 	// Send the message
// 	response, err := client.SendMulticast(ctx, message)
// 	if err != nil {
// 		log.Fatalf("error sending FCM message: %v\n", err)
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to send FCM notification",
// 		})
// 	}

// 	log.Printf("Successfully sent FCM message: %v\n", response)
// 	return c.JSON(fiber.Map{
// 		"message": "Notification sent successfully",
// 		"success": response.SuccessCount,
// 		"failure": response.FailureCount,
// 		"errors":  response.Responses,
// 	})
// }

func NewCurl(data map[string]string, method string, url string, c *fiber.Ctx) map[string]interface{} {

	client := &http.Client{}

	dataSend, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	// Create a POST request with a JSON payload
	req, err := http.NewRequest(method, url, bytes.NewReader(dataSend))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	responseData, err := ByteResponse(responseBody)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	return responseData
}

func RefreshUser(userId string) ([]map[string]interface{}, error) {

	datas, err := ExecuteQuery(fmt.Sprintf(`WITH data_customer AS (SELECT c.user_id, JSONB_AGG(c.*) as datas FROM tk.get_customer_by_userid(%v) c GROUP BY c.user_id)

											SELECT NULL as employee,
												u.id,
												u.full_name as name,
												u.username,
												u.profile_photo,
												p.email,
												p.phone,
												p.ktp,
												ARRAY[]::varchar[] as permission,
												dc.datas as "userInfo"
											FROM public.user u
											LEFT JOIN data_customer dc
												ON u.id = dc.user_id
											LEFT JOIN hr.person p
												ON u.id = p.user_id
											WHERE u.id = %v
											GROUP BY u.id, p.id, dc.datas`, userId, userId))

	if err != nil {
		return nil, err
	}

	return datas, nil
}

func SendCurl(data []byte, method string, url string) (map[string]interface{}, error) {

	client := &http.Client{}

	// req, err := http.NewRequest("GET", "https://rest.pt-bks.com/pluto-mobile/md/getDataTodayMD2", bytes.NewReader(dataSend))
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	responseData, err := ByteResponse(responseBody)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	if responseData["data"] == nil {
		return nil, fiber.ErrNotFound
	}

	switch responseData["data"].(type) {
	case map[string]interface{}:
		if len(responseData["data"].(map[string]interface{})) == 0 {
			responseData["data"] = nil
		}
	case []interface{}:
		if len(responseData["data"].([]interface{})) == 0 {
			responseData["data"] = nil
		}
	case string:
		responseData["data"] = responseData["data"].(string)
	default:
		responseData["data"] = nil

	}

	return responseData, nil
}

func RemoveZeroFromMap(data map[string]interface{}) (map[string]interface{}, []string, []string) {
	var remainingKeys, removedKeys []string
	for key, val := range data {
		if val.(int) < 1 {
			removedKeys = append(removedKeys, key)
			delete(data, key)
		} else {
			remainingKeys = append(remainingKeys, key)
		}
	}

	return data, remainingKeys, removedKeys
}

func ItemExists(slice []interface{}, value int) bool {
	for _, v := range slice {
		switch v := v.(type) {
		case int:
			if v == value {
				return true
			}
		case float64:
			if int(v) == value { // Cast float64 to int for comparison
				return true
			}
		}
	}
	return false
}
