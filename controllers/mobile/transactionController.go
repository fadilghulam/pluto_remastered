package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	db "pluto_remastered/config"
	"pluto_remastered/helpers"
	"pluto_remastered/structs"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

// func InsertTransactions(c *fiber.Ctx) error {

// 	type TemplateInputUser struct {
// 		Data map[string]interface{} `json:"transaction"`
// 	}

// 	inputUser := new(TemplateInputUser)
// 	err := c.BodyParser(inputUser)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 			Message: "Gagal mendapatkan input data",
// 			Success: false,
// 		})
// 	}

// 	tableInsert := make(map[string]interface{})
// 	tableDelete := make(map[string]interface{})
// 	for key, value := range inputUser.Data {
// 		if helpers.IsDeletedIds(key) {
// 			pattern := regexp.MustCompile(`_deleted_ids`)

// 			// Perform the replacement
// 			tempTableNameDelete := pattern.ReplaceAllString(key, "")

// 			tableDelete[tempTableNameDelete] = value
// 			// fmt.Printf("%v", value.([]interface{})[0])
// 		} else {
// 			tableInsert[key] = value
// 		}
// 	}

// 	// fmt.Println(tableInsert)
// 	// fmt.Println(tableDelete)

// 	result := make(map[string][]map[string]interface{})

// 	tx := db.DB.Begin()

// 	for tableName, records := range tableDelete {
// 		instanceSliceDelete, err := structs.GetStructInstanceByTableName(tableName)
// 		if err != nil {
// 			tx.Rollback()
// 			fmt.Println(err)
// 			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 				Message: "Gagal mendapatkan tabel data",
// 				Success: false,
// 			})
// 		}
// 		var tempIdsDelete []string

// 		for _, id := range records.([]interface{}) {
// 			if idStr, ok := id.(string); ok {
// 				tempIdsDelete = append(tempIdsDelete, idStr)
// 			}
// 		}

// 		if err := tx.Clauses(clause.Returning{}).Where("id IN (?)", tempIdsDelete).Delete(instanceSliceDelete).Error; err != nil {
// 			tx.Rollback()
// 			fmt.Println(err)
// 			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 				Message: "Gagal delete data",
// 				Success: false,
// 			})
// 		}

// 		recordsValue := reflect.ValueOf(instanceSliceDelete).Elem() // dereference the pointer to slice
// 		for i := 0; i < recordsValue.Len(); i++ {
// 			record := recordsValue.Index(i).Interface() // access the individual record

// 			// Use reflection to get id, sync_key, and created_at fields from the record
// 			id := reflect.ValueOf(record).FieldByName("ID").Interface()
// 			createdAtField := reflect.ValueOf(record).FieldByName("CreatedAt")
// 			dtmCrtField := reflect.ValueOf(record).FieldByName("DtmCrt")
// 			syncKeyField := reflect.ValueOf(record).FieldByName("SyncKey")
// 			var syncKey interface{}

// 			if createdAtField.IsValid() {
// 				syncKey = createdAtField.Interface()
// 			}

// 			if dtmCrtField.IsValid() {
// 				syncKey = dtmCrtField.Interface()
// 			}

// 			if syncKeyField.IsValid() {
// 				syncKey = syncKeyField.Interface()
// 			}

// 			result[tableName] = append(result[tableName], map[string]interface{}{
// 				"id":       id,
// 				"sync_key": syncKey,
// 			})
// 		}

// 		for i := 0; i < len(tempIdsDelete); i++ {
// 			found := false
// 			for _, value := range result[tableName] {
// 				if value["id"] == tempIdsDelete[i] {
// 					found = true
// 					break
// 				}
// 			}

// 			if found == false {
// 				result[tableName] = append(result[tableName], map[string]interface{}{
// 					"id":       tempIdsDelete[i],
// 					"sync_key": time.Now().Format("2006-01-02 15:04:05"),
// 				})
// 			}
// 		}
// 	}

// 	tx.Commit()

// 	tx = db.DB.Begin()
// 	for tableName, records := range tableInsert {

// 		instanceSlice, err := structs.GetStructInstanceByTableName(tableName)
// 		if err != nil {
// 			tx.Rollback()
// 			fmt.Println(err)
// 			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 				Message: "Gagal mendapatkan tabel data",
// 				Success: false,
// 			})
// 		}

// 		recordsBytes, err := json.Marshal(records)
// 		if err != nil {
// 			tx.Rollback()
// 			fmt.Println(err)
// 			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 				Message: "Gagal konversi data tabel",
// 				Success: false,
// 			})
// 		}

