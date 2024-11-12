package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	db "pluto_remastered/config"
	"pluto_remastered/helpers"
	"pluto_remastered/structs"
	"reflect"
	"strconv"
	"strings"
	"time"

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

func CompleteSales(c *fiber.Ctx) error {

	// start := time.Now()
	type inputData struct {
		UserID          int32  `json:"userId"`
		UserIdSubtitute *int32 `json:"userIdSubtitute"`
		Date            string `json:"date"`
		ConfirmKey      string `json:"confirmKey"`
		BranchID        int16  `json:"branchId"`
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
	if userInput.ConfirmKey != "3nagaterbangbersama" {
		// getData = $this->db->query("SELECT * FROM stok_salesman WHERE confirm_key = '$confirmKey'")->result_array();
		if err := db.DB.Where("confirm_key = ? AND gudang_id = ?", userInput.ConfirmKey, gudang.ID).Find(&stokUser).Error; err != nil {
			fmt.Println(err.Error() + " err 2")
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal mendapatkan data stok user",
				Success: false,
			})
		}

		if stokUser.ID != 0 {
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

			tx := db.DB.Begin()

			if err := tx.Model(&stokUser).
				Where(whereUpdate+" AND confirm_key = ? AND gudang_id = ?", userInput.ConfirmKey, gudang.ID).
				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal so stok user",
					Success: false,
				})
			}

			if err := tx.Where(whereUpdate+" AND confirm_key = ? AND gudang_id = ?", userInput.ConfirmKey, gudang.ID).
				Find(&stokUser).Error; err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal mendapat data so update",
					Success: false,
				})
			}

			if err := tx.Model(&stokSalesman).
				Where("stok_user_id = ?", stokUser.ID).
				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal so stok user",
					Success: false,
				})
			}

			if err := tx.Model(&stokMerchandiser).
				Where("stok_user_id = ?", stokUser.ID).
				Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal mendapat data so update",
					Success: false,
				})
			}

			if err := tx.Commit().Error; err != nil {
				tx.Rollback()
				fmt.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal so stok user",
					Success: false,
				})
			}

			// if tx.RowsAffected == 0 {
			// 	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			// 		Message: "Tidak ada data so",
			// 		Success: true,
			// 	})
			// }
		}
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

		tx := db.DB.Begin()

		if err := tx.Model(&stokUser).
			Where(whereUpdate+" AND is_complete = 0 AND gudang_id = ?", gudang.ID).
			Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error() + "error 1")
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal so stok user",
				Success: false,
			})
		}

		if err := tx.Where(whereUpdate+" AND gudang_id = ?", gudang.ID).
			Find(&stokUser).Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error() + "error 2")
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal mendapat data so update",
				Success: false,
			})
		}

		if err := tx.Model(&stokSalesman).
			Where("stok_user_id = ?", stokUser.ID).
			Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error() + "error 3")
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal so stok user",
				Success: false,
			})
		}

		if err := tx.Model(&stokMerchandiser).
			Where("stok_user_id = ?", stokUser.ID).
			Updates(map[string]interface{}{"confirm_key": nil, "tanggal_so": time.Now(), "is_complete": 1}).Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error() + "error 4")
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal mendapat data so update",
				Success: false,
			})
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error() + "error 5")
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal so stok user",
				Success: false,
			})
		}

		// if tx.RowsAffected == 0 {
		// 	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
		// 		Message: "Tidak ada data so",
		// 		Success: true,
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
		return err
	}

	// fmt.Println(datas)

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Berhasil so",
		Success: true,
		Data:    datas[0],
	})
}
