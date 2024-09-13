package controllers

import (
	"encoding/json"
	"fmt"
	db "pluto_remastered/config"
	"pluto_remastered/helpers"
	"pluto_remastered/structs"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetGudang(c *fiber.Ctx) error {
	type TemplateInputUser struct {
		GudangId *string `json:"gudangId"`
		BranchId *string `json:"branchId"`
	}

	inputUser := new(TemplateInputUser)
	err := c.QueryParser(inputUser)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan input data",
			Success: false,
		})
	}

	where := ""

	if inputUser.GudangId != nil {
		where = where + " AND g.id IN (" + *inputUser.GudangId + ")"
	}

	if inputUser.BranchId != nil {
		where = where + " AND g.branch_id IN (" + *inputUser.BranchId + ")"
	}

	datas, err := helpers.NewExecuteQuery(fmt.Sprintf(`SELECT g.id, 
															g.name, 
															g.deskripsi, 
															g.sr_id, 
															g.rayon_id, 
															g.branch_id
														FROM gudang g
														WHERE g.is_salesman = 1 %s ORDER BY id`, where))

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data",
			Success: false,
		})
	}

	if len(datas) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data not found",
			Success: true,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data has been loaded",
		Success: true,
		Data:    datas,
	})
}

func GetProdukByGudang(c *fiber.Ctx) error {

	type TemplateInputUser struct {
		GudangId  *string `json:"gudangId"`
		ProdukId  *string `json:"produkId"`
		Pita      *string `json:"pita"`
		Condition *string `json:"condition"`
	}

	inputUser := new(TemplateInputUser)
	err := c.QueryParser(inputUser)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan input data",
			Success: false,
		})
	}

	where := ""

	if inputUser.GudangId != nil {
		where = where + " AND gudang_id IN (" + *inputUser.GudangId + ")"
	}

	if inputUser.ProdukId != nil {
		where = where + " AND produk_id IN (" + *inputUser.ProdukId + ")"
	}

	if inputUser.Pita != nil {
		where = where + " AND pita IN (" + *inputUser.Pita + ")"
	}

	if inputUser.Condition != nil {
		where = where + " AND condition = '" + *inputUser.Condition + "'"
	}

	datas, err := helpers.NewExecuteQuery(fmt.Sprintf(`SELECT sg.id, 
												JSONB_BUILD_OBJECT('id', sg.produk_id, 'name', p.name, 'code', p.code, 'photo', p.foto) as produk,
												harga,
												jumlah,
												batch,
												condition,
												pita,
												gudang_id
											FROM stok_gudang sg 
											JOIN produk p
												ON sg.produk_id = p.id
											WHERE TRUE %s 
											ORDER BY pita DESC`, where))

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data",
			Success: false,
		})
	}

	if len(datas) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data not found",
			Success: true,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data has been loaded",
		Success: true,
		Data:    datas,
	})
}

func GetItemByGudang(c *fiber.Ctx) error {

	type TemplateInputUser struct {
		GudangId *string `json:"gudangId"`
		ItemId   *string `json:"itemId"`
	}

	inputUser := new(TemplateInputUser)
	err := c.QueryParser(inputUser)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan input data",
			Success: false,
		})
	}

	where := ""

	if inputUser.GudangId != nil {
		where = where + " AND gudang_id IN (" + *inputUser.GudangId + ")"
	}

	if inputUser.ItemId != nil {
		where = where + " AND item_id IN (" + *inputUser.ItemId + ")"
	}

	datas, err := helpers.NewExecuteQuery(fmt.Sprintf(`SELECT sg.id, 
															sg.item_id as item_id,
															i.name as item_name,
															i.code as item_code,
															i.category_id as category_id,
															ic.name as category_name,
															i.brand_id,
															pb.name as brand_name,
															harga,
															jumlah,
															batch,
															gudang_id
														FROM md.stok_gudang_item sg 
														JOIN md.item i
															ON sg.item_id = i.id
														JOIN md.item_category ic
															ON i.category_id = ic.id
														JOIN produk_brand pb
															ON i.brand_id = pb.id
														WHERE TRUE %s 
														ORDER BY i.id`, where))

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data",
			Success: false,
		})
	}

	if len(datas) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data not found",
			Success: true,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data has been loaded",
		Success: true,
		Data:    datas,
	})
}

func ConfirmOrder(c *fiber.Ctx) error {

	inputUser := new(structs.StokSalesmanRiwayat)
	err := c.BodyParser(inputUser)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan input data",
			Success: false,
		})
	}

	whereConfirmKey := ""
	if inputUser.ConfirmKey != "3kucingjantan" {
		whereConfirmKey = " AND confirm_key = '" + inputUser.ConfirmKey + "'"
	}

	typeVar := c.FormValue("type", "PRODUK")
	stokSalesmanRiwayats := []structs.StokSalesmanRiwayat{}
	stokMerchandiserRiwayats := []structs.StokMerchandiserRiwayat{}

	if typeVar == "PRODUK" {
		err = db.DB.Where("parent_id = ? AND user_id = ?"+whereConfirmKey, inputUser.ParentId, inputUser.UserId).Find(&stokSalesmanRiwayats).Error
		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal mendapatkan data order / retur",
				Success: false,
			})
		}
	} else {
		err = db.DB.Where("parent_id = ? AND user_id = ?"+whereConfirmKey, inputUser.ParentId, inputUser.UserId).Find(&stokMerchandiserRiwayats).Error
		if err != nil {
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Gagal mendapatkan data order / retur",
				Success: false,
			})
		}
	}

	if len(stokSalesmanRiwayats) == 0 || len(stokMerchandiserRiwayats) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data order / retur tidak ditemukan",
			Success: false,
		})
	}

	if len(stokSalesmanRiwayats) > 0 {
		db.DB.Where("parent_id = ? AND user_id = ?"+whereConfirmKey, inputUser.ParentId, inputUser.UserId).Updates(structs.StokSalesmanRiwayat{
			IsValidate: 1,
			ConfirmKey: "",
			DtmUpd:     time.Now(),
		})
	}

	if len(stokMerchandiserRiwayats) > 0 {
		db.DB.Where("parent_id = ? AND user_id = ?"+whereConfirmKey, inputUser.ParentId, inputUser.UserId).Updates(structs.StokMerchandiserRiwayat{
			IsValidate: 1,
			ConfirmKey: "",
			DtmUpd:     time.Now(),
		})
	}

	params := map[string]interface{}{
		"parentId": inputUser.ParentId,
	}

	dataSend, err := json.Marshal(params)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Error marshaling JSON",
			Success: false,
		})
	}

	_, err = helpers.SendCurl(dataSend, "POST", "https://api.gudangku.pt-bks.com/order/product-notify")
	if err != nil {
		fmt.Println("Error sending request:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengirim notification",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
		Message: "Order barhasil dikonfirmasi",
		Success: true,
	})

}