// 		if err := json.Unmarshal(recordsBytes, instanceSlice); err != nil {
// 			tx.Rollback()
// 			// return c.Status(fiber.StatusBadRequest).SendString("Failed to parse records: " + err.Error())
// 			fmt.Println(err)
// 			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 				Message: "Gagal konversi data tabel 2",
// 				Success: false,
// 			})
// 		}

// 		var tempIds []string
// 		recordsValue := reflect.ValueOf(instanceSlice).Elem() // dereference the pointer to slice
// 		for i := 0; i < recordsValue.Len(); i++ {
// 			record := recordsValue.Index(i).Interface() // access the individual record

// 			// Use reflection to get id, sync_key, and created_at fields from the record
// 			id := reflect.ValueOf(record).FieldByName("ID").Interface()

// 			tempIds = append(tempIds, fmt.Sprintf("%v", id))
// 		}

// 		if err := tx.Clauses(clause.Returning{}).Save(instanceSlice).Error; err != nil {
// 			tx.Rollback()
// 			fmt.Println(err)
// 			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 				Message: "Gagal insert data",
// 				Success: false,
// 			})
// 		}

// 		tx.Where("id IN (?)", tempIds).Find(instanceSlice)

// 		recordsValue = reflect.ValueOf(instanceSlice).Elem() // dereference the pointer to slice
// 		for i := 0; i < recordsValue.Len(); i++ {
// 			record := recordsValue.Index(i).Interface() // access the individual record

// 			// Use reflection to get id, sync_key, and created_at fields from the record
// 			id := reflect.ValueOf(record).FieldByName("ID").Interface()
// 			createdAtField := reflect.ValueOf(record).FieldByName("CreatedAt")
// 			dtmCrtField := reflect.ValueOf(record).FieldByName("DtmCrt")
// 			syncKeyField := reflect.ValueOf(record).FieldByName("SyncKey")
// 			var syncKey interface{}

// 			if createdAtField.IsValid() {
// 				syncKey = createdAtField.Interface()
// 			}

// 			if dtmCrtField.IsValid() {
// 				syncKey = dtmCrtField.Interface()
// 			}

// 			if syncKeyField.IsValid() {
// 				syncKey = syncKeyField.Interface()
// 			}

// 			result[tableName] = append(result[tableName], map[string]interface{}{
// 				"id":       id,
// 				"sync_key": syncKey,
// 			})
// 		}
// 	}
// 	tx.Commit()

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "success",
// 		"data":    result,
// 		// "data": nil,
// 	})
// }

