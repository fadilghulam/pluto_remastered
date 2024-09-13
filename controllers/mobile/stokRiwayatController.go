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

func GetListOrder(c *fiber.Ctx) error {
	type TemplateInputUser struct {
		UserId   *string `json:"userId"`
		ChildId  *string `json:"childId"`
		ParentId *string `json:"parentId"`
		Type     *string `json:"type"`
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
	if inputUser.UserId != nil {
		where = where + " AND ssr.user_id IN ( " + *inputUser.UserId + " )"
	}
	if inputUser.ChildId != nil {
		where = where + " AND ssr.id IN ( " + *inputUser.ChildId + " )"
	}
	if inputUser.ParentId != nil {
		where = where + " AND ssr.parent_id IN ( " + *inputUser.ParentId + " )"
	}
	if inputUser.Type != nil {
		where = where + " AND UPPER(ssr.type) = UPPER('" + *inputUser.Type + "')"
	}

	data, err := helpers.NewExecuteQuery(fmt.Sprintf(`SELECT ssr.parent_id as order_id, 
											CASE WHEN MIN(ssr.is_validate) = 1 THEN 'Approve'
												WHEN MIN(ssr.is_validate) = 0 AND ssr.confirm_key IS NOT NULL THEN 'Processed' END as status,
											ssr.condition,
											ssr.gudang_id,
											ssr.salesman_id,
											ssr.tanggal_riwayat,
											ssr.aksi,
											JSONB_AGG( DISTINCT
												JSONB_BUILD_OBJECT(
													'id_order_child', ssr.id,
													'produk', JSONB_BUILD_OBJECT(
																			'id_produk', p.id,
																			'code', p.code,
																			'name', p.name,
																			'photo', p.foto
																		),
													'jumlah', ssr.jumlah,
													'condition', ssr.condition,
													'pita', ssr.pita,
													'gudang_id', ssr.gudang_id
												) --ORDER BY ssr.pita DESC, ssr.id
											) as datas
										FROM stok_salesman_riwayat ssr
										JOIN produk p
											ON ssr.produk_id = p.id
										JOIN stok_gudang sg
											ON p.id = sg.produk_id
											AND ssr.condition = sg.condition
											AND ssr.pita = sg.pita
										WHERE TRUE AND ssr.parent_id IS NOT NULL AND DATE(tanggal_riwayat) BETWEEN CURRENT_DATE -'1 month'::interval AND CURRENT_DATE %s
										GROUP BY ssr.parent_id, ssr.confirm_key, ssr.condition, ssr.gudang_id, ssr.salesman_id, ssr.tanggal_riwayat, ssr.aksi
										ORDER BY ssr.tanggal_riwayat DESC`, where))

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data",
			Success: false,
		})
	}

	if len(data) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data tidak ada",
			Success: true,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data has been loaded",
		Success: true,
		Data:    data,
	})
}

func GetListOrderMD(c *fiber.Ctx) error {
	type TemplateInputUser struct {
		UserId   *string `json:"userId"`
		ChildId  *string `json:"childId"`
		ParentId *string `json:"parentId"`
		Type     *string `json:"type"`
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
	if inputUser.UserId != nil {
		where = where + " AND ssr.user_id IN ( " + *inputUser.UserId + " )"
	}
	if inputUser.ChildId != nil {
		where = where + " AND ssr.id IN ( " + *inputUser.ChildId + " )"
	}
	if inputUser.ParentId != nil {
		where = where + " AND ssr.parent_id IN ( " + *inputUser.ParentId + " )"
	}
	if inputUser.Type != nil {
		where = where + " AND UPPER(ssr.type) = UPPER('" + *inputUser.Type + "')"
	}

	data, err := helpers.NewExecuteQuery(fmt.Sprintf(`SELECT ssr.parent_id as order_id, 
														CASE WHEN MIN(ssr.is_validate) = 1 THEN 'Approve'
															WHEN MIN(ssr.is_validate) = 0 AND ssr.confirm_key IS NOT NULL THEN 'Processed' END as status,
														ssr.gudang_id,
														ssr.merchandiser_id,
														ssr.tanggal_riwayat,
														ssr.aksi,
														JSONB_AGG( DISTINCT
															JSONB_BUILD_OBJECT(
																'id_order_child', ssr.id,
																'item', JSONB_BUILD_OBJECT(
																			'id_produk', i.id,
																			'code', i.code,
																			'name', i.name,
																			'category_id', ic.id,
																			'category_name', ic.name,
																			'brand_id', pb.id,
																			'brand_name', pb.name
																		),
																'jumlah', ssr.jumlah,
																'gudang_id', ssr.gudang_id
															) --ORDER BY ssr.pita DESC, ssr.id
														) as datas
													FROM md.stok_merchandiser_riwayat ssr
													JOIN md.item i
														ON ssr.item_id = i.id
													JOIN md.item_category ic
														ON i.category_id = ic.id
													JOIN produk_brand pb
														ON i.brand_id = pb.id
													JOIN md.stok_gudang_item sg
														ON i.id = sg.item_id
													WHERE TRUE AND ssr.parent_id IS NOT NULL AND DATE(tanggal_riwayat) BETWEEN CURRENT_DATE -'1 month'::interval AND CURRENT_DATE %s
													GROUP BY ssr.parent_id, ssr.confirm_key, ssr.gudang_id, ssr.merchandiser_id, ssr.tanggal_riwayat, ssr.aksi
													ORDER BY ssr.tanggal_riwayat DESC`, where))

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data",
			Success: false,
		})
	}

	if len(data) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data tidak ada",
			Success: true,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data has been loaded",
		Success: true,
		Data:    data,
	})
}

