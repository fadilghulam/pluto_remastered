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

	// fmt.Println(where)

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
															JSONB_BUILD_OBJECT(
																'id', sg.item_id, 
																'name', i.name, 
																'code', i.code
															) as item,
															JSONB_BUILD_OBJECT(
																'id', i.category_id,
																'name', ic.name
															) as category,
															JSONB_BUILD_OBJECT(
																'id', i.brand_id,
																'name', pb.name
															) as brand,
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
			Message: "Something's wrong with your input",
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

func GetStokProduk(c *fiber.Ctx) error {
	type TemplateInputUser struct {
		UserId *string `json:"userId"`
		Date   *string `json:"date"`
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

	if inputUser.Date != nil {
		where = " AND DATE(ss2.tanggal_stok) <= DATE('" + *inputUser.Date + "')"
	} else {
		where = " AND DATE(ss2.tanggal_stok) = CURRENT_DATE"
	}

	dataMax, err := helpers.ExecuteQuery(
		fmt.Sprintf(`SELECT 
									MAX (ss2.tanggal_stok) AS tgl, 
									ss2.user_id 
									FROM 
									PUBLIC.stok_salesman ss2 
									WHERE ss2.user_id IN(%s) %s 
									GROUP BY ss2.user_id
									LIMIT 1`, *inputUser.UserId, where))

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data stok",
			Success: false,
		})
	}

	if len(dataMax) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.Response{
			Message: "Data stok tidak ditemukan",
			Success: false,
			Data:    make([]interface{}, 0),
		})
	}

	templateQuery := `WITH ssr AS
                (
                    SELECT ssr.user_id, ssr.produk_id, ssr.condition, ssr.pita,
                    SUM(CASE WHEN aksi ='ORDER' THEN COALESCE(ssr.jumlah,0) ELSE 0 END) AS order,
                    SUM(CASE WHEN aksi ='RETUR' THEN COALESCE(ssr.jumlah,0) ELSE 0 END) AS retur
                    FROM stok_salesman_riwayat ssr 
                    WHERE TRUE AND ssr.is_validate = 1 AND ssr.user_id = {{.QDataMaxUserId}} 
                    AND DATE(ssr.tanggal_riwayat) = DATE('{{.QDataMaxTanggal}}')
                    GROUP BY ssr.user_id, ssr.produk_id, ssr.condition, ssr.pita
                ), penjualans AS
				(
					SELECT p.user_id, pd.produk_id, pd.condition, pd.pita, SUM(pd.jumlah) as jumlah
					FROM penjualan p
					JOIN penjualan_detail pd
					ON p.id = pd.penjualan_id
					WHERE p.user_id = {{.QDataMaxUserId}} AND DATE(p.tanggal_penjualan) = DATE('{{.QDataMaxTanggal}}')
					GROUP BY p.user_id, pd.produk_id, pd.condition, pd.pita
				), pengembalians AS
				(
					SELECT p.user_id, pd.produk_id, pd.condition, pd.pita, SUM(pd.jumlah) as jumlah
					FROM pengembalian p
					JOIN pengembalian_detail pd
					ON p.id = pd.pengembalian_id
					WHERE p.user_id = {{.QDataMaxUserId}} AND DATE(p.tanggal_pengembalian) = DATE('{{.QDataMaxTanggal}}')
					GROUP BY p.user_id, pd.produk_id, pd.condition, pd.pita
				)

                    SELECT  ss.id, 
                            ss.stok_gudang_id, 
                            ss.user_id, 
                            ss.produk_id, 
                            ss.tanggal_stok, 
                            ss.dtm_crt, 
                            ss.dtm_upd, 
                            ss.confirm_key, 
                            ss.is_complete, 
                            ss.tanggal_so, 
                            ss.so_admin_gudang_id, 
                            ss.condition, 
                            ss.pita, 
                            (ss.stok_awal - SUM(COALESCE(ssr.order,0))) as stok_awal, 
                            SUM(COALESCE(ssr.order,0)) orders, 
                            SUM(COALESCE(ssr.order,0)) as returs,
							SUM(COALESCE(pj.jumlah,0)) as penjualan,
							SUM(COALESCE(pg.jumlah,0)) as pengembalian,
                            ss.stok_akhir 
                    FROM
                    PUBLIC.stok_salesman ss
                    LEFT JOIN ssr
                        ON ss.user_id = ssr.user_id
                        AND ss.produk_id = ssr.produk_id
                        AND ss.condition = ssr.condition
                        AND ss.pita = ssr.pita
					LEFT JOIN penjualans pj
						ON ss.user_id = pj.user_id
						AND ss.produk_id = pj.produk_id
                        AND ss.condition = pj.condition
                        AND ss.pita = pj.pita
					LEFT JOIN pengembalians pg
						ON ss.user_id = pg.user_id
						AND ss.produk_id = pg.produk_id
                        AND ss.condition = pg.condition
                        AND ss.pita = pg.pita
                    WHERE ss.condition = ('GOOD') AND ss.user_id = {{.QDataMaxUserId}} AND DATE(ss.tanggal_stok) = DATE('{{.QDataMaxTanggal}}')
                    GROUP BY ss.id`

	templateParamQuery := map[string]interface{}{
		"QDataMaxUserId":  dataMax[0]["user_id"],
		"QDataMaxTanggal": dataMax[0]["tgl"],
	}

	query1, err := helpers.PrepareQuery(templateQuery, templateParamQuery)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	returnData, err := helpers.NewExecuteQuery(query1)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika eksekusi query",
			Success: false,
		})
	}

	if len(returnData) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data stok tidak ditemukan",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data stok berhasil diambil",
		Data:    returnData,
		Success: true,
	})
}

