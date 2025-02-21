package controllers

import (
	"fmt"
	"pluto_remastered/helpers"
	"time"

	db "pluto_remastered/config"
	"pluto_remastered/structs"

	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	// "gorm.io/gorm/clause"
)

func UpdatePengiriman(c *fiber.Ctx) error {

	type Date struct {
		Receive string `json:"receive"`
		Send    string `json:"send"`
	}

	type Product struct {
		Code   string                      `json:"code"`
		Detail []structs.StokGudangRiwayat `json:"detail"`
		ID     int64                       `json:"id"`
		Pita   string                      `json:"pita"`
		Qty    int                         `json:"qty"`
	}

	type RefNumber struct {
		System string `json:"system"`
		Vendor string `json:"vendor"`
	}

	type TemplateInputUser struct {
		Date                   Date      `json:"date"`
		LatestUpdate           *string   `json:"latest_update"`
		Product                []Product `json:"product"`
		Receiver               string    `json:"receiver"`
		RefNumber              RefNumber `json:"ref_number"`
		Sender                 string    `json:"sender"`
		Status                 string    `json:"status"`
		StokGudangPengirimanID int64     `json:"stok_gudang_pengiriman_id"`
		Tag                    string    `json:"tag"`
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

	// fmt.Println(inputUser.Data)

	stokGudangRiwayat := new(structs.StokGudangRiwayat)

	tx := db.DB.Begin()
	i := 0
	for _, value := range inputUser.Product {
		for _, dataRiwayat := range value.Detail {

			if err := mapstructure.Decode(dataRiwayat, &stokGudangRiwayat); err != nil {
				fmt.Println(err.Error())
			}

			if err := tx.Where("id = ?", stokGudangRiwayat.ID).Delete(stokGudangRiwayat).Error; err != nil {
				tx.Rollback()
				fmt.Println(err)
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal menghapus data",
					Success: false,
				})
			}

			output := fmt.Sprintf("%d%s%d", stokGudangRiwayat.StokGudangPengirimanID, "", i)
			i++

			stokGudangRiwayat.ID = helpers.ConvertStringToInt64(output)

			// fmt.Println(stokGudangRiwayat)

			if err := tx.Save(stokGudangRiwayat).Error; err != nil {
				tx.Rollback()
				fmt.Println(err)
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Gagal menyimpan update data",
					Success: false,
				})
			}
		}
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	if err := tx.Model(&structs.StokGudangPengiriman{}).
		Where("id = ?", inputUser.StokGudangPengirimanID).
		Updates(structs.StokGudangPengiriman{
			UpdatePriceAt: &currentTime,
			DtmUpd:        time.Now().Format("2006-01-02 15:04:05"),
		}).Error; err != nil {
		tx.Rollback()
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal menyimpan update data pengiriman",
			Success: false,
		})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal menyimpan data",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
		Message: "Berhasil menyimpan data",
		Success: true,
	})
}