//	func InsertTransactions(c *fiber.Ctx) error {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"message": "success",
//			"success": false,
//		})
//	}
func InsertTransactions(c *fiber.Ctx) error {

	type TemplateInputUser struct {
		Data map[string]interface{} `json:"transaction"`
		// Date string                 `json:"date"`
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

	// if inputUser.Date != time.Now().Format("2006-01-02") {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
	// 		Message: "Berbeda tanggal",
	// 		Success: false,
	// 	})
	// }

	// fmt.Println(inputUser)

	tableInsert := make(map[string]interface{})
	tableDelete := make(map[string]interface{})
	for key, value := range inputUser.Data {
		if helpers.IsDeletedIds(key) {
			pattern := regexp.MustCompile(`_deleted_ids`)

			// Perform the replacement
			tempTableNameDelete := pattern.ReplaceAllString(key, "")

			tableDelete[tempTableNameDelete] = value
			// fmt.Printf("%v", value.([]interface{})[0])
		} else {
			tableInsert[key] = value
		}
	}

	// fmt.Println(tableInsert)
	// fmt.Println(tableDelete)

	result := make(map[string][]map[string]interface{})

	tx := db.DB.Begin()

	for tableName, records := range tableDelete {
		instanceSliceDelete, err := structs.GetStructInstanceByTableName(tableName)
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal mendapatkan tabel data",
				Success: false,
			})
		}
		var tempIdsDelete []string

		for _, id := range records.([]interface{}) {
			if idStr, ok := id.(string); ok {
				tempIdsDelete = append(tempIdsDelete, idStr)
			}
		}

		if err := tx.Clauses(clause.Returning{}).Where("id IN (?)", tempIdsDelete).Delete(instanceSliceDelete).Error; err != nil {
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

		for i := 0; i < len(tempIdsDelete); i++ {
			found := false
			for _, value := range result[tableName] {
				if value["id"] == tempIdsDelete[i] {
					found = true
					break
				}
			}

			if found == false {
				result[tableName] = append(result[tableName], map[string]interface{}{
					"id":       tempIdsDelete[i],
					"sync_key": time.Now().Format("2006-01-02 15:04:05"),
				})
			}
		}
	}

	tx.Commit()

	keys := make([]string, 0, len(tableInsert))
	for key := range tableInsert {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	tx = db.DB.Begin()
	for _, key := range keys {

		for _, records := range tableInsert[key].([]interface{}) {

			// fmt.Println(records)
			// instanceSlice, err := structs.GetStructInstanceByTableName(key)
			// if err != nil {
			// 	tx.Rollback()
			// 	fmt.Println(err)
			// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			// 		Message: "Gagal mendapatkan tabel data",
			// 		Success: false,
			// 	})
			// }

			// recordsBytes, err := json.Marshal(records)
			// if err != nil {
			// 	tx.Rollback()
			// 	fmt.Println(err)
			// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			// 		Message: "Gagal konversi data tabel",
			// 		Success: false,
			// 	})
			// }

			// if err := json.Unmarshal(recordsBytes, instanceSlice); err != nil {
			// 	tx.Rollback()
			// 	// return c.Status(fiber.StatusBadRequest).SendString("Failed to parse records: " + err.Error())
			// 	fmt.Println(err)
			// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			// 		Message: "Gagal konversi data tabel 2",
			// 		Success: false,
			// 	})
			// }

			instanceSlice, structType, err := structs.GetStructInstanceByTableNameSingle(key)
			if err != nil {
				tx.Rollback()
				fmt.Println(err)
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal mendapatkan tabel data",
					Success: false,
				})
			}

			// Convert `records` into JSON
			recordsBytes, err := json.Marshal(records)
			if err != nil {
				tx.Rollback()
				fmt.Println(err)
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal konversi data tabel",
					Success: false,
				})
			}

			// Create a new instance of the struct
			element := reflect.New(structType).Interface()
			if err := json.Unmarshal(recordsBytes, element); err != nil {
				tx.Rollback()
				fmt.Println(err)
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal konversi data tabel 2",
					Success: false,
				})
			}

			reflect.ValueOf(instanceSlice).Elem().Set(reflect.Append(
				reflect.ValueOf(instanceSlice).Elem(),
				reflect.ValueOf(element).Elem(),
			))

			var tempIds []string
			var tempTime any
			recordsValue := reflect.ValueOf(instanceSlice).Elem() // dereference the pointer to slice
			for i := 0; i < recordsValue.Len(); i++ {
				record := recordsValue.Index(i).Interface() // access the individual record

				// Use reflection to get id, sync_key, and created_at fields from the record
				id := reflect.ValueOf(record).FieldByName("ID").Interface()
				deletedCreatedAtField := reflect.ValueOf(record).FieldByName("CreatedAt")
				deletedDtmCrtField := reflect.ValueOf(record).FieldByName("DtmCrt")

				tempIds = append(tempIds, fmt.Sprintf("%v", id))

				if deletedCreatedAtField.IsValid() {
					tempTime = deletedCreatedAtField.Interface()
				}

				if deletedDtmCrtField.IsValid() {
					tempTime = deletedDtmCrtField.Interface()
				}
			}

			dateNow := time.Now().Format("2006-01-02")

			convertedTempTime, err := time.Parse("2006-01-02T15:04:05", tempTime.(string))
			if err != nil {
				fmt.Println("ConvertedTempTime error:", err)
			}

			// if convertedTempTime.Format("2006-01-02") == dateNow && key != "customer" {
			if convertedTempTime.Format("2006-01-02") == dateNow {

				if err := tx.Clauses(clause.Returning{}).Save(instanceSlice).Error; err != nil {
					tx.Rollback()
					fmt.Println("error table " + key)
					// fmt.Println(instanceSlice)
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
					updatedAtField := reflect.ValueOf(record).FieldByName("UpdatedAt")
					dtmCrtField := reflect.ValueOf(record).FieldByName("DtmCrt")
					DtmUpdField := reflect.ValueOf(record).FieldByName("DtmUpd")
					syncKeyField := reflect.ValueOf(record).FieldByName("SyncKey")
					var syncKey interface{}
					var tempSyncKey interface{}

					if createdAtField.IsValid() {
						syncKey = createdAtField.Interface()
					}

					if dtmCrtField.IsValid() {
						syncKey = dtmCrtField.Interface()
					}

					if syncKeyField.IsValid() {
						syncKey = syncKeyField.Interface()
					}

					if updatedAtField.IsValid() {
						tempSyncKey = updatedAtField.Interface()
					}

					if DtmUpdField.IsValid() {
						tempSyncKey = DtmUpdField.Interface()
					}

					syncKeyTime, _ := syncKey.(time.Time)
					tempSyncKeyTime, _ := tempSyncKey.(time.Time)

					if tempSyncKeyTime.After(syncKeyTime) {
						syncKey = syncKeyTime
					} else {
						syncKey = tempSyncKeyTime
					}

					if syncKey.(time.Time).IsZero() {
						syncKey = time.Now().Format("2006-01-02 15:04:05")
					}

					result[key] = append(result[key], map[string]interface{}{
						"id":       id,
						"sync_key": syncKey,
					})
				}
			} else {
				for i := 0; i < recordsValue.Len(); i++ {
					record := recordsValue.Index(i).Interface()
					id := reflect.ValueOf(record).FieldByName("ID").Interface()

					result[key] = append(result[key], map[string]interface{}{
						"id":         id,
						"sync_key":   time.Now().Format("2006-01-02 15:04:05"),
						"deleted_at": time.Now().Format("2006-01-02 15:04:05"),
					})
				}
			}
		}
	}
	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    result,
		// "data": nil,
	})
}