func GetStokItem(c *fiber.Ctx) error {
	type TemplateInputUser struct {
		UserId *string `json:"userId"`
		Date   *string `json:"date"`
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

	if inputUser.Date != nil {
		where = " AND DATE(ss2.tanggal_stok) <= DATE('" + *inputUser.Date + "')"
	} else {
		where = " AND DATE(ss2.tanggal_stok) = CURRENT_DATE"
	}

	dataMax, err := helpers.ExecuteQuery(
		fmt.Sprintf(`SELECT 
						MAX (ss2.tanggal_stok) AS tgl, 
						ss2.user_id 
						FROM 
						md.stok_merchandiser ss2 
						WHERE ss2.user_id IN(%s) %s 
						GROUP BY ss2.user_id
						LIMIT 1`, *inputUser.UserId, where))

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data stok",
			Success: false,
		})
	}

	if len(dataMax) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.Response{
			Message: "Data stok tidak ditemukan",
			Success: false,
			Data:    make([]interface{}, 0),
		})
	}

	templateQuery := `WITH ssr AS
                (
                    SELECT ssr.user_id, ssr.item_id,
                    SUM(CASE WHEN aksi ='ORDER' THEN COALESCE(ssr.jumlah,0) ELSE 0 END) AS order,
                    SUM(CASE WHEN aksi ='RETUR' THEN COALESCE(ssr.jumlah,0) ELSE 0 END) AS retur
					FROM
            		md.stok_merchandiser_riwayat ssr 
                    WHERE TRUE AND ssr.is_validate = 1 AND ssr.user_id = {{.QDataMaxUserId}} 
                    AND DATE(ssr.tanggal_riwayat) = DATE('{{.QDataMaxTanggal}}')
                    GROUP BY ssr.user_id, ssr.item_id
                ), transactions as (
					SELECT tr.user_id, trd.item_id, SUM(trd.qty) as jumlah
					FROM md.transaction tr
					JOIN md.transaction_detail trd
					ON tr.id = trd.transaction_id
					WHERE tr.user_id = {{.QDataMaxUserId}} AND DATE(tr.datetime) = DATE('{{.QDataMaxTanggal}}')
                    GROUP BY tr.user_id, trd.item_id
				)

                    SELECT  ss.id, 
                            ss.stok_gudang_id, 
                            ss.user_id, 
                            ss.item_id, 
                            ss.tanggal_stok, 
                            ss.dtm_crt, 
                            ss.dtm_upd, 
                            ss.confirm_key, 
                            ss.is_complete, 
                            ss.tanggal_so, 
                            ss.so_admin_gudang_id,
                            (ss.stok_awal - SUM(COALESCE(ssr.order,0))) as stok_awal, 
                            SUM(COALESCE(ssr.order,0)) orders, 
                            SUM(COALESCE(ssr.order,0)) as returs, 
							SUM(COALESCE(tr.jumlah,0)) as penjualan,
							0 as pengembalian,
                            ss.stok_akhir 
                    FROM
                    md.stok_merchandiser ss
                    LEFT JOIN ssr
                        ON ss.user_id = ssr.user_id
                        AND ss.item_id = ssr.item_id
					LEFT JOIN transactions tr
						ON ss.user_id = tr.user_id
						AND ss.item_id = tr.item_id
                    WHERE ss.user_id = {{.QDataMaxUserId}} AND DATE(ss.tanggal_stok) = DATE('{{.QDataMaxTanggal}}')
                    GROUP BY ss.id`

	templateParamQuery := map[string]interface{}{
		"QDataMaxUserId":  dataMax[0]["user_id"],
		"QDataMaxTanggal": dataMax[0]["tgl"],
	}

	query1, err := helpers.PrepareQuery(templateQuery, templateParamQuery)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	returnData, err := helpers.NewExecuteQuery(query1)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika eksekusi query",
			Success: false,
		})
	}

	if len(returnData) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data stok tidak ditemukan",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data stok berhasil diambil",
		Data:    returnData,
		Success: true,
	})
}
