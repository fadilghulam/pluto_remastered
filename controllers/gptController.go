package controllers

import (
	"context"
	"fmt"
	"log"
	db "pluto_remastered/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sashabaranov/go-openai"
)

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