// func CompleteSales(c *fiber.Ctx) error {

// 	// start := time.Now()
// 	type inputData struct {
// 		UserID          int32  `json:"userId"`
// 		UserIdSubtitute *int32 `json:"userIdSubtitute"`
// 		Date            string `json:"date"`
// 		ConfirmKey      string `json:"confirmKey"`
// 		BranchID        int16  `json:"branchId"`
// 	}
// 	var userInput inputData
// 	if err := c.BodyParser(&userInput); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseWithoutData{
// 			Success: false,
// 			Message: "Gagal mendapatkan input data",
// 		})
// 	}

// 	var dataSend, whereUpdate string
// 	if userInput.UserIdSubtitute != nil {
// 		dataSend = "?userId=" + strconv.Itoa(int(*userInput.UserIdSubtitute)) + "&date=" + userInput.Date
// 		whereUpdate = fmt.Sprintf("user_id_subtitute = %v AND user_id = %v AND DATE(tanggal_stok) = DATE('%s')", strconv.Itoa(int(*userInput.UserIdSubtitute)), userInput.UserID, userInput.Date)
// 	} else {
// 		dataSend = "?userId=" + strconv.Itoa(int(userInput.UserID)) + "&date=" + userInput.Date
// 		whereUpdate = fmt.Sprintf("user_id = %v AND (user_id_subtitute = 0 OR user_id_subtitute IS NULL) AND DATE(tanggal_stok) = DATE('%s')", userInput.UserID, userInput.Date)
// 	}

// 	// fmt.Println("https://rest.pt-bks.com/pluto-mobile/completeSalesQuery" + dataSend)

// 	responseData, err := helpers.SendCurl(nil, "GET", "https://rest.pt-bks.com/pluto-mobile/completeSalesQuery"+dataSend)
// 	if err != nil && err.Error() != "Not Found" {
// 		fmt.Println(err.Error() + " err 1")
// 		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 			Message: "Error mengambil data pengajuan",
// 			Success: false,
// 		})
// 	}

// 	if responseData["success"] == false {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": responseData["message"],
// 			"success": responseData["success"],
// 		})
// 	}

// 	stokUser := new(structs.StokUser)
// 	stokSalesman := new(structs.StokSalesman)
// 	stokMerchandiser := new(structs.StokMerchandiser)
// 	gudang := new(structs.Gudang)

// 	if err := db.DB.Where("branch_id = ? ", userInput.BranchID).First(&gudang).Error; err != nil && err.Error() != "record not found" {
// 		fmt.Println(err.Error())
// 		return c.Status(http.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 			Message: "Gagal mendapatkan data gudang",
// 			Success: false,
// 		})
// 	}
// 	if userInput.ConfirmKey != "3nagaterbangbersama" {
// 		// getData = $this->db->query("SELECT * FROM stok_salesman WHERE confirm_key = '$confirmKey'")->result_array();
// 		if err := db.DB.Where("confirm_key = ? AND gudang_id = ?", userInput.ConfirmKey, gudang.ID).Find(&stokUser).Error; err != nil {
// 			fmt.Println(err.Error() + " err 2")
// 			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 				Message: "Gagal mendapatkan data stok user",
// 				Success: false,
// 			})
// 		}

// 		if stokUser.ID != 0 {
// 			params := map[string]interface{}{
// 				"confirmKey": userInput.ConfirmKey,
// 			}

// 			dataSend, err := json.Marshal(params)
// 			if err != nil {
// 				fmt.Println("Error marshaling JSON:", err)
// 				return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
// 					Message: "Error marshaling JSON",
// 					Success: false,
// 				})
// 			}

// 			_, err = helpers.SendCurl(dataSend, "POST", "https://api.gudangku.pt-bks.com/order/closing-sales-notify")
// 			if err != nil && err.Error() != "Not Found" {
// 				fmt.Println("Error sending curl:", err)
// 			}

// 			tx := db.DB.Begin()

// 			if err := tx.Model(&stokUser).
// 				Where(whereUpdate+" AND confirm_key = ? AND gudang_id = ?", userInput.ConfirmKey, gudang.ID).
// 				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
// 				tx.Rollback()
// 				fmt.Println(err.Error())
// 				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 					Message: "Gagal so stok user",
// 					Success: false,
// 				})
// 			}

