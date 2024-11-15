package controllers

import (
	"encoding/json"
	"fmt"
	db "pluto_remastered/config"
	"pluto_remastered/helpers"
	"pluto_remastered/structs"
	"strconv"
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

func getStokParent(userId *string, date *string, userIdSubtitute *string, gudangId *string, c *fiber.Ctx) ([]map[string]interface{}, error) {

	// fmt.Println("from function 1", userId, date, userIdSubtitute, gudangId)
	// fmt.Println("from function 2",*userId, *date, *userIdSubtitute, *gudangId)

	templateQuery := `WITH penjualans as (
							SELECT p.user_id, DATE(p.tanggal_penjualan) as dates, pd.produk_id, pd.condition, pd.pita, SUM(pd.jumlah) as qty
							FROM penjualan p
							JOIN penjualan_detail pd
								ON p.id = pd.penjualan_id
							WHERE p.user_id = {{.QInputUserId}} AND DATE(p.tanggal_penjualan) = '{{.QInputDate}}' {{.QPWhereSubtitute}}
							GROUP BY p.user_id, DATE(p.tanggal_penjualan), pd.produk_id, pd.condition, pd.pita
						), pengembalians as (
							SELECT p.user_id, DATE(p.tanggal_pengembalian) as dates, pd.produk_id, pd.condition, pd.pita, SUM(pd.jumlah) as qty
							FROM pengembalian p
							JOIN pengembalian_detail pd
								ON p.id = pd.pengembalian_id
							WHERE p.user_id = {{.QInputUserId}} AND DATE(p.tanggal_pengembalian) = '{{.QInputDate}}' {{.QPWhereSubtitute}}
							GROUP BY p.user_id, DATE(p.tanggal_pengembalian), pd.produk_id, pd.condition, pd.pita
						), order_gudangs as (
							SELECT ssr.user_id, DATE(ssr.tanggal_riwayat) as dates, ssr.produk_id, ssr.condition, ssr.pita, SUM(ssr.jumlah) as qty
							FROM stok_salesman_riwayat ssr
							WHERE ssr.user_id = {{.QInputUserId}} AND DATE(ssr.tanggal_riwayat) = '{{.QInputDate}}' AND ssr.aksi = 'ORDER' {{.QSsrWhereSubtitute}}
							GROUP BY ssr.user_id, DATE(ssr.tanggal_riwayat) , ssr.produk_id, ssr.condition, ssr.pita
						), retur_gudangs as (
							SELECT ssr.user_id, DATE(ssr.tanggal_riwayat) as dates, ssr.produk_id, ssr.condition, ssr.pita, SUM(ssr.jumlah) as qty
							FROM stok_salesman_riwayat ssr
							WHERE ssr.user_id = {{.QInputUserId}} AND DATE(ssr.tanggal_riwayat) = '{{.QInputDate}}' AND ssr.aksi = 'RETUR' {{.QSsrWhereSubtitute}}
							GROUP BY ssr.user_id, DATE(ssr.tanggal_riwayat) , ssr.produk_id, ssr.condition, ssr.pita
						), transactions as (
							SELECT tr.user_id, DATE(tr.datetime) as dates, trd.item_id, SUM(trd.qty) as qty
							FROM md.transaction tr
							JOIN md.transaction_detail trd
								ON tr.id = trd.transaction_id
							WHERE tr.user_id = {{.QInputUserId}} AND DATE(tr.datetime) = '{{.QInputDate}}' {{.QTrWhereSubtitute}}
							GROUP BY tr.user_id, DATE(tr.datetime), trd.item_id
						), order_items as (
							SELECT smr.user_id, DATE(smr.tanggal_riwayat) as dates, smr.item_id, SUM(smr.jumlah) as qty
							FROM md.stok_merchandiser_riwayat smr
							WHERE smr.user_id = {{.QInputUserId}} AND DATE(smr.tanggal_riwayat) = '{{.QInputDate}}' AND smr.aksi = 'ORDER' {{.QSmrWhereSubtitute}}
							GROUP BY smr.user_id, DATE(smr.tanggal_riwayat), smr.item_id
						), retur_items as (
							SELECT smr.user_id, DATE(smr.tanggal_riwayat) as dates, smr.item_id, SUM(smr.jumlah) as qty
							FROM md.stok_merchandiser_riwayat smr
							WHERE smr.user_id = {{.QInputUserId}} AND DATE(smr.tanggal_riwayat) = '{{.QInputDate}}' AND smr.aksi = 'RETUR' {{.QSmrWhereSubtitute}}
							GROUP BY smr.user_id, DATE(smr.tanggal_riwayat), smr.item_id
						), stok_salesmans as (
							SELECT ss.stok_user_id,
											JSONB_AGG(
												JSONB_BUILD_OBJECT(
													'id', ss.id,
													'stok_gudang_id', ss.stok_gudang_id,
													'user_id', ss.user_id,
													'produk_id', ss.produk_id,
													'tanggal_stok', ss.tanggal_stok,
													'dtm_crt', ss.dtm_crt,
													'dtm_upd', ss.dtm_upd,
													'confirm_key', ss.confirm_key,
													'is_complete', ss.is_complete,
													'tanggal_so', ss.tanggal_so,
													'so_admin_gudang_id', ss.so_admin_gudang_id,
													'condition', ss.condition,
													'pita', ss.pita,
													'stok_awal', ss.stok_awal,
													'orders', sso.qty,
													'returs', ssr.qty,
													'penjualan', pj.qty,
													'pengembalian', pg.qty,
													'stok_akhir', ss.stok_akhir
												)
											) as detail_produks
							FROM stok_salesman ss
							LEFT JOIN penjualans pj
								ON ss.produk_id = pj.produk_id
								AND ss.condition = pj.condition
								AND ss.pita = pj.pita
								AND DATE(ss.tanggal_stok) = pj.dates
								AND ss.user_id = pj.user_id
							LEFT JOIN pengembalians pg
								ON ss.produk_id = pg.produk_id
								AND ss.condition = pg.condition
								AND ss.pita = pg.pita
								AND DATE(ss.tanggal_stok) = pg.dates
								AND ss.user_id = pg.user_id
							LEFT JOIN order_gudangs sso
								ON ss.produk_id = sso.produk_id
								AND ss.condition = sso.condition
								AND ss.pita = sso.pita
								AND DATE(ss.tanggal_stok) = sso.dates
								AND ss.user_id = sso.user_id
							LEFT JOIN retur_gudangs ssr
								ON ss.produk_id = ssr.produk_id
								AND ss.condition = ssr.condition
								AND ss.pita = ssr.pita
								AND DATE(ss.tanggal_stok) = ssr.dates
								AND ss.user_id = ssr.user_id
							WHERE ss.user_id = {{.QInputUserId}} AND DATE(ss.tanggal_stok) = '{{.QInputDate}}' {{.QSSWhereSubtitute}}
							GROUP BY ss.stok_user_id
						), stok_merchandisers as (
							SELECT ss.stok_user_id,
											JSONB_AGG(
												JSONB_BUILD_OBJECT(
													'id', ss.id,
													'stok_gudang_id', ss.stok_gudang_id,
													'user_id', ss.user_id,
													'item_id', ss.item_id,
													'tanggal_stok', ss.tanggal_stok,
													'dtm_crt', ss.dtm_crt,
													'dtm_upd', ss.dtm_upd,
													'confirm_key', ss.confirm_key,
													'is_complete', ss.is_complete,
													'tanggal_so', ss.tanggal_so,
													'so_admin_gudang_id', ss.so_admin_gudang_id,
													'stok_awal', ss.stok_awal,
													'orders', sso.qty,
													'returs', ssr.qty,
													'penjualan', tr.qty,
													'pengembalian', 0,
													'stok_akhir', ss.stok_akhir
												)
											) as detail_items
							FROM md.stok_merchandiser ss 
							JOIN md.item i
								ON ss.item_id = i.id
							LEFT JOIN item_unit iu
								ON i.unit_id = iu.id
							LEFT JOIN transactions tr
								ON ss.item_id = tr.item_id
								AND DATE(ss.tanggal_stok) = tr.dates
								AND ss.user_id = tr.user_id
							LEFT JOIN order_items sso
								ON ss.item_id = sso.item_id
								AND DATE(ss.tanggal_stok) = sso.dates
								AND ss.user_id = sso.user_id
							LEFT JOIN retur_items ssr
								ON ss.item_id = ssr.item_id
								AND DATE(ss.tanggal_stok) = ssr.dates
								AND ss.user_id = ssr.user_id
							WHERE ss.user_id = {{.QInputUserId}} AND DATE(ss.tanggal_stok) = '{{.QInputDate}}' {{.QSSWhereSubtitute}}
							GROUP BY ss.stok_user_id
						)

						SELECT su.tanggal_stok,
								su.gudang_id,
								su.is_complete,
								su.tanggal_so,
								CASE WHEN su.user_id_subtitute IS NOT NULL AND su.user_id_subtitute <> 0 THEN su.user_id_subtitute ELSE su.user_id END as user_id,
								COALESCE(subs.full_name, u.full_name) as name,
								CASE WHEN su.user_id_subtitute IS NOT NULL AND su.user_id_subtitute <> 0 THEN 1 ELSE 0 END as is_subtitute,
								CASE WHEN su.user_id_subtitute IS NOT NULL AND su.user_id_subtitute <> 0 THEN u.id ELSE NULL END as account_owner_id,
								CASE WHEN su.user_id_subtitute IS NOT NULL AND su.user_id_subtitute <> 0 THEN u.full_name ELSE NULL END as account_owner_name,
								ss.detail_produks,
								sm.detail_items
						FROM stok_user su
						JOIN public.user u
							ON su.user_id = u.id
						LEFT JOIN public.user subs
							ON su.user_id_subtitute = subs.id
							AND su.user_id_subtitute <> 0
						LEFT JOIN stok_salesmans ss
							ON su.id = ss.stok_user_id
						LEFT JOIN stok_merchandisers sm
							ON su.id = sm.stok_user_id
						WHERE DATE(su.tanggal_stok) = DATE('{{.QInputDate}}')
							AND su.user_id = {{.QInputUserId}}
							AND su.gudang_id = {{.QInputGudangId}}
							{{.QInputUserIdSubtitute}}`

	where := ""
	if userIdSubtitute != nil {
		if *userIdSubtitute != "" {
			where = " AND su.user_id_subtitute = " + *userIdSubtitute
		}
	} else {
		where = " AND (su.user_id_subtitute = 0 OR su.user_id_subtitute IS NULL)"
	}

	templateParamQuery := make(map[string]interface{})

	if userIdSubtitute != nil {
		if *userIdSubtitute != "" {
			templateParamQuery = map[string]interface{}{
				"QInputUserId":          *userId,
				"QInputDate":            *date,
				"QInputUserIdSubtitute": where,
				"QInputGudangId":        *gudangId,
				"QPWhereSubtitute":      " AND p.user_id_subtitute = " + *userIdSubtitute,
				"QSsrWhereSubtitute":    " AND ssr.user_id_subtitute = " + *userIdSubtitute,
				"QTrWhereSubtitute":     " AND tr.user_id_subtitute = " + *userIdSubtitute,
				"QSmrWhereSubtitute":    " AND smr.user_id_subtitute = " + *userIdSubtitute,
				"QSSWhereSubtitute":     " AND ss.user_id_subtitute = " + *userIdSubtitute,
			}
		} else {
			templateParamQuery = map[string]interface{}{
				"QInputUserId":          *userId,
				"QInputDate":            *date,
				"QInputUserIdSubtitute": where,
				"QInputGudangId":        *gudangId,
				"QPWhereSubtitute":      "",
				"QSsrWhereSubtitute":    "",
				"QTrWhereSubtitute":     "",
				"QSmrWhereSubtitute":    "",
				"QSSWhereSubtitute":     "",
			}
		}
	} else {
		templateParamQuery = map[string]interface{}{
			"QInputUserId":          *userId,
			"QInputDate":            *date,
			"QInputUserIdSubtitute": where,
			"QInputGudangId":        *gudangId,
			"QPWhereSubtitute":      "",
			"QSsrWhereSubtitute":    "",
			"QTrWhereSubtitute":     "",
			"QSmrWhereSubtitute":    "",
			"QSSWhereSubtitute":     "",
		}
	}

	query1, err := helpers.PrepareQuery(templateQuery, templateParamQuery)

	if err != nil {
		fmt.Println(err)
		return nil, c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika generate query",
			Success: false,
		})
	}

	// fmt.Println(query1)

	returnData, err := helpers.ExecuteQuery(query1)

	if err != nil {
		fmt.Println(err)
		return nil, c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika eksekusi query",
			Success: false,
		})
	}

	if len(returnData) == 0 {
		return nil, c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data stok tidak ditemukan",
			Success: false,
		})
	}

	return returnData, nil
}

func GetStoks(c *fiber.Ctx) error {
	type TemplateInputUser struct {
		UserId          *string `json:"userId"`
		Date            *string `json:"date"`
		UserIdSubtitute *string `json:"userIdSubtitute"`
		BranchID        *string `json:"branchId"`
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

	gudang := new(structs.Gudang)
	if err := db.DB.Where("branch_id = ? ", *inputUser.BranchID).First(&gudang).Error; err != nil && err.Error() != "record not found" {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal mendapatkan data gudang",
			Success: false,
		})
	}

	sGudangID := strconv.Itoa(int(gudang.ID))
	// fmt.Println(inputUser.UserId, inputUser.Date, inputUser.UserIdSubtitute, sGudangID)
	datas, err := getStokParent(inputUser.UserId, inputUser.Date, inputUser.UserIdSubtitute, &sGudangID, c)

	// fmt.Println(datas)
	if err != nil {
		return err
	}

	if len(datas) > 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.Response{
			Message: "Data tidak ditemukan",
			Success: true,
			Data:    make(map[string]interface{}),
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Berhasil mendapatkan data",
		Success: true,
		Data:    datas[0],
		// Data:    nil,
	})
}