func PostOrder(c *fiber.Ctx) error {

	type Products struct {
		Id        *string `json:"id"`
		Qty       *string `json:"qty"`
		Pita      *string `json:"pita"`
		Condition *string `json:"condition"`
		Aksi      *string `json:"aksi"`
	}

	type TemplateInputUser struct {
		Date     *string    `json:"date"`
		GudangId *string    `json:"gudangId"`
		UserId   *string    `json:"userId"`
		Type     *string    `json:"type"`
		Products []Products `json:"products"`
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

	var stokSalesmanRiwayat []structs.StokSalesmanRiwayat

	parentId := int64(helpers.ParseInt(*inputUser.UserId)) + time.Now().Unix()
	for i := 0; i < len(inputUser.Products); i++ {
		// 	tempString := *inputUser.Products[i].Id + "-" + *inputUser.Products[i].Pita + "-" + *inputUser.Products[i].Condition + "-" + *inputUser.Products[i].Aksi + *inputUser.UserId

		// 	if
		stokSalesmanRiwayat = append(stokSalesmanRiwayat, structs.StokSalesmanRiwayat{
			ProdukId:       int16(helpers.ParseInt(*inputUser.Products[i].Id)),
			Jumlah:         int32(helpers.ParseInt(*inputUser.Products[i].Qty)),
			Pita:           int32(helpers.ParseInt(*inputUser.Products[i].Pita)),
			Condition:      *inputUser.Products[i].Condition,
			IsValidate:     0,
			GudangId:       int16(helpers.ParseInt(*inputUser.GudangId)),
			UserId:         int32(helpers.ParseInt(*inputUser.UserId)),
			TanggalRiwayat: helpers.ParseDate(*inputUser.Date),
			ParentId:       parentId,
			Aksi:           *inputUser.Products[i].Aksi,
		})
	}

	tx := db.DB.Begin()

	err = tx.Create(&stokSalesmanRiwayat).Error
	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan input data",
			Success: false,
		})
	}

	tx.Commit()

	params := map[string]interface{}{
		"warehouseId": *inputUser.GudangId,
	}

	dataSend, err := json.Marshal(params)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Error marshaling JSON",
			Success: false,
		})
	}

	_, err = helpers.SendCurl(dataSend, "POST", "https://api.gudangku.pt-bks.com/order/new-order-product-notify")
	if err != nil {
		fmt.Println("Error sending request:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengirim notification",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
		Message: "Success insert data",
		Success: true,
	})
}

func PostOrderMD(c *fiber.Ctx) error {

	type Items struct {
		Id   *string `json:"id"`
		Qty  *string `json:"qty"`
		Aksi *string `json:"aksi"`
	}

	type TemplateInputUser struct {
		Date     *string `json:"date"`
		GudangId *string `json:"gudangId"`
		UserId   *string `json:"userId"`
		Type     *string `json:"type"`
		Items    []Items `json:"items"`
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

	var StokMerchandiserRiwayat []structs.StokMerchandiserRiwayat

	parentId := 5 + int64(helpers.ParseInt(*inputUser.UserId)) + time.Now().Unix()
	for i := 0; i < len(inputUser.Items); i++ {
		// 	tempString := *inputUser.Items[i].Id + "-" + *inputUser.Items[i].Pita + "-" + *inputUser.Items[i].Condition + "-" + *inputUser.Items[i].Aksi + *inputUser.UserId

		// 	if
		StokMerchandiserRiwayat = append(StokMerchandiserRiwayat, structs.StokMerchandiserRiwayat{
			ItemId:         int16(helpers.ParseInt(*inputUser.Items[i].Id)),
			Jumlah:         int32(helpers.ParseInt(*inputUser.Items[i].Qty)),
			IsValidate:     0,
			GudangId:       int16(helpers.ParseInt(*inputUser.GudangId)),
			UserId:         int32(helpers.ParseInt(*inputUser.UserId)),
			TanggalRiwayat: helpers.ParseDate(*inputUser.Date),
			ParentId:       parentId,
			Aksi:           *inputUser.Items[i].Aksi,
		})
	}

	tx := db.DB.Begin()

	err = tx.Create(&StokMerchandiserRiwayat).Error
	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan input data",
			Success: false,
		})
	}

	tx.Commit()

	params := map[string]interface{}{
		"warehouseId": *inputUser.GudangId,
	}

	dataSend, err := json.Marshal(params)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Error marshaling JSON",
			Success: false,
		})
	}

	_, err = helpers.SendCurl(dataSend, "POST", "https://api.gudangku.pt-bks.com/order/new-order-item-notify")
	if err != nil {
		fmt.Println("Error sending request:", err)
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mengirim notification",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
		Message: "Success insert data",
		Success: true,
	})
}