// 			if err := tx.Where(whereUpdate+" AND confirm_key = ? AND gudang_id = ?", userInput.ConfirmKey, gudang.ID).
// 				Find(&stokUser).Error; err != nil {
// 				tx.Rollback()
// 				fmt.Println(err.Error())
// 				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 					Message: "Gagal mendapat data so update",
// 					Success: false,
// 				})
// 			}

// 			if err := tx.Model(&stokSalesman).
// 				Where("stok_user_id = ?", stokUser.ID).
// 				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
// 				tx.Rollback()
// 				fmt.Println(err.Error())
// 				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 					Message: "Gagal so stok user",
// 					Success: false,
// 				})
// 			}

// 			if err := tx.Model(&stokMerchandiser).
// 				Where("stok_user_id = ?", stokUser.ID).
// 				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
// 				tx.Rollback()
// 				fmt.Println(err.Error())
// 				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 					Message: "Gagal mendapat data so update",
// 					Success: false,
// 				})
// 			}

// 			if err := tx.Commit().Error; err != nil {
// 				tx.Rollback()
// 				fmt.Println(err.Error())
// 				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 					Message: "Gagal so stok user",
// 					Success: false,
// 				})
// 			}

// 			// if tx.RowsAffected == 0 {
// 			// 	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
// 			// 		Message: "Tidak ada data so",
// 			// 		Success: true,
// 			// 	})
// 			// }
// 		}
// 	} else {
// 		params := map[string]interface{}{
// 			"confirmKey": userInput.ConfirmKey,
// 		}

// 		dataSend, err := json.Marshal(params)
// 		if err != nil {
// 			fmt.Println("Error marshaling JSON:", err)
// 			return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
// 				Message: "Error marshaling JSON",
// 				Success: false,
// 			})
// 		}

// 		_, err = helpers.SendCurl(dataSend, "POST", "https://api.gudangku.pt-bks.com/order/closing-sales-notify")
// 		if err != nil && err.Error() != "Not Found" {
// 			fmt.Println("Error sending curl:", err)
// 		}

// 		tx := db.DB.Begin()

// 		if err := tx.Model(&stokUser).
// 			Where(whereUpdate+" AND is_complete = 0 AND gudang_id = ?", gudang.ID).
// 			Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
// 			tx.Rollback()
// 			fmt.Println(err.Error() + "error 1")
// 			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 				Message: "Gagal so stok user",
// 				Success: false,
// 			})
// 		}

// 		if err := tx.Where(whereUpdate).
// 			Find(&stokUser).Error; err != nil {
// 			tx.Rollback()
// 			fmt.Println(err.Error() + "error 2")
// 			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 				Message: "Gagal mendapat data so update",
// 				Success: false,
// 			})
// 		}

// 		if err := tx.Model(&stokSalesman).
// 			Where("stok_user_id = ?", stokUser.ID).
// 			Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
// 			tx.Rollback()
// 			fmt.Println(err.Error() + "error 3")
// 			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 				Message: "Gagal so stok user",
// 				Success: false,
// 			})
// 		}

// 		if err := tx.Model(&stokMerchandiser).
// 			Where("stok_user_id = ?", stokUser.ID).
// 			Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
// 			tx.Rollback()
// 			fmt.Println(err.Error() + "error 4")
// 			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 				Message: "Gagal mendapat data so update",
// 				Success: false,
// 			})
// 		}

// 		if err := tx.Commit().Error; err != nil {
// 			tx.Rollback()
// 			fmt.Println(err.Error() + "error 5")
// 			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 				Message: "Gagal so stok user",
// 				Success: false,
// 			})
// 		}

// 		// if tx.RowsAffected == 0 {
// 		// 	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
// 		// 		Message: "Tidak ada data so",
// 		// 		Success: true,
// 		// 	})
// 		// }
// 	}

// 	sUserID := strconv.Itoa(int(userInput.UserID))
// 	sGudangID := strconv.Itoa(int(gudang.ID))
// 	var sUserIDSubtitute string
// 	if userInput.UserIdSubtitute != nil {
// 		sUserIDSubtitute = strconv.Itoa(int(*userInput.UserIdSubtitute))
// 	} else {
// 		sUserIDSubtitute = ""
// 	}
// 	datas, err := getStokParent(&sUserID, &userInput.Date, &sUserIDSubtitute, &sGudangID, c)

// 	if err != nil {
// 		return err
// 	}

// 	// fmt.Println(datas)

// 	return c.Status(fiber.StatusOK).JSON(helpers.Response{
// 		Message: "Berhasil so",
// 		Success: true,
// 		Data:    datas[0],
// 	})
// }

