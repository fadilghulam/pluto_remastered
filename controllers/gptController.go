package controllers

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	db "pluto_remastered/config"
	structBot "pluto_remastered/structs/gpt_structs"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

func FetchRandomRow[T any](db *gorm.DB, tableName string, result *T) error {

	// Count rows in the table
	var count int64
	if err := db.Table(tableName).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("no rows found in table %s", tableName)
	}

	// Get a random offset
	offset := rand.Intn(int(count))

	// Fetch a single row at the random offset
	if err := db.Table(tableName).Offset(offset).Limit(1).Find(result).Error; err != nil {
		return err
	}

	return nil
}

func FetchRandomItems(db *gorm.DB, maxItems int) ([]structBot.Item, error) {

	// Count the total number of items
	var count int64
	if err := db.Model(&structBot.Item{}).Count(&count).Error; err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("no items found in the database")
	}

	// Generate a random number of items to fetch (up to maxItems)
	numItems := rand.Intn(maxItems) + 1

	// Fetch random items with their item types
	var items []structBot.Item
	if err := db.
		Order("RANDOM()").
		Limit(numItems).
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func FetchLinkedDataSr(db *gorm.DB) (*structBot.Sr, *structBot.Rayon, *structBot.Branch, error) {
	// Step 1: Fetch a random SR
	var sr structBot.Sr
	if err := FetchRandomRow(db, "sr", &sr); err != nil {
		return nil, nil, nil, err
	}

	// Step 2: Fetch a random Rayon linked to the selected SR
	var rayon structBot.Rayon
	if err := db.Where("sr_id = ?", sr.ID).Limit(1).Find(&rayon).Error; err != nil {
		return &sr, nil, nil, err
	}

	// Step 3: Fetch a random Branch linked to the selected Rayon
	var branch structBot.Branch
	if err := db.Where("rayon_id = ?", rayon.ID).Limit(1).Find(&branch).Error; err != nil {
		return &sr, &rayon, nil, err
	}

	return &sr, &rayon, &branch, nil
}

func FetchLinkedDataProvince(db *gorm.DB) (*structBot.Province, *structBot.Regency, *structBot.District, *structBot.SubDistrict, error) {
	// Step 1: Fetch a random SR
	var province structBot.Province
	if err := FetchRandomRow(db, "province", &province); err != nil {
		return nil, nil, nil, nil, err
	}

	// Step 2: Fetch a random Rayon linked to the selected SR
	var regency structBot.Regency
	if err := db.Where("province_id = ?", province.ID).Limit(1).Find(&regency).Error; err != nil {
		return &province, nil, nil, nil, err
	}

	// Step 3: Fetch a random Branch linked to the selected Rayon
	var district structBot.District
	if err := db.Where("regency_id = ?", regency.ID).Limit(1).Find(&district).Error; err != nil {
		return &province, &regency, nil, nil, err
	}

	var subDistrict structBot.SubDistrict
	if err := db.Where("district_id = ?", district.ID).Limit(1).Find(&subDistrict).Error; err != nil {
		return &province, &regency, &district, nil, err
	}

	return &province, &regency, &district, &subDistrict, nil
}

func FetchMasters(db *gorm.DB) (map[string]interface{}, error) {

	var customer structBot.Customer
	if err := FetchRandomRow(db, "customer", &customer); err != nil {
		return nil, err
	}

	var customerType structBot.CustomerType
	if err := db.Where("id = ?", customer.CustomerTypeID).Limit(1).Find(&customerType).Error; err != nil {
		return nil, err
	}

	var transactionType structBot.TransactionType
	if err := FetchRandomRow(db, "transaction_type", &transactionType); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})

	result["customer"] = customer
	result["customerType"] = customerType
	result["transactionType"] = transactionType

	return result, nil
}

