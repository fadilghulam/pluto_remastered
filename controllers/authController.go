package controllers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	db "pluto_remastered/config"
	"pluto_remastered/helpers"
	"pluto_remastered/structs"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func SendOtp(c *fiber.Ctx) error {
	type OtpRequest struct {
		Phone string `json:"phone"`
	}
	var otpReq OtpRequest
	if err := c.BodyParser(&otpReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Success: false,
			Message: "Something Went Wrong",
		})
	}

	if otpReq.Phone == "" {
		return c.Status(http.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Success: false,
			Message: "Phone number is required",
		})
	}

	params := map[string]interface{}{
		"sendTo":  otpReq.Phone,
		"appName": "PLUTO MOBILE",
	}

	dataSend, err := json.Marshal(params)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Error marshaling JSON",
			Success: false,
		})
	}

	responseData, err := helpers.SendCurl(dataSend, "POST", "https://rest.pt-bks.com/olympus/sendOtp")
	if err != nil {
		fmt.Println("Error sending request:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengirim notification",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": responseData["success"], "message": responseData["message"], "otpLength": 5})
}

func Login(c *fiber.Ctx) error {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	type LoginRequest struct {
		Username *string `json:"username"`
		Password *string `json:"password"`
		Otp      *string `json:"otp"`
		SendTo   *string `json:"sendTo"`
	}
	var loginReq LoginRequest
	if err := c.BodyParser(&loginReq); err != nil {
		fmt.Println(err.Error())
		return c.Status(http.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Success: false,
			Message: "Gagal mendapat input data",
		})
	}

	userSession := structs.UserSession{}

	if err := c.BodyParser(&userSession); err != nil {
		fmt.Println(err.Error())
		return c.Status(http.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Success: false,
			Message: "Gagal mendapat input data",
		})
	}

	// params := make(map[string]interface{})

	params := map[string]interface{}{
		"username": *loginReq.Username,
		"password": fmt.Sprintf("%x", md5.Sum([]byte(*loginReq.Password))),
		"appName":  "PLUTO MOBILE",
	}

	if loginReq.Otp != nil {
		params["otp"] = *loginReq.Otp
		params["sendTo"] = *loginReq.SendTo
	}

	dataSend, err := json.Marshal(params)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Error marshaling JSON",
			Success: false,
		})
	}

	returnedData, err := helpers.SendCurl(dataSend, "POST", "https://rest.pt-bks.com/olympus/login")
	if err != nil {
		fmt.Println("Error sending request:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengirim notification",
			Success: false,
		})
	}

	if returnedData["data"] != nil {

		tempUserID, _ := strconv.Atoi(returnedData["data"].(map[string]interface{})["id"].(string))

		*userSession.UserID = int32(tempUserID)

		tx := db.DB.Begin()

		if err := tx.Save(&userSession).Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal menyimpan data",
				Success: false,
			})
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal menyimpan data",
				Success: false,
			})
		}

	} else {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data user",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(returnedData)
}

func LoginAs(c *fiber.Ctx) error {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	type LoginAsRequest struct {
		ToUserId        *int32  `json:"toUserId"`
		FromUserId      *int32  `json:"fromUserId"`
		RequestDateTime *string `json:"requestDatetime"`
		IsDone          *int16  `json:"isDone"`
	}

	var loginReq LoginAsRequest
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Success: false,
			Message: "Gagal mendapat input data",
		})
	}

	userSession := structs.UserSession{}

	if err := c.BodyParser(&userSession); err != nil {
		fmt.Println(err.Error())
		return c.Status(http.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Success: false,
			Message: "Gagal mendapat input data",
		})
	}

	// params := make(map[string]interface{})

	params := map[string]interface{}{
		"toUserId":        *loginReq.ToUserId,
		"fromUserId":      *loginReq.FromUserId,
		"requestDatetime": *loginReq.RequestDateTime,
		"appName":         "Pluto.",
	}

	if loginReq.IsDone != nil {
		params["isDone"] = *loginReq.IsDone
	}

	dataSend, err := json.Marshal(params)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Error marshaling JSON",
			Success: false,
		})
	}

	returnedData, err := helpers.SendCurl(dataSend, "POST", "https://rest.pt-bks.com/olympus/login")
	if err != nil {
		fmt.Println("Error sending request:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data login",
			Success: false,
		})
	}

	if returnedData["data"] != nil {

		// tempUserID, _ := strconv.Atoi(returnedData["data"].(map[string]interface{})["id"].(string))

		userSession.UserID = loginReq.ToUserId
		userSession.UserIDSubtitute = *loginReq.FromUserId

		tx := db.DB.Begin()

		if err := tx.Save(&userSession).Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal menyimpan data",
				Success: false,
			})
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal menyimpan data",
				Success: false,
			})
		}

	} else {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data user",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(returnedData)
}

func RefreshDataUser(c *fiber.Ctx) error {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	type LoginRequest struct {
		RefreshUserId *int32  `json:"refreshUserId"`
		Mode          *string `json:"mode"`
	}
	var loginReq LoginRequest
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	params := make(map[string]interface{})
	if loginReq.RefreshUserId != nil {
		keyRefresh := os.Getenv("KEY_REFRESH_USER")
		params[keyRefresh] = *loginReq.RefreshUserId
		params["appName"] = "PLUTO MOBILE"
	}

	dataSend, err := json.Marshal(params)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Error marshaling JSON",
			Success: false,
		})
	}

	returnedData, err := helpers.SendCurl(dataSend, "POST", "https://rest.pt-bks.com/olympus/login")
	if err != nil {
		fmt.Println("Error sending request:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengirim notification",
			Success: false,
		})
	}

	dataMap := make(map[string]interface{})
	if loginReq.Mode != nil && strings.ToLower(*loginReq.Mode) == "permission" {
		dataMap["permission"] = returnedData["data"].(map[string]interface{})["permission"].(map[string]interface{})
		dataMap["serverTime"] = time.Now().UTC()
	} else {
		dataMap = returnedData
		dataMap["serverTime"] = time.Now().UTC()
	}

	dataMap["success"] = true
	dataMap["message"] = "Success"

	return c.Status(fiber.StatusOK).JSON(dataMap)
}

func Logout(c *fiber.Ctx) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	type LogoutRequest struct {
		UserID   *int32  `json:"userId"`
		DeviceID *string `json:"deviceId"`
		AppName  *string `json:"appName"`
	}

	var logoutReq LogoutRequest
	if err := c.BodyParser(&logoutReq); err != nil {
		fmt.Println(err.Error())
		return c.Status(http.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Success: false,
			Message: "Gagal mendapat input data",
		})
	}

	tokenFcm := new(structs.TokenFcm)
	if err := db.DB.
		Where("user_id = ? AND app_name = ? AND device_id = ?", *logoutReq.UserID, *logoutReq.AppName, *logoutReq.DeviceID).
		Delete(&tokenFcm).Error; err != nil {
		fmt.Println(err.Error())
		return c.Status(http.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Success: false,
			Message: "Gagal menghapus data",
		})
	}

	return c.Status(http.StatusOK).JSON(helpers.ResponseWithoutData{
		Success: true,
		Message: "Success menghapus token fcm",
	})
}
