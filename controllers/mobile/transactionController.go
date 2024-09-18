package controllers

import (
	"encoding/json"
	"fmt"
	db "pluto_remastered/config"
	"pluto_remastered/helpers"
	"pluto_remastered/structs"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func InsertTransactions(c *fiber.Ctx) error {

	type TemplateInputUser struct {
		Data      map[string]interface{} `json:"data"`
		DeletedID map[string]interface{} `json:"deletedId"`
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

	result := make(map[string][]map[string]interface{})

	tx := db.DB.Begin()

	for tableName, records := range inputUser.DeletedID {
		instanceSliceDelete, err := structs.GetStructInstanceByTableName(tableName)
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal mendapatkan tabel data",
				Success: false,
			})
		}

		whereIdIn := strings.Split(records.(string), ",")

		if err := tx.Clauses(clause.Returning{}).Where("id IN (?)", whereIdIn).Delete(instanceSliceDelete).Error; err != nil {
			tx.Rollback()
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal delete data",
				Success: false,
			})
		}

		recordsValue := reflect.ValueOf(instanceSliceDelete).Elem() // dereference the pointer to slice
		for i := 0; i < recordsValue.Len(); i++ {
			record := recordsValue.Index(i).Interface() // access the individual record

			// Use reflection to get id, sync_key, and created_at fields from the record
			id := reflect.ValueOf(record).FieldByName("ID").Interface()
			createdAtField := reflect.ValueOf(record).FieldByName("CreatedAt")
			dtmCrtField := reflect.ValueOf(record).FieldByName("DtmCrt")
			syncKeyField := reflect.ValueOf(record).FieldByName("SyncKey")
			var syncKey interface{}

			if createdAtField.IsValid() {
				syncKey = createdAtField.Interface()
			}

			if dtmCrtField.IsValid() {
				syncKey = dtmCrtField.Interface()
			}

			if syncKeyField.IsValid() {
				syncKey = syncKeyField.Interface()
			}

			result[tableName] = append(result[tableName], map[string]interface{}{
				"id":       id,
				"sync_key": syncKey,
			})
		}
	}

	tx.Commit()

	tx = db.DB.Begin()
	for tableName, records := range inputUser.Data {

		instanceSlice, err := structs.GetStructInstanceByTableName(tableName)
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal mendapatkan tabel data",
				Success: false,
			})
		}

		recordsBytes, err := json.Marshal(records)
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal konversi data tabel",
				Success: false,
			})
		}

		if err := json.Unmarshal(recordsBytes, instanceSlice); err != nil {
			tx.Rollback()
			// return c.Status(fiber.StatusBadRequest).SendString("Failed to parse records: " + err.Error())
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal konversi data tabel",
				Success: false,
			})
		}

		var tempIds []string
		recordsValue := reflect.ValueOf(instanceSlice).Elem() // dereference the pointer to slice
		for i := 0; i < recordsValue.Len(); i++ {
			record := recordsValue.Index(i).Interface() // access the individual record

			// Use reflection to get id, sync_key, and created_at fields from the record
			id := reflect.ValueOf(record).FieldByName("ID").Interface()

			tempIds = append(tempIds, fmt.Sprintf("%v", id))
		}

		if err := tx.Clauses(clause.Returning{}).Save(instanceSlice).Error; err != nil {
			tx.Rollback()
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal insert data",
				Success: false,
			})
		}

		tx.Where("id IN (?)", tempIds).Find(instanceSlice)

		recordsValue = reflect.ValueOf(instanceSlice).Elem() // dereference the pointer to slice
		for i := 0; i < recordsValue.Len(); i++ {
			record := recordsValue.Index(i).Interface() // access the individual record

			// Use reflection to get id, sync_key, and created_at fields from the record
			id := reflect.ValueOf(record).FieldByName("ID").Interface()
			createdAtField := reflect.ValueOf(record).FieldByName("CreatedAt")
			dtmCrtField := reflect.ValueOf(record).FieldByName("DtmCrt")
			syncKeyField := reflect.ValueOf(record).FieldByName("SyncKey")
			var syncKey interface{}

			if createdAtField.IsValid() {
				syncKey = createdAtField.Interface()
			}

			if dtmCrtField.IsValid() {
				syncKey = dtmCrtField.Interface()
			}

			if syncKeyField.IsValid() {
				syncKey = syncKeyField.Interface()
			}

			result[tableName] = append(result[tableName], map[string]interface{}{
				"id":       id,
				"sync_key": syncKey,
			})
		}
	}
	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    result,
	})
}