func FetchRandomUser(db *gorm.DB) (map[string]interface{}, error) {
	var user structBot.User
	if err := FetchRandomRow(db, "user", &user); err != nil {
		return nil, err
	}

	var userExecutor structBot.User
	if err := db.Where("user_id <> ? AND branch_id = ?", user.ID, user.BranchID).Limit(1).Find(&userExecutor).Error; err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	result["userHolder"] = user
	result["userExecutor"] = userExecutor

	return result, nil
}

func MoveData[T any](sourceDB *gorm.DB, targetDB *gorm.DB, query string, targetSlice *[]T) error {
	// Fetch data from the source database
	rows := []map[string]interface{}{}
	if err := sourceDB.Raw(query).Scan(&rows).Error; err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}

	// Prepare the target slice
	*targetSlice = make([]T, len(rows))

	// Loop through each row and map the data to the struct dynamically
	for i, row := range rows {
		// Get the struct type and value for dynamic assignment
		structType := reflect.TypeOf((*T)(nil)).Elem()
		structValue := reflect.New(structType).Elem()

		// Loop through each field in the struct and map the value from the row
		for j := 0; j < structValue.NumField(); j++ {
			field := structValue.Field(j)
			fieldName := structType.Field(j).Name

			// Match the struct field name to the column name (case insensitive)
			for key, value := range row {
				// Convert column name to match struct field name (snake_case to PascalCase)
				if strings.ToLower(key) == strings.ToLower(fieldName) {
					// Set the field value dynamically based on type
					if value != nil {
						// Convert value to the struct field's type
						field.Set(reflect.ValueOf(value))
					}
					break
				}
			}
		}

		// Assign the mapped struct to the target slice
		(*targetSlice)[i] = structValue.Interface().(T)
	}

	// Save the mapped data into the target database
	if err := targetDB.Save(targetSlice).Error; err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	return nil
}
func MoveData2[T any](sourceDB *gorm.DB, targetDB *gorm.DB, query string, targetSlice *[]T) ([]T, error) {
	// Fetch data from the source database
	rows := []map[string]interface{}{}
	if err := sourceDB.Raw(query).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	// Prepare the target slice
	*targetSlice = make([]T, len(rows))

	// Loop through each row and map the data to the struct dynamically
	for i, row := range rows {
		// Get the struct type and value for dynamic assignment
		structType := reflect.TypeOf((*T)(nil)).Elem()
		structValue := reflect.New(structType).Elem()

		// Loop through each field in the struct and map the value from the row
		for j := 0; j < structValue.NumField(); j++ {
			field := structValue.Field(j)
			fieldName := structType.Field(j).Name

			// Match the struct field name to the column name (case insensitive)
			for key, value := range row {
				// Convert column name to match struct field name (snake_case to PascalCase)
				if strings.ToLower(key) == strings.ToLower(fieldName) {
					// Set the field value dynamically based on type
					if value != nil {
						// Convert value to the struct field's type
						field.Set(reflect.ValueOf(value))
					}
					break
				}
			}
		}

		// Assign the mapped struct to the target slice
		(*targetSlice)[i] = structValue.Interface().(T)
	}

	// Save the mapped data into the target database
	// if err := targetDB.Save(targetSlice).Error; err != nil {
	// 	return nil ,fmt.Errorf("failed to save data: %w", err)
	// }

	return *targetSlice, nil
}

func TestGenerateDataRandom(c *fiber.Ctx) error {

	datas := []map[string]interface{}{}
	if err := db.DB.Raw("SELECT * FROM produk").Scan(&datas).Error; err != nil {
		return err
	}

	// fmt.Println(datas)

	datasBot := make([]structBot.Item, len(datas))
	for i := range datas {
		// datasBot[i].ID = int32(datas[i]["id"].(int64))
		datasBot[i].ID = datas[i]["id"].(int32)
		*datasBot[i].Name = datas[i]["name"].(string)
		*datasBot[i].ItemTypeID = 1
		// datasBot[i].ProvinceID = datas[i]["id_prov"].(string)

		// fmt.Println(datas[i]["id"].(string))

		// tempInt2, _ := strconv.Atoi(datas[i]["id"].(string))
		// datasBot[i].ID = int16(tempInt2)

		// fmt.Println(datasBot[i].ID)
		// tempInt, _ := strconv.Atoi(datas[i]["id_kec"].(string))
		// datasBot[i].DistrictID = int32(tempInt)
	}

	if err := db.DBBot.Save(&datasBot).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"datas": datasBot})

}

