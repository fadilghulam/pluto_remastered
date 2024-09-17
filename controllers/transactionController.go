package controllers

import (
	"fmt"
	"pluto_remastered/helpers"
	"pluto_remastered/structs"

	"github.com/gofiber/fiber/v2"
)

func InsertTransactions(c *fiber.Ctx) error {

	type TemplateInputUser struct {
		Data map[string]interface{} `json:"data"`
	}

	inputUser := new(TemplateInputUser)
	err := c.BodyParser(inputUser)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan input data",
			Success: false,
		})
	}

	for k, v := range inputUser.Data {
		fmt.Println(k, v)
		test, err := structs.GetStructInstanceByTableName(k)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Instance: %T\n", test)
	}

	return c.Status(fiber.StatusOK).JSON("Success")
}