func CompleteSales(c *fiber.Ctx) error {

	// start := time.Now()
	type inputData struct {
		UserID          int32   `json:"userId"`
		UserIdSubtitute *int32  `json:"userIdSubtitute"`
		Date            string  `json:"date"`
		ConfirmKey      string  `json:"confirmKey"`
		BranchID        int16   `json:"branchId"`
		Type            *string `json:"type"`
	}
	var userInput inputData
	if err := c.BodyParser(&userInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Success: false,
			Message: "Gagal mendapatkan input data",
		})
	}

	var dataSend, whereUpdate string
	if userInput.UserIdSubtitute != nil {
		dataSend = "?userId=" + strconv.Itoa(int(*userInput.UserIdSubtitute)) + "&date=" + userInput.Date
		whereUpdate = fmt.Sprintf("user_id_subtitute = %v AND user_id = %v AND DATE(tanggal_stok) = DATE('%s')", strconv.Itoa(int(*userInput.UserIdSubtitute)), userInput.UserID, userInput.Date)
	} else {
		dataSend = "?userId=" + strconv.Itoa(int(userInput.UserID)) + "&date=" + userInput.Date
		whereUpdate = fmt.Sprintf("user_id = %v AND (user_id_subtitute = 0 OR user_id_subtitute IS NULL) AND DATE(tanggal_stok) = DATE('%s')", userInput.UserID, userInput.Date)
	}

	// fmt.Println("https://rest.pt-bks.com/pluto-mobile/completeSalesQuery" + dataSend)

	responseData, err := helpers.SendCurl(nil, "GET", "https://rest.pt-bks.com/pluto-mobile/completeSalesQuery"+dataSend)
	if err != nil && err.Error() != "Not Found" {
		fmt.Println(err.Error() + " err 1")
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Error mengambil data pengajuan",
			Success: false,
		})
	}

	if responseData["success"] == false {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": responseData["message"],
			"success": responseData["success"],
		})
	}

	stokUser := new(structs.StokUser)
	stokSalesman := new(structs.StokSalesman)
	stokMerchandiser := new(structs.StokMerchandiser)
	gudang := new(structs.Gudang)

	if err := db.DB.Where("branch_id = ? ", userInput.BranchID).First(&gudang).Error; err != nil && err.Error() != "record not found" {
		fmt.Println(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data gudang",
			Success: false,
		})
	}

	tx := db.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	if strings.TrimSpace(userInput.ConfirmKey) == "3nagaterbangbersama" {
		// if err := db.DB.Where("confirm_key = ? AND gudang_id = ?", userInput.ConfirmKey, gudang.ID).Find(&stokUser).Error; err != nil {
		// 	fmt.Println(err.Error() + " err 2")
		// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 		Message: "Gagal mendapatkan data stok user",
		// 		Success: false,
		// 	})
		// }

		// if stokUser.ID != 0 {
		// 	params := map[string]interface{}{
		// 		"confirmKey": userInput.ConfirmKey,
		// 	}

		// 	dataSend, err := json.Marshal(params)
		// 	if err != nil {
		// 		fmt.Println("Error marshaling JSON:", err)
		// 		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
		// 			Message: "Error marshaling JSON",
		// 			Success: false,
		// 		})
		// 	}

		// 	_, err = helpers.SendCurl(dataSend, "POST", "https://api.gudangku.pt-bks.com/order/closing-sales-notify")
		// 	if err != nil && err.Error() != "Not Found" {
		// 		fmt.Println("Error sending curl:", err)
		// 	}

		// 	tx := db.DB.Begin()

		// 	if err := tx.Model(&stokUser).
		// 		Where(whereUpdate+" AND confirm_key = ? AND gudang_id = ?", userInput.ConfirmKey, gudang.ID).
		// 		Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
		// 		tx.Rollback()
		// 		fmt.Println(err.Error())
		// 		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 			Message: "Gagal so stok user",
		// 			Success: false,
		// 		})
		// 	}

		// 	if err := tx.Where(whereUpdate+" AND confirm_key = ? AND gudang_id = ?", userInput.ConfirmKey, gudang.ID).
		// 		Find(&stokUser).Error; err != nil {
		// 		tx.Rollback()
		// 		fmt.Println(err.Error())
		// 		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 			Message: "Gagal mendapat data so update",
		// 			Success: false,
		// 		})
		// 	}

		// 	if err := tx.Model(&stokSalesman).
		// 		Where("stok_user_id = ?", stokUser.ID).
		// 		Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
		// 		tx.Rollback()
		// 		fmt.Println(err.Error())
		// 		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 			Message: "Gagal so stok user",
		// 			Success: false,
		// 		})
		// 	}

		// 	if err := tx.Model(&stokMerchandiser).
		// 		Where("stok_user_id = ?", stokUser.ID).
		// 		Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
		// 		tx.Rollback()
		// 		fmt.Println(err.Error())
		// 		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 			Message: "Gagal mendapat data so update",
		// 			Success: false,
		// 		})
		// 	}

		// 	if err := tx.Commit().Error; err != nil {
		// 		tx.Rollback()
		// 		fmt.Println(err.Error())
		// 		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 			Message: "Gagal so stok user",
		// 			Success: false,
		// 		})
		// 	}
		// }

		var tempId int64
		doUpdate := 1
		if userInput.Type != nil && *userInput.Type == "produk" {

			if err := tx.Model(&stokSalesman).
				Clauses(clause.Returning{}).
				Where(whereUpdate).
				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal update stok so",
					Success: false,
				})
			}

			if err := tx.Where(whereUpdate + " AND is_complete = 0").
				First(&stokMerchandiser).Error; err != nil && err.Error() != "record not found" {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal mendapat data",
					Success: false,
				})
			}

			if stokMerchandiser.ID != 0 {
				doUpdate = 0
			}

			tempId = *stokSalesman.StokUserId
		} else if userInput.Type != nil && *userInput.Type == "item" {
			if err := tx.Model(&stokMerchandiser).
				Clauses(clause.Returning{}).
				Where(whereUpdate).
				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal mendapat data so update",
					Success: false,
				})
			}

			if err := tx.Where(whereUpdate + " AND is_complete = 0").
				First(&stokSalesman).Error; err != nil && err.Error() != "record not found" {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal mendapat data",
					Success: false,
				})
			}

			if stokSalesman.ID != 0 {
				doUpdate = 0
			}

			tempId = *stokMerchandiser.StokUserId
		} else {
			if err := tx.Model(&stokSalesman).
				Clauses(clause.Returning{}).
				Where(whereUpdate).
				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal so stok user",
					Success: false,
				})
			}

			if err := tx.Model(&stokMerchandiser).
				Clauses(clause.Returning{}).
				Where(whereUpdate).
				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal mendapat data so update",
					Success: false,
				})
			}

			if *stokSalesman.StokUserId != 0 || stokSalesman.StokUserId != nil {
				tempId = *stokSalesman.StokUserId
			} else {
				tempId = *stokMerchandiser.StokUserId
			}
		}

		if doUpdate == 1 && tempId != 0 {
			if err := tx.Model(&stokUser).
				Where("id = ?", tempId).
				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal so stok user",
					Success: false,
				})
			}
		}

		// if err := tx.Commit().Error; err != nil {
		// 	// tx.Rollback()
		// 	fmt.Println(err.Error())
		// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 		Message: "Gagal so",
		// 		Success: false,
		// 	})
		// }
	} else {
		params := map[string]interface{}{
			"confirmKey": userInput.ConfirmKey,
		}

		dataSend, err := json.Marshal(params)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
				Message: "Error marshaling JSON",
				Success: false,
			})
		}

		_, err = helpers.SendCurl(dataSend, "POST", "https://api.gudangku.pt-bks.com/order/closing-sales-notify")
		if err != nil && err.Error() != "Not Found" {
			fmt.Println("Error sending curl:", err)
		}

		// if err := tx.Model(&stokUser).
		// 	Where(whereUpdate+" AND is_complete = 0 AND gudang_id = ?", gudang.ID).
		// 	Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
		// 	tx.Rollback()
		// 	fmt.Println(err.Error() + "error 1")
		// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 		Message: "Gagal so stok user",
		// 		Success: false,
		// 	})
		// }

		// if err := tx.Where(whereUpdate).
		// 	Find(&stokUser).Error; err != nil {
		// 	tx.Rollback()
		// 	fmt.Println(err.Error() + "error 2")
		// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 		Message: "Gagal mendapat data so update",
		// 		Success: false,
		// 	})
		// }

		// if err := tx.Model(&stokSalesman).
		// 	Where("stok_user_id = ?", stokUser.ID).
		// 	Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
		// 	tx.Rollback()
		// 	fmt.Println(err.Error() + "error 3")
		// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 		Message: "Gagal so stok user",
		// 		Success: false,
		// 	})
		// }

		// if err := tx.Model(&stokMerchandiser).
		// 	Where("stok_user_id = ?", stokUser.ID).
		// 	Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
		// 	tx.Rollback()
		// 	fmt.Println(err.Error() + "error 4")
		// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 		Message: "Gagal mendapat data so update",
		// 		Success: false,
		// 	})
		// }

		var tempId int64
		doUpdate := 1
		if userInput.Type != nil && *userInput.Type == "produk" {

			tempCheck := tx.Model(&stokSalesman).
				Where(whereUpdate+" AND confirm_key = ?", userInput.ConfirmKey).
				Find(&stokSalesman)

			if tempCheck.Error != nil || tempCheck.RowsAffected == 0 {
				tx.Rollback()
				return c.Status(fiber.StatusNoContent).JSON(helpers.ResponseWithoutData{
					Message: "Data tidak ditemukan",
					Success: false,
				})
			}

			if err := tx.Model(&stokSalesman).
				Clauses(clause.Returning{}).
				Where(whereUpdate+" AND confirm_key = ?", userInput.ConfirmKey).
				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal so stok user",
					Success: false,
				})
			}

			if err := tx.Where(whereUpdate + " AND is_complete = 0").
				First(&stokMerchandiser).Error; err != nil && err.Error() != "record not found" {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal mendapat data",
					Success: false,
				})
			}

			if stokMerchandiser.ID != 0 {
				doUpdate = 0
			}

			tempId = *stokSalesman.StokUserId
		} else if userInput.Type != nil && *userInput.Type == "item" {

			tempCheck := tx.Model(&stokMerchandiser).
				Where(whereUpdate+" AND confirm_key = ?", userInput.ConfirmKey).
				Find(&stokMerchandiser)

			if tempCheck.Error != nil || tempCheck.RowsAffected == 0 {
				tx.Rollback()
				return c.Status(fiber.StatusNoContent).JSON(helpers.ResponseWithoutData{
					Message: "Data tidak ditemukan",
					Success: false,
				})
			}

			if err := tx.Model(&stokMerchandiser).
				Clauses(clause.Returning{}).
				Where(whereUpdate+" AND confirm_key = ?", userInput.ConfirmKey).
				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal mendapat data so update",
					Success: false,
				})
			}

			if err := tx.Where(whereUpdate + " AND is_complete = 0").
				First(&stokSalesman).Error; err != nil && err.Error() != "record not found" {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal mendapat data",
					Success: false,
				})
			}

			if stokSalesman.ID != 0 {
				doUpdate = 0
			}

			tempId = *stokMerchandiser.StokUserId
		} else {

			tempCheck := tx.Model(&stokSalesman).
				Where(whereUpdate+" AND confirm_key = ?", userInput.ConfirmKey).
				Find(&stokSalesman)

			tempCheck2 := tx.Model(&stokMerchandiser).
				Where(whereUpdate+" AND confirm_key = ?", userInput.ConfirmKey).
				Find(&stokMerchandiser)

			if tempCheck.Error != nil || (tempCheck.RowsAffected == 0 && tempCheck2.RowsAffected == 0) {
				tx.Rollback()
				return c.Status(fiber.StatusNoContent).JSON(helpers.ResponseWithoutData{
					Message: "Data tidak ditemukan",
					Success: false,
				})
			}

			if err := tx.Model(&stokSalesman).
				Clauses(clause.Returning{}).
				Where(whereUpdate+" AND confirm_key = ?", userInput.ConfirmKey).
				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal so stok user",
					Success: false,
				})
			}

			if err := tx.Model(&stokMerchandiser).
				Clauses(clause.Returning{}).
				Where(whereUpdate+" AND confirm_key = ?", userInput.ConfirmKey).
				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal mendapat data so update",
					Success: false,
				})
			}

			if *stokSalesman.StokUserId != 0 || stokSalesman.StokUserId != nil {
				tempId = *stokSalesman.StokUserId
			} else {
				tempId = *stokMerchandiser.StokUserId
			}
		}

		if doUpdate == 1 && tempId != 0 {
			if err := tx.Model(&stokUser).
				Where("id = ?", tempId).
				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal so stok user",
					Success: false,
				})
			}
		}

		// if err := tx.Commit().Error; err != nil {
		// 	tx.Rollback()
		// 	fmt.Println(err.Error() + "error 5")
		// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
		// 		Message: "Gagal so stok user",
		// 		Success: false,
		// 	})
		// }
	}

	sUserID := strconv.Itoa(int(userInput.UserID))
	sGudangID := strconv.Itoa(int(gudang.ID))
	var sUserIDSubtitute string
	if userInput.UserIdSubtitute != nil {
		sUserIDSubtitute = strconv.Itoa(int(*userInput.UserIdSubtitute))
	} else {
		sUserIDSubtitute = ""
	}
	datas, err := getStokParent(&sUserID, &userInput.Date, &sUserIDSubtitute, &sGudangID, c)

	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data stok terbaru",
			Success: false,
		})
	}

	if err := tx.Commit().Error; err != nil {
		// tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal so",
			Success: false,
		})
	}

	// fmt.Println(datas)

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Berhasil so",
		Success: true,
		Data:    datas[0],
	})
}