func getDatabaseSchema() string {
	var schemaBuilder strings.Builder

	rows, err := db.DB.Raw(`
		SELECT table_schema, table_name, column_name, data_type
		FROM information_schema.columns
		WHERE table_schema IN ('public', 'hr')
		ORDER BY table_schema, table_name, ordinal_position;
	`).Rows()
	if err != nil {
		log.Fatal("Failed to fetch schema:", err)
	}
	defer rows.Close()

	currentSchema := ""
	currentTable := ""
	for rows.Next() {
		var schemaName, tableName, columnName, dataType string
		rows.Scan(&schemaName, &tableName, &columnName, &dataType)

		if currentSchema != schemaName || currentTable != tableName {
			if currentSchema != "" || currentTable != "" {
				schemaBuilder.WriteString("\n")
			}
			currentSchema = schemaName
			currentTable = tableName
			schemaBuilder.WriteString(fmt.Sprintf("- Schema: %s, Table: %s\n", schemaName, tableName))
		}
		schemaBuilder.WriteString(fmt.Sprintf("    - %s (%s)\n", columnName, dataType))
	}

	return schemaBuilder.String()
}

func generateSQLQueryUsingGPT(userQuery string) string {
	schema := getDatabaseSchema() // Fetch schema information

	prompt := fmt.Sprintf(`
		Using the following database schema (note the schema names):
		%s
		Translate the following natural language query into an SQL query:
		"%s"
	`, schema, userQuery)

	resp, err := db.OpenaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a database assistant that generates SQL queries.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 150,
		},
	)
	if err != nil {
		log.Println("Error generating SQL query:", err)
		return ""
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content)
}

func executeSQLQuery(query string) ([]map[string]interface{}, error) {
	rows, err := db.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	columns, _ := rows.Columns()
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		row := map[string]interface{}{}
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	return results, nil
}

func summarizeQueryResults(results []map[string]interface{}) string {
	if len(results) == 0 {
		return "No data found."
	}

	var builder strings.Builder
	builder.WriteString("Results:\n")

	for i, row := range results {
		builder.WriteString(fmt.Sprintf("Row %d:\n", i+1))
		for key, value := range row {
			builder.WriteString(fmt.Sprintf("  - %s: %v\n", key, value))
		}
		if i == 4 { // Limit to 5 rows for readability
			builder.WriteString("  ...\n")
			break
		}
	}

	return builder.String()
}

func generateNaturalLanguageResponse(userQuery string, queryResult []map[string]interface{}) string {
	// Serialize query results into a string
	resultSummary := summarizeQueryResults(queryResult)

	// Build the prompt
	prompt := fmt.Sprintf(`
		The user asked: "%s"
		The following are the results of the SQL query:
		%s

		Explain the results in natural language, as if you were speaking to the user.
	`, userQuery, resultSummary)

	resp, err := db.OpenaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an assistant that explains database query results in a simple and natural way.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 150,
		},
	)
	if err != nil {
		log.Println("Error generating natural language response:", err)
		return "Sorry, I couldn't generate a natural language response."
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content)
}

func PostPrompt(c *fiber.Ctx) error {
	var body struct {
		UserQuery string `json:"userQuery"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	sqlQuery := generateSQLQueryUsingGPT(body.UserQuery)
	if sqlQuery == "" {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate SQL query"})
	}

	results, err := executeSQLQuery(sqlQuery)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to execute SQL query", "details": err.Error()})
	}

	naturalLanguageResponse := generateNaturalLanguageResponse(body.UserQuery, results)

	return c.JSON(fiber.Map{
		"query":    sqlQuery,
		"results":  results,
		"response": naturalLanguageResponse,
	})

}
